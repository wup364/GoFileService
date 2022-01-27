// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件预览接口模块

package controller

import (
	"fileservice/biz/modules/filedatas"
	"fileservice/biz/service"
	"io"
	"net/http"
	"pakku/ipakku"
	"pakku/utils/logs"
	"pakku/utils/serviceutil"
	"pakku/utils/strutil"
	"strconv"
	"strings"
)

// Preview 文件预览模块
type Preview struct {
	fm  service.FileDatas           `@autowired:"FileDatas"`
	um  service.UserAuth4Rpc        `@autowired:"User4RPC"`
	tt  service.TransportToken      `@autowired:"TransportToken"`
	pmc service.FilePermissionCheck `@autowired:"FilePermission"`
}

// AsController 实现 AsController 接口
func (ctl *Preview) AsController() ipakku.ControllerConfig {
	return ipakku.ControllerConfig{
		RequestMapping: "/filepreview/v1",
		RouterConfig: ipakku.RouterConfig{
			ToLowerCase: true,
			HandlerFunc: [][]interface{}{
				{"GET", "status/:" + `[\s\S]*`, ctl.Status},
				{"GET", ctl.Asktoken},
				{"GET", "samedirfiles/:" + `[\s\S]*`, ctl.SameDirFiles},
				{"ANY", "stream/:" + `[\s\S]*`, ctl.Stream},
			},
		},
		FilterConfig: ipakku.FilterConfig{
			FilterFunc: [][]interface{}{
				{`/:[\s\S]*`, func(rw http.ResponseWriter, r *http.Request) bool {
					if strings.Contains(r.URL.Path, "/stream/") {
						return true
					}
					return ctl.um.GetAuthFilterFunc()(rw, r)
				}},
			},
		},
	}
}

// Status 保持会话用
func (ctl *Preview) Status(w http.ResponseWriter, r *http.Request) {
	if token, err := ctl.getSteamToken(strutil.GetPathName(r.URL.Path)); nil != err || nil == token {
		serviceutil.SendBadRequest(w, err.Error())
	} else if _, err = ctl.tt.RefreshToken(token.Token); nil != err {
		serviceutil.SendBadRequest(w, err.Error())
	} else {
		serviceutil.SendSuccess(w, "")
	}
}

// Asktoken 预览令牌申请
func (ctl *Preview) Asktoken(w http.ResponseWriter, r *http.Request) {
	if qpath := r.FormValue("path"); len(qpath) == 0 {
		serviceutil.SendBadRequest(w, ErrorOprationFailed.Error())
	} else if !ctl.fm.IsFile(qpath) {
		serviceutil.SendBadRequest(w, ErrorFileNotExist.Error())
	} else if !ctl.checkPermision(ctl.getUserID4Request(r), qpath, service.FPM_Read) {
		serviceutil.SendBadRequest(w, service.ErrorPermissionInsufficient.Error())
	} else {
		if token, err := ctl.tt.AskReadToken(qpath); nil != err {
			serviceutil.SendServerError(w, err.Error())
		} else {
			serviceutil.SendSuccess(w, token.Token)
		}
	}
}

// SameDirFiles 同目录文件
func (ctl *Preview) SameDirFiles(w http.ResponseWriter, r *http.Request) {
	var prentPath string
	token, err := ctl.getSteamToken(strutil.GetPathName(r.URL.Path))
	if nil != err || nil == token {
		serviceutil.SendBadRequest(w, err.Error())
	} else {
		if prentPath = strutil.GetPathParent(token.FilePath); !ctl.fm.IsExist(prentPath) {
			serviceutil.SendBadRequest(w, ErrorFileNotExist.Error())
		} else if !ctl.checkPermision(ctl.getUserID4Request(r), token.FilePath, service.FPM_Read) {
			serviceutil.SendBadRequest(w, service.ErrorPermissionInsufficient.Error())
		} else {
			if list, err := ctl.fm.GetDirNodeList(prentPath); err != nil {
				serviceutil.SendServerError(w, err.Error())
			} else {
				res := make([]service.FNodeDto, 0)
				if len(list) > 0 {
					for i := 0; i < len(list); i++ {
						res = append(res, *list[i].ToDto())
					}
				}
				//
				serviceutil.SendSuccess(w, PreviewSameDirFiles{
					Path: token.FilePath,
					PeerDatas: filedatas.FileListSorter{
						Asc:       true,
						SortField: "Path",
					}.Sort(res),
				})
			}
		}
	}
}

// Stream 下载
func (ctl *Preview) Stream(w http.ResponseWriter, r *http.Request) {
	var filePath string
	token, err := ctl.getSteamToken(strutil.GetPathName(r.URL.Path))
	if nil != err || nil == token {
		serviceutil.SendBadRequest(w, err.Error())
		return
	} else {
		// 校验
		parent := strutil.GetPathParent(token.FilePath)
		if filePath = strutil.Parse2UnixPath(parent + "/" + strutil.GetPathName(r.FormValue("fileName"))); !ctl.fm.IsFile(filePath) {
			serviceutil.SendBadRequest(w, ErrorFileNotExist.Error())
			return
		}
	}
	//
	maxSize := ctl.fm.GetFileSize(filePath)
	start, end, hasRange := serviceutil.GetRequestRange(r, maxSize)
	//
	if fr, err := ctl.fm.DoRead(filePath, start); nil != err {
		serviceutil.SendServerError(w, err.Error())
	} else {
		defer fr.Close()
		ctLength := end - start
		if ctLength == 0 || ctLength < 0 {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.Header().Set("Content-Length", strconv.Itoa(int(ctLength)))
			if hasRange {
				w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(start))+"-"+strconv.Itoa(int(end-1))+"/"+strconv.Itoa(int(maxSize)))
				w.WriteHeader(http.StatusPartialContent)
			}
			if _, err := io.Copy(w, io.LimitReader(fr, ctLength)); nil != err && err != io.EOF {
				logs.Errorln(err)
				serviceutil.SendServerError(w, err.Error())
			}
		}
	}
}

// checkPermision 检查权限
func (ctl *Preview) checkPermision(userID, path string, permission int64) bool {
	return ctl.pmc.HashPermission(userID, path, permission)
}

// getUserID4Request 获取登录用户
func (ctl *Preview) getUserID4Request(r *http.Request) string {
	if askstr := ctl.um.GetAccessKey4Request(r); len(askstr) > 0 {
		if ack, err := ctl.um.GetUserAccess(askstr); nil == err {
			return ack.UserID
		}
	}
	return ""
}

// getSteamToken 获取文件传输Token对象
func (ctl *Preview) getSteamToken(token string) (*service.StreamToken, error) {
	if tokenBody, err := ctl.tt.QueryToken(token); nil != err || nil == tokenBody {
		return nil, ErrorOprationExpires
	} else {
		return tokenBody, nil
	}
}
