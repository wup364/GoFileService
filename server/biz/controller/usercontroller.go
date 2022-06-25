// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// UserAPI 用户管理api

package controller

import (
	"fileservice/biz/service"
	"net/http"
	"strings"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/serviceutil"
)

// UserCtrl 用户管理api
type UserCtrl struct {
	um service.User4RPC `@autowired:"User4RPC"`
}

// AsController 实现 AsController 接口
func (ctl *UserCtrl) AsController() ipakku.ControllerConfig {
	return ipakku.ControllerConfig{
		RequestMapping: "/user/v1",
		RouterConfig: ipakku.RouterConfig{
			ToLowerCase: true,
			HandlerFunc: [][]interface{}{
				{"GET", ctl.ListAllUsers},
				{"GET", ctl.QueryUser},
				{"POST", ctl.AddUser},
				{"DELETE", ctl.DelUser},
				{"POST", ctl.UpdateUserName},
				{"POST", ctl.UpdateUserPwd},
				{"POST", ctl.CheckPwd},
				{"POST", ctl.Logout},
			},
		},
		FilterConfig: ipakku.FilterConfig{
			FilterFunc: [][]interface{}{
				{`/:[\s\S]*`, func(rw http.ResponseWriter, r *http.Request) bool {
					if strings.HasSuffix(r.URL.Path, "/checkpwd") {
						return true
					}
					return ctl.um.GetAuthFilterFunc()(rw, r)
				}},
			},
		},
	}
}

// checkPermission 检查是否是管理员
func (ctl *UserCtrl) checkPermission(w http.ResponseWriter, r *http.Request) bool {
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

// ListAllUsers 列出所有用户数据, 无分页
func (ctl *UserCtrl) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	if !ctl.checkPermission(w, r) {
		return
	}
	if users, err := ctl.um.ListAllUsers(); nil == err {
		serviceutil.SendSuccess(w, users)
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// QueryUser 根据用户ID查询详细信息
func (ctl *UserCtrl) QueryUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	if len(userID) == 0 {
		serviceutil.SendBadRequest(w, ErrorUserIDIsNil.Error())
		return
	}
	if !ctl.checkPermission(w, r) {
		return
	}
	if user, err := ctl.um.QueryUser(userID); nil == err {
		serviceutil.SendSuccess(w, user)
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// AddUser 添加用户
func (ctl *UserCtrl) AddUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	userName := r.FormValue("username")
	userPwd := r.FormValue("userpwd")
	if len(userID) == 0 {
		serviceutil.SendBadRequest(w, ErrorUserIDIsNil.Error())
		return
	}
	if len(userName) == 0 {
		serviceutil.SendBadRequest(w, ErrorUserNameIsNil.Error())
		return
	}
	if !ctl.checkPermission(w, r) {
		return
	}
	uinfo := service.UserInfo{
		UserID:   userID,
		UserName: userName,
		UserPWD:  userPwd,
		UserType: service.UserType_Normal,
	}

	if err := ctl.um.AddUser(&uinfo); nil == err {
		serviceutil.SendSuccess(w, "")
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// UpdateUserPwd 修改用户密码
func (ctl *UserCtrl) UpdateUserPwd(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	userPwd := r.FormValue("userpwd")
	if len(userID) == 0 {
		serviceutil.SendBadRequest(w, ErrorUserIDIsNil.Error())
		return
	}
	if !ctl.checkPermission(w, r) {
		return
	}
	if err := ctl.um.UpdatePWD(userID, userPwd); nil == err {
		serviceutil.SendSuccess(w, "")
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// UpdateUserName 修改用户昵称
func (ctl *UserCtrl) UpdateUserName(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	userName := r.FormValue("username")
	if len(userID) == 0 {
		serviceutil.SendBadRequest(w, ErrorUserIDIsNil.Error())
		return
	}
	if len(userName) == 0 {
		serviceutil.SendBadRequest(w, ErrorUserNameIsNil.Error())
		return
	}
	if !ctl.checkPermission(w, r) {
		return
	}
	if err := ctl.um.UpdateUserName(userID, userName); nil == err {
		serviceutil.SendSuccess(w, "")
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// DelUser 根据userID删除用户
func (ctl *UserCtrl) DelUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	if len(userID) == 0 {
		serviceutil.SendBadRequest(w, ErrorUserIDIsNil.Error())
		return
	}
	if !ctl.checkPermission(w, r) {
		return
	}
	users, _ := ctl.um.ListAllUsers()
	if nil == users || len(users) == 1 {
		serviceutil.SendBadRequest(w, "cannot delete the last user")
		return
	}
	{
		count := 0
		lauid := ""
		for _, val := range users {
			if val.UserType == service.UserType_Admin {
				count++
				lauid = val.UserID
			}
		}
		if count <= 1 && (userID == lauid || len(lauid) == 0) {
			serviceutil.SendBadRequest(w, "cannot delete the last admin user")
			return
		}
	}
	if err := ctl.um.DelUser(userID); nil == err {
		serviceutil.SendSuccess(w, "")
		// ctl.um.DestroyAccess(ctl.um.GetAccessKey4Request(r))
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// CheckPwd 校验密码是否一致,  校验成功返回session
func (ctl *UserCtrl) CheckPwd(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	pwd := r.FormValue("pwd")
	if len(userID) == 0 {
		serviceutil.SendBadRequest(w, ErrorUserIDIsNil.Error())
		return
	}
	// 检查密码是否正确, 如果正确需要返回签名信息
	if ack, err := ctl.um.AskAccess(userID, pwd); nil == err {
		serviceutil.SendSuccess(w, ack)
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// Logout 注销会话
func (ctl *UserCtrl) Logout(w http.ResponseWriter, r *http.Request) {
	if ack := ctl.um.GetAccessKey4Request(r); len(ack) > 0 {
		if err := ctl.um.DestroyAccess(ack); nil == err {
			serviceutil.SendSuccess(w, "")
		} else {
			serviceutil.SendServerError(w, err.Error())
		}
	}
}
