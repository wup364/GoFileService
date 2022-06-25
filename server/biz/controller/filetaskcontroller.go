// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 异步任务接口

package controller

import (
	"fileservice/biz/service"
	"net/http"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/serviceutil"
)

// AsyncTaskCtrl 异步处理任务接口
type AsyncTaskCtrl struct {
	um  service.UserAuth4Rpc `@autowired:"User4RPC"`
	ast service.AsyncTask    `@autowired:"AsyncTask"`
}

// AsController 实现 AsController 接口
func (ctl *AsyncTaskCtrl) AsController() ipakku.ControllerConfig {
	return ipakku.ControllerConfig{
		RequestMapping: "/filetask/v1",
		RouterConfig: ipakku.RouterConfig{
			ToLowerCase: true,
			HandlerFunc: [][]interface{}{
				{http.MethodPost, ctl.AsyncExec},
				{http.MethodPost, ctl.AsyncExecToken},
			},
		},
		FilterConfig: ipakku.FilterConfig{
			FilterFunc: [][]interface{}{
				{`/:[\s\S]*`, ctl.um.GetAuthFilterFunc()},
			},
		},
	}
}

// AsyncExec 发起一个异步操作, 返回一个可以查询的tooken
func (ctl *AsyncTaskCtrl) AsyncExec(w http.ResponseWriter, r *http.Request) {
	if executor, err := ctl.ast.GetTaskObject(r.FormValue("func")); nil != err {
		serviceutil.SendServerError(w, err.Error())
	} else {
		if token, err := executor.Execute(r); nil != err {
			serviceutil.SendServerError(w, err.Error())
		} else {
			serviceutil.SendSuccess(w, token)
		}
	}
}

// AsyncExecToken 查询由AsyncExec返回的token状态
func (ctl *AsyncTaskCtrl) AsyncExecToken(w http.ResponseWriter, r *http.Request) {
	if executor, err := ctl.ast.GetTaskObject(r.FormValue("func")); nil != err {
		serviceutil.SendServerError(w, err.Error())
	} else {
		executor.Status(w, r)
	}
}
