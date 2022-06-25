// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// FPmsAPI 文件权限管理api

package controller

import (
	"encoding/json"
	"fileservice/biz/constants"
	"fileservice/biz/service"
	"net/http"
	"strconv"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/serviceutil"
)

// FilePermissionCtrl 用户权限管理
type FilePermissionCtrl struct {
	um  service.UserAuth4Rpc   `@autowired:"User4RPC"`
	pms service.FilePermission `@autowired:"FilePermission"`
}

// AsController 实现 AsController 接口
func (ctl *FilePermissionCtrl) AsController() ipakku.ControllerConfig {
	return ipakku.ControllerConfig{
		RequestMapping: "/filepms/v1",
		RouterConfig: ipakku.RouterConfig{
			ToLowerCase: true,
			HandlerFunc: [][]interface{}{
				{http.MethodGet, ctl.ListFPermissions},
				{http.MethodPost, ctl.GetUserPermissionSum},
				{http.MethodGet, ctl.ListUserFPermissions},
				{http.MethodPost, ctl.AddFPermission},
				{http.MethodDelete, ctl.DelFPermission},
				{http.MethodPost, ctl.UpdateFPermission},
			},
		},
		FilterConfig: ipakku.FilterConfig{
			FilterFunc: [][]interface{}{
				{`/:[\s\S]*`, ctl.um.GetAuthFilterFunc()},
			},
		},
	}
}

// checkPermission 检查是否是管理员
func (ctl *FilePermissionCtrl) checkPermission(w http.ResponseWriter, r *http.Request) bool {
	if accesskey := ctl.um.GetAccessKey4Request(r); len(accesskey) > 0 {
		if ack, err := ctl.um.GetUserAccess(accesskey); nil == err {
			if len(ack.UserID) > 0 {
				// 管理员
				if ack.UserType == service.UserType_Admin {
					return true
					// 自己的账户
				} else if qUserID := r.FormValue("userid"); len(qUserID) > 0 && qUserID == ack.UserID {
					return true
				}
			}
		}
	}
	w.WriteHeader(http.StatusForbidden)
	return false
}

// ListFPermissions 列出所有数据, 无分页
func (ctl *FilePermissionCtrl) ListFPermissions(w http.ResponseWriter, r *http.Request) {
	if !ctl.checkPermission(w, r) {
		return
	}
	if pms, err := ctl.pms.ListFPermissions(); nil == err {
		pmsdto := make([]*service.PermissionInfoDto, len(pms))
		for i := 0; i < len(pms); i++ {
			pmsdto[i] = pms[i].ToDto()
		}
		serviceutil.SendSuccess(w, pmsdto)
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// GetUserPermissionSum 根据用户ID查询权限值
func (ctl *FilePermissionCtrl) GetUserPermissionSum(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	if len(userID) == 0 {
		serviceutil.SendBadRequest(w, service.ErrorUserIDIsNil.Error())
		return
	}
	var paths []string
	if err := json.Unmarshal([]byte(r.FormValue("paths")), &paths); nil != err || len(paths) == 0 {
		serviceutil.SendBadRequest(w, service.ErrorPermissionPathIsNil.Error())
		return
	}
	if !ctl.checkPermission(w, r) {
		return
	}
	permissions := make(map[string]int64, len(paths))
	for i := 0; i < len(paths); i++ {
		permissions[paths[i]] = ctl.pms.GetUserPermissionSum(userID, paths[i])
	}

	serviceutil.SendSuccess(w, permissions)
}

// ListUserFPermissions 根据用户ID查询详细信息
func (ctl *FilePermissionCtrl) ListUserFPermissions(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	if len(userID) == 0 {
		serviceutil.SendBadRequest(w, service.ErrorUserIDIsNil.Error())
		return
	}
	if !ctl.checkPermission(w, r) {
		return
	}
	if pms, err := ctl.pms.ListUserFPermissions(userID); nil == err {
		pmsdto := make([]*service.PermissionInfoDto, len(pms))
		for i := 0; i < len(pms); i++ {
			pmsdto[i] = pms[i].ToDto()
		}
		serviceutil.SendSuccess(w, pmsdto)
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// AddFPermission 添加
func (ctl *FilePermissionCtrl) AddFPermission(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	path := r.FormValue("path")
	permissionStr := r.FormValue("permission")
	if len(userID) == 0 {
		serviceutil.SendBadRequest(w, service.ErrorUserIDIsNil.Error())
		return
	}
	if userID == constants.AdminUserID {
		serviceutil.SendBadRequest(w, "不能给管理员添加权限")
		return
	}
	if len(path) == 0 {
		serviceutil.SendBadRequest(w, service.ErrorPermissionPathIsNil.Error())
		return
	}
	if len(permissionStr) == 0 {
		serviceutil.SendBadRequest(w, service.ErrorPermissionIsNil.Error())
		return
	}
	permission, err := strconv.Atoi(permissionStr)
	if nil != err {
		serviceutil.SendBadRequest(w, err.Error())
		return
	}
	if !ctl.checkPermission(w, r) {
		return
	}
	pmsInfo := service.PermissionInfo{
		UserID:     userID,
		Path:       path,
		Permission: int64(permission),
	}

	if err := ctl.pms.AddFPermission(pmsInfo); nil == err {
		serviceutil.SendSuccess(w, "")
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// UpdateFPermission 修改
func (ctl *FilePermissionCtrl) UpdateFPermission(w http.ResponseWriter, r *http.Request) {
	permissionID := r.FormValue("permissionid")
	permissionStr := r.FormValue("permission")
	if len(permissionStr) == 0 {
		serviceutil.SendBadRequest(w, service.ErrorPermissionIDIsNil.Error())
		return
	}
	permission, err := strconv.Atoi(permissionStr)
	if nil != err {
		serviceutil.SendBadRequest(w, err.Error())
		return
	}

	if !ctl.checkPermission(w, r) {
		return
	}
	pmsInfo := service.PermissionInfo{
		PermissionID: permissionID,
		Permission:   int64(permission),
	}
	if err := ctl.pms.UpdateFPermission(pmsInfo); nil == err {
		serviceutil.SendSuccess(w, "")
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// DelFPermission 根据ID删除
func (ctl *FilePermissionCtrl) DelFPermission(w http.ResponseWriter, r *http.Request) {
	permissionID := r.FormValue("permissionid")
	if len(permissionID) == 0 {
		serviceutil.SendBadRequest(w, service.ErrorPermissionIDIsNil.Error())
		return
	}
	if !ctl.checkPermission(w, r) {
		return
	}

	if err := ctl.pms.DelFPermission(permissionID); nil == err {
		serviceutil.SendSuccess(w, "")
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}
