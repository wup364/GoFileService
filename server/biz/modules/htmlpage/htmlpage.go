// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 静态资源加载器 配置参数(mloader.GetParam): htmlpage.static

package htmlpage

import (
	"fileservice/biz/modules/user4rpc"
	"fileservice/biz/service"
	"pakku/ipakku"
	"pakku/utils/logs"
	"strings"

	"net/http"
)

// HTMLPage 静态资源加载器
type HTMLPage struct {
	static string
	sg     *user4rpc.Signature
	ch     ipakku.AppCache   `@autowired:"AppCache"`
	sv     ipakku.AppService `@autowired:"AppService"`
}

// AsModule 模块加载器接口实现, 返回模块信息&配置
func (html *HTMLPage) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "HTMLPage",
		Version:     1.0,
		Description: "静态资源",
		OnReady: func(mctx ipakku.Loader) {
			html.static = mctx.GetParam("htmlpage.static").ToString("./webapps")
		},
		OnInit: func() {
			html.sg = user4rpc.NewApiSignature(html.ch)
			html.registerServlet()
		},
	}
}

// registerServlet 注册servlet
func (html *HTMLPage) registerServlet() {
	// 1.首页重定向
	// if err := html.sv.Any("/", PageDispatch{}.Index); err != nil {
	// 	logs.Panicln(err)
	// }
	// 2. 页面资源
	if err := html.sv.SetStaticDIR("/", html.static, html.staticFilter); nil != err {
		logs.Panicln(err)
	}
}

// staticFilter 静态资源过滤器
func (html *HTMLPage) staticFilter(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if acrhs := r.Header["Access-Control-Request-Headers"]; len(acrhs) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(acrhs, ","))
	}
	if acrms := r.Header["Access-Control-Request-Method"]; len(acrms) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(acrms, ","))
	}
	if strings.ToUpper(r.Method) == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return false
	}
	// fmt.Println("Static Path: ", r.URL.Path)
	if strings.HasSuffix(r.URL.Path, ".html") {
		if ack, err := r.Cookie(service.AuthHeader_AccessKey); nil == err && len(ack.Value) > 0 {
			// 是否只要 cookie 里面有合法 accessKey 就可以了
			if _, err := html.sg.GetUserAccess(ack.Value); nil == err {
				return true
			}
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return false
	}
	return true
}
