// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件传输接口

package controller

import (
	"fileservice/biz/service"
	"io"
	"mime/multipart"
	"net/http"
	"pakku/ipakku"
	"pakku/utils/logs"
	"pakku/utils/serviceutil"
	"pakku/utils/strutil"
	"strconv"
	"strings"
)

// TransportCtrl 文件传输接口
type TransportCtrl struct {
	fm  service.FileDatas           `@autowired:"FileDatas"`
	um  service.UserAuth4Rpc        `@autowired:"User4RPC"`
	tt  service.TransportToken      `@autowired:"TransportToken"`
	pms service.FilePermissionCheck `@autowired:"FilePermission"`
}

// AsController 实现 AsController 接口
func (ctl *TransportCtrl) AsController() ipakku.ControllerConfig {
	return ipakku.ControllerConfig{
		RequestMapping: "/filestream/v1",
		RouterConfig: ipakku.RouterConfig{
			ToLowerCase: true,
			HandlerFunc: [][]interface{}{
				{"GET", "read/:" + `[\s\S]*`, ctl.Read},
				{"ANY", "put/:" + `[\s\S]*`, ctl.Put},
				{"GET", ctl.Token},
			},
		},
		FilterConfig: ipakku.FilterConfig{
			FilterFunc: [][]interface{}{
				{`/:[\s\S]*`, func(rw http.ResponseWriter, r *http.Request) bool {
					if strings.Contains(r.URL.Path, "/read/") || strings.Contains(r.URL.Path, "/put/") {
						return true
					}
					return ctl.um.GetAuthFilterFunc()(rw, r)
				}},
			},
		},
	}
}

// checkPermision 检查权限
func (ctl *TransportCtrl) checkPermision(userID, path string, permission int64) bool {
	return ctl.pms.HashPermission(userID, path, permission)
}

// getUserID4Request 获取登录用户
func (ctl *TransportCtrl) getUserID4Request(r *http.Request) string {
	if askstr := ctl.um.GetAccessKey4Request(r); len(askstr) > 0 {
		if ack, err := ctl.um.GetUserAccess(askstr); nil == err {
			return ack.UserID
		}
	}
	return ""
}

// Token 传输令牌申请
func (ctl *TransportCtrl) Token(w http.ResponseWriter, r *http.Request) {
	qdata := r.FormValue("data")
	qtype := r.FormValue("type")
	if len(qtype) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var err error
	var token *service.StreamToken
	//
	if qtype == "download" {
		if !ctl.checkPermision(ctl.getUserID4Request(r), qdata, service.FPM_Read) {
			err = ErrorPermissionInsufficient
		} else {
			token, err = ctl.tt.AskReadToken(qdata)
		}
	} else if qtype == "upload" {
		if !ctl.checkPermision(ctl.getUserID4Request(r), qdata, service.FPM_Write) {
			err = ErrorPermissionInsufficient
		} else {
			token, err = ctl.tt.AskWriteToken(qdata)
		}
	} else {
		err = ErrorNotSupport
	}
	if nil != token {
		serviceutil.SendSuccess(w, token.ToDto())
	} else if nil != err {
		serviceutil.SendBadRequest(w, err.Error())
	}
}

// Put 文件上传, 支持Form和Body上传方式
func (ctl *TransportCtrl) Put(w http.ResponseWriter, r *http.Request) {
	if reqMethod := strings.ToLower(r.Method); reqMethod != "post" && reqMethod != "put" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	token, err := ctl.getSteamToken(strutil.GetPathName(r.URL.Path))
	if nil != err || nil == token {
		serviceutil.SendServerError(w, err.Error())
		return
	}
	if len(token.FilePath) == 0 {
		serviceutil.SendServerError(w, "path is empty")
		return
	}
	//
	filename := strutil.GetPathName(token.FilePath)
	if mr, err := r.MultipartReader(); err == nil {
		var pName string
		if pName = r.Header.Get(headerFormNameFile); len(pName) == 0 {
			pName = filename
		}
		hasfile := false
		for {
			var p *multipart.Part
			if p, err = mr.NextPart(); nil == p || err == io.EOF {
				break
			}
			if p.FormName() != pName {
				continue
			}
			hasfile = true
			if err = ctl.fm.DoWrite(token.FilePath, p); nil != err {
				serviceutil.SendServerError(w, err.Error())
			} else {
				serviceutil.SendSuccess(w, "")
			}
			break
		}
		if !hasfile {
			serviceutil.SendServerError(w, "file not found from the form")
		}
	} else if nil != err && err == http.ErrNotMultipart {
		if err := ctl.fm.DoWrite(token.FilePath, r.Body); nil != err {
			serviceutil.SendServerError(w, err.Error())
		} else {
			serviceutil.SendSuccess(w, "")
		}
	} else {
		serviceutil.SendServerError(w, err.Error())
	}
}

// Read 下载
func (ctl *TransportCtrl) Read(w http.ResponseWriter, r *http.Request) {
	qToken := strutil.GetPathName(r.URL.Path)
	token, err := ctl.getSteamToken(qToken)
	if nil != err || nil == token {
		serviceutil.SendServerError(w, err.Error())
		return
	} else {
		// 校验
		if !ctl.fm.IsFile(token.FilePath) {
			serviceutil.SendBadRequest(w, ErrorFileNotExist.Error())
			return
		} else {
			// 刷新token, 使其不过期
			ctl.tt.RefreshToken(qToken)
		}
	}
	// 默认设置下载头
	if r.FormValue("action") != "stream" {
		w.Header().Set("Content-Disposition", "attachment; filename="+strutil.GetPathName(token.FilePath))
		w.Header().Set("Content-Type", strutil.GetMimeTypeBySuffix(strutil.GetPathSuffix(token.FilePath)))
	}
	//
	maxSize := ctl.fm.GetFileSize(token.FilePath)
	start, end := serviceutil.GetRequestRange(r, maxSize)
	//
	if fr, err := ctl.fm.DoRead(token.FilePath, start); nil != err {
		serviceutil.SendServerError(w, err.Error())
	} else {
		defer fr.Close()
		ctLength := end - start
		if ctLength == 0 || ctLength < 0 {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.Header().Set("Content-Length", strconv.Itoa(int(ctLength)))
			if ctLength != maxSize {
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

// getSteamToken 获取文件传输Token对象
func (ctl *TransportCtrl) getSteamToken(token string) (*service.StreamToken, error) {
	if tokenBody, err := ctl.tt.QueryToken(token); nil != err || nil == tokenBody {
		return nil, ErrorOprationExpires
	} else {
		return tokenBody, nil
	}
}
