// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件操作接口, 文件的新建、删除、移动、复制等操作

package controller

import (
	"errors"
	"fileservice/biz/modules/filedatas"
	"fileservice/biz/service"
	"net/http"
	"path"
	"strings"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/serviceutil"
	"github.com/wup364/pakku/utils/strutil"
)

// FileOptsCtrl 文件操作接口
type FileOptsCtrl struct {
	fm  service.FileDatas           `@autowired:"FileDatas"`
	um  service.UserAuth4Rpc        `@autowired:"User4RPC"`
	pms service.FilePermissionCheck `@autowired:"FilePermission"`
}

// AsController 实现 AsController 接口
func (ctl *FileOptsCtrl) AsController() ipakku.ControllerConfig {
	return ipakku.ControllerConfig{
		RequestMapping: "/file/v1",
		RouterConfig: ipakku.RouterConfig{
			ToLowerCase: true,
			HandlerFunc: [][]interface{}{
				{http.MethodGet, ctl.Info},
				{http.MethodGet, ctl.List},
				{http.MethodDelete, ctl.Del},
				{http.MethodPost, ctl.ReName},
				{http.MethodPost, ctl.NewFolder},
			},
		},
		FilterConfig: ipakku.FilterConfig{
			FilterFunc: [][]interface{}{
				{`/:[\s\S]*`, ctl.um.GetAuthFilterFunc()},
			},
		},
	}
}

// checkPermision 检查权限
func (ctl *FileOptsCtrl) checkPermision(userID, path string, permission int64) bool {
	if len(path) == 0 || !strings.HasPrefix(path, "/") {
		return false
	}
	return ctl.pms.HashPermission(userID, path, permission)
}

// GetUserID4Request 获取登录用户
func (ctl *FileOptsCtrl) GetUserID4Request(r *http.Request) string {
	if askstr := ctl.um.GetAccessKey4Request(r); len(askstr) > 0 {
		if ack, err := ctl.um.GetUserAccess(askstr); nil == err {
			return ack.UserID
		}
	}
	return ""
}

// Info 获取文件|文件夹信息
func (ctl *FileOptsCtrl) Info(w http.ResponseWriter, r *http.Request) {
	qpath := path.Clean(r.FormValue("path"))
	if !ctl.checkPermision(ctl.GetUserID4Request(r), qpath, service.FPM_Visible) {
		serviceutil.SendBadRequest(w, ErrorPermissionInsufficient.Error())
		return
	}
	if len(qpath) == 0 {
		serviceutil.SendBadRequest(w, "path is empty")
		return
	}
	if !ctl.fm.IsExist(qpath) {
		serviceutil.SendBadRequest(w, qpath+" is not exist")
		return
	}
	serviceutil.SendSuccess(w, ctl.fm.GetNode(qpath).ToDto())
}

// List 查询路径下的列表以及基本信息
func (ctl *FileOptsCtrl) List(w http.ResponseWriter, r *http.Request) {
	qpath := r.FormValue("path")
	qlimit := r.FormValue("limit")
	qoffset := r.FormValue("offset")
	qSort := r.FormValue("sort")
	qAsc := r.FormValue("asc")
	// fmt.Println(r.URL, r.Body, qpath)
	if len(qpath) == 0 {
		serviceutil.SendBadRequest(w, errors.New("path is empty").Error())
		return
	}
	if len(qAsc) == 0 {
		qAsc = "true"
	}
	userID := ctl.GetUserID4Request(r)
	if !ctl.checkPermision(userID, qpath, service.FPM_VisibleChild) {
		serviceutil.SendBadRequest(w, ErrorPermissionInsufficient.Error())
		return
	}

	list, err := ctl.fm.GetDirNodeList(qpath, strutil.String2Int(qlimit, -1), strutil.String2Int(qoffset, -1))
	if err != nil {
		serviceutil.SendServerError(w, err.Error())
		return
	}
	// 如果当前或上级路径有可见以上权限, 则文件默认可见
	canVisible := ctl.checkPermision(userID, qpath, service.FPM_Visible)
	res := make([]service.FNodeDto, 0)
	if len(list) > 0 {
		for i := 0; i < len(list); i++ {
			if list[i].IsFile && canVisible {
				res = append(res, list[i].ToDto())
				continue
			}
			if ctl.checkPermision(userID, list[i].Path, service.FPM_VisibleChild) {
				res = append(res, list[i].ToDto())
			}
		}
	}
	//
	serviceutil.SendSuccess(w, filedatas.FileListSorter{
		Asc:       strutil.String2Bool(qAsc),
		SortField: qSort,
	}.Sort(res))
}

// Del 批量|单个删除文件|文件夹
func (ctl *FileOptsCtrl) Del(w http.ResponseWriter, r *http.Request) {
	qpath := r.FormValue("path")
	if len(qpath) == 0 {
		serviceutil.SendBadRequest(w, "path is empty")
		return
	}
	if !ctl.checkPermision(ctl.GetUserID4Request(r), qpath, service.FPM_Write) {
		serviceutil.SendBadRequest(w, ErrorPermissionInsufficient.Error())
		return
	}
	if !ctl.fm.IsExist(qpath) {
		serviceutil.SendBadRequest(w, ErrorFileNotExist.Error())
		return
	}
	if err := ctl.fm.DoDelete(qpath); nil == err {
		serviceutil.SendSuccess(w, "")
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// ReName 重命名
func (ctl *FileOptsCtrl) ReName(w http.ResponseWriter, r *http.Request) {
	qSrcPath := r.FormValue("path")
	qName := r.FormValue("name")
	if len(qName) == 0 {
		serviceutil.SendBadRequest(w, ErrorNewNameIsEmpty.Error())
		return
	}
	if !ctl.checkPermision(ctl.GetUserID4Request(r), qSrcPath, service.FPM_Write) {
		serviceutil.SendBadRequest(w, ErrorPermissionInsufficient.Error())
		return
	}
	if !ctl.fm.IsExist(qSrcPath) {
		serviceutil.SendBadRequest(w, ErrorFileNotExist.Error())
		return
	}
	if err := ctl.fm.DoRename(qSrcPath, qName); nil == err {
		serviceutil.SendSuccess(w, "")
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// NewFolder 新建文件夹
func (ctl *FileOptsCtrl) NewFolder(w http.ResponseWriter, r *http.Request) {
	qPath := r.FormValue("path")
	qPath = strutil.Parse2UnixPath(qPath)
	if !ctl.checkPermision(ctl.GetUserID4Request(r), qPath, service.FPM_Write) {
		serviceutil.SendBadRequest(w, ErrorPermissionInsufficient.Error())
		return
	}
	if !ctl.fm.IsExist(strutil.GetPathParent(qPath)) {
		serviceutil.SendBadRequest(w, ErrorParentFolderNotExist.Error())
		return
	}
	if err := ctl.fm.DoMkDir(qPath); nil == err {
		serviceutil.SendSuccess(w, "")
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}
