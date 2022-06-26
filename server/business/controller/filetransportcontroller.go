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
	"fileservice/business/modules/filetransport"
	"fileservice/business/service"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/httpclient"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/serviceutil"
	"github.com/wup364/pakku/utils/strutil"
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
				{http.MethodHead, "read/:[A-Za-z0-9]+$", ctl.ReadHead},
				{http.MethodGet, "read/:[A-Za-z0-9]+$", ctl.Read},
				{http.MethodPost, "put/:[A-Za-z0-9]+$", ctl.Put},
				{http.MethodPut, "put/:[A-Za-z0-9]+$", ctl.Put},
				{http.MethodGet, "token", ctl.GetToken},
				{http.MethodPost, "token", ctl.PostToken},
			},
		},
		FilterConfig: ipakku.FilterConfig{
			FilterFunc: [][]interface{}{
				{`/:[\s\S]*`, func(rw http.ResponseWriter, r *http.Request) bool {
					if strings.Contains(r.URL.Path, "/read/") || strings.Contains(r.URL.Path, "/put/") {
						return true
					} else if strings.Contains(r.URL.Path, "/token") && r.FormValue("type") == "submitupload" {
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

// GetToken 传输令牌申请
func (ctl *TransportCtrl) GetToken(w http.ResponseWriter, r *http.Request) {
	qdata := strutil.Parse2UnixPath(r.FormValue("data"))
	qtype := r.FormValue("type")
	if len(qtype) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var err error
	var token *service.StreamToken
	//
	if qtype == "stream" {
		if !ctl.checkPermision(ctl.getUserID4Request(r), qdata, service.FPM_Read) {
			err = ErrorPermissionInsufficient
		} else {
			if token, err = ctl.tt.AskReadToken(qdata, nil); nil == err {
				token.TokenURL = "/filestream/v1/read/" + token.Token
			}
		}
	} else if qtype == "download" {
		if !ctl.checkPermision(ctl.getUserID4Request(r), qdata, service.FPM_Read) {
			err = ErrorPermissionInsufficient
		} else {
			if token, err = ctl.tt.AskReadToken(qdata, nil); nil == err {
				token.TokenURL = httpclient.BuildURLWithArray("/filestream/v1/read/"+token.Token, [][]string{{"name", strutil.GetPathName(qdata)}})
			}
		}
	} else if qtype == "upload" {
		if !ctl.checkPermision(ctl.getUserID4Request(r), qdata, service.FPM_Write) {
			err = ErrorPermissionInsufficient
		} else {
			if token, err = ctl.tt.AskWriteToken(qdata, nil); nil == err {
				token.TokenURL = "/filestream/v1/put/" + token.Token
			}
		}
	} else {
		err = ErrorNotSupport
	}
	if nil != err {
		serviceutil.SendBadRequest(w, err.Error())
	} else {
		serviceutil.SendSuccess(w, token.ToDto())
	}
}

// PostToken 传输令牌操作
func (ctl *TransportCtrl) PostToken(w http.ResponseWriter, r *http.Request) {
	qtype := r.FormValue("type")
	if len(qtype) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var result interface{}
	var err error
	if qtype == "submitupload" {
		err = ErrorNotSupport
	} else {
		err = ErrorNotSupport
	}
	if nil != err {
		serviceutil.SendBadRequest(w, err.Error())
	} else {
		serviceutil.SendSuccess(w, result)
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
		serviceutil.SendBadRequest(w, err.Error())
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
				p.Close()
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
func (ctl *TransportCtrl) ReadHead(w http.ResponseWriter, r *http.Request) {
	if token, err := ctl.getSteamToken(strutil.GetPathName(r.URL.Path)); nil != err || nil == token {
		serviceutil.SendBadRequest(w, err.Error())
		return
	} else {
		maxSize := ctl.fm.GetFileSize(token.FilePath)
		start, end, hasRange := serviceutil.GetRequestRange(r, maxSize)
		w.Header().Set("Content-Length", strconv.Itoa(int(end-start)))
		if hasRange {
			w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(start))+"-"+strconv.Itoa(int(end-1))+"/"+strconv.Itoa(int(maxSize)))
		}
	}

}

// Read 下载
func (ctl *TransportCtrl) Read(w http.ResponseWriter, r *http.Request) {
	qToken := strutil.GetPathName(r.URL.Path)
	token, err := ctl.getSteamToken(qToken)
	if nil != err || nil == token {
		serviceutil.SendBadRequest(w, err.Error())
		return
	}
	// 校验
	if !ctl.fm.IsFile(token.FilePath) {
		serviceutil.SendBadRequest(w, ErrorFileNotExist.Error())
		return
	}
	// name
	if name := r.FormValue("name"); len(name) > 0 {
		w.Header().Set("Content-Type", strutil.GetMimeTypeBySuffix(name))
		w.Header().Set("Content-Disposition", "attachment; filename="+name)
	}
	//
	maxSize := ctl.fm.GetFileSize(token.FilePath)
	start, end, hasRange := serviceutil.GetRequestRange(r, maxSize)
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
			if hasRange {
				w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(start))+"-"+strconv.Itoa(int(end-1))+"/"+strconv.Itoa(int(maxSize)))
				w.WriteHeader(http.StatusPartialContent)
			}
			if _, err := io.Copy(w, filetransport.NewTokenReaderWarp(token, io.LimitReader(fr, ctLength), ctl.tt)); nil != err && err != io.EOF {
				if opErr, ok := err.(*net.OpError); ok && nil != opErr.Err {
					if opErr, ok = opErr.Err.(*net.OpError); ok && opErr.Op == "write" {
						return
					}
				}
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
