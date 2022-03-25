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
	"encoding/json"
	"fileservice/biz/modules/filedatas"
	"fileservice/biz/service"
	"net/http"
	"pakku/ipakku"
	"pakku/utils/serviceutil"
	"pakku/utils/strutil"
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
				{http.MethodGet, ctl.Asktoken},
				{http.MethodGet, "status/:" + `[\s\S]*`, ctl.Status},
				{http.MethodGet, "samedirfiles/:" + `[\s\S]*`, ctl.SameDirFiles},
				{http.MethodPost, "samedirtoken/:" + `[\s\S]*`, ctl.AskSameDIRToken},
			},
		},
		FilterConfig: ipakku.FilterConfig{
			FilterFunc: [][]interface{}{
				{`/:[\s\S]*`, func(rw http.ResponseWriter, r *http.Request) bool {
					if strings.Contains(r.URL.Path, "/asktoken") {
						return ctl.um.GetAuthFilterFunc()(rw, r)
					}
					return true
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
		if token, err := ctl.tt.AskReadToken(qpath, nil); nil != err {
			serviceutil.SendServerError(w, err.Error())
		} else {
			serviceutil.SendSuccess(w, token.Token)
		}
	}
}

// AskSameDIRToken 同级目录下的预览令牌申请
func (ctl *Preview) AskSameDIRToken(w http.ResponseWriter, r *http.Request) {
	var nameList []string
	if qnames := r.FormValue("names"); len(qnames) == 0 {
		serviceutil.SendBadRequest(w, ErrorOprationFailed.Error())
		return
	} else {
		if err := json.Unmarshal([]byte(qnames), &nameList); nil != err {
			serviceutil.SendBadRequest(w, err.Error())
			return
		} else if len(nameList) > 50 {
			serviceutil.SendBadRequest(w, "max limit 50")
			return
		}
	}
	//
	if token, err := ctl.getSteamToken(strutil.GetPathName(r.URL.Path)); nil != err || nil == token {
		serviceutil.SendBadRequest(w, err.Error())
	} else {
		// 校验 & 获取token
		tokenList := make(map[string]service.StreamTokenDto)
		parent := strutil.GetPathParent(token.FilePath)
		for i := 0; i < len(nameList); i++ {
			if filePath := strutil.Parse2UnixPath(parent + "/" + nameList[i]); !ctl.fm.IsFile(filePath) {
				serviceutil.SendBadRequest(w, ErrorFileNotExist.Error())
				return
			} else {
				if token, err := ctl.tt.AskReadToken(filePath, nil); nil == err {
					tokenList[nameList[i]] = *token.ToDto()
				}
			}
		}
		serviceutil.SendSuccess(w, tokenList)
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
			if list, err := ctl.fm.GetDirNodeList(prentPath, -1, -1); err != nil {
				serviceutil.SendServerError(w, err.Error())
			} else {
				res := make([]service.FNodeDto, 0)
				if len(list) > 0 {
					for i := 0; i < len(list); i++ {
						res = append(res, list[i].ToDto())
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
