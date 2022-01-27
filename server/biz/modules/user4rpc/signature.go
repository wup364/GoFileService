// Copyright (C) 2021 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 鉴权过滤器

package user4rpc

import (
	"fileservice/biz/service"
	"net/http"
	"os"
	"pakku/ipakku"
	"pakku/utils/logs"
	"pakku/utils/serviceutil"
	"pakku/utils/strutil"
	"sort"
	"strings"
)

// NewApiSignature NewApiSignature
func NewApiSignature(ch ipakku.AppCache) *Signature {
	vs := true
	if os.Getenv("validation-sign") == "false" {
		vs = false
	}
	return &Signature{validationSign: vs, ch: ch}
}

// Signature 签名校验
type Signature struct {
	validationSign bool
	ch             ipakku.AppCache `@autowired:"AppCache"`
}

// RestfulAPIFilter Api签名拦截器
func (st *Signature) RestfulAPIFilter(w http.ResponseWriter, r *http.Request) bool {
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
	// 从请求中获取accessKey, 不能为空
	accessKey := r.Header.Get(service.AuthHeader_AccessKey)
	// 填充Form对象
	if nil == r.Form {
		err := r.ParseForm()
		if nil != err {
			// 出现异常, 不继续处理
			serviceutil.SendServerError(w, err.Error())
			return false
		}
	}
	// 从请求中获取signature, 不能为空
	sign := r.Header.Get(service.AuthHeader_Sign)
	if len(sign) == 0 && len(r.Form[service.AuthHeader_Sign]) > 0 {
		sign = r.Form[service.AuthHeader_Sign][0]
	}
	//
	playload := ""
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" || contentType == "text/plain" || contentType == "application/javascript" || contentType == "application/xml" {
		playload = strutil.ReadAsString(r.Body)
		r.Body = &textReader{strings.NewReader(playload)}
	} else {
		playload = st.getFormPlayLoad(r)
	}
	// 校验参数合法性
	if err := st.VerificationSignature(r.URL.String(), sign, accessKey, playload); nil != err {
		serviceutil.SendErrorAndStatus(w, http.StatusUnauthorized, err.Error())
		return false
	}
	return true // 校验参数合法性-通过
}

// AskAccess 获取access
func (st *Signature) AskAccess(user service.UserInfo) (*service.UserAccessDto, error) {
	access := &service.UserAccess{
		UserID:    user.UserID,
		UserType:  user.UserType,
		UserName:  user.UserName,
		AccessKey: strutil.GetUUID(),
		SecretKey: strutil.GetUUID(),
		Props:     make(map[string]string),
	}
	if err := st.ch.Set(service.Cachelib_UserAccessToken, access.AccessKey, access); nil != err {
		return nil, err
	}
	return access.ToDto(), nil

}

// GetUserAccess 获取 useraccess
func (st *Signature) GetUserAccess(accessKey string) (*service.UserAccessDto, error) {
	var val service.UserAccess
	if err := st.ch.Get(service.Cachelib_UserAccessToken, accessKey, &val); nil == err {
		return val.ToDto(), nil
	} // else {
	// logs.Errorln(err)
	//}
	return nil, service.ErrorAuthentication
}

// GetAccessKey4Request 从http中获取accesskey
func (st *Signature) GetAccessKey4Request(r *http.Request) string {
	var accessKey string
	if accessKey = r.Header.Get(service.AuthHeader_AccessKey); len(accessKey) == 0 {
		if accessKey = r.FormValue(service.AuthHeader_AccessKey); len(accessKey) > 0 {
			r.Header.Set(service.AuthHeader_AccessKey, accessKey)
		}
	}
	if len(accessKey) == 0 {
		if ack, err := r.Cookie(service.AuthHeader_AccessKey); nil == err {
			if accessKey = ack.Value; len(accessKey) > 0 {
				r.Header.Set(service.AuthHeader_AccessKey, accessKey)
			}
		}
	}
	return accessKey
}

// DestroyAccess 销毁 useraccess
func (st *Signature) DestroyAccess(accessKey string) error {
	if err := st.ch.Del(service.Cachelib_UserAccessToken, accessKey); nil != err {
		if _, gerr := st.GetUserAccess(accessKey); nil == gerr {
			logs.Errorln(err)
			return err
		}
	}
	return nil
}

// RefreshAccessKey 刷新 access
func (st *Signature) RefreshAccessKey(accessKey string) error {
	var val service.UserAccess
	if err := st.ch.Get(service.Cachelib_UserAccessToken, accessKey, &val); nil == err {
		if err := st.ch.Set(service.Cachelib_UserAccessToken, val.AccessKey, val); nil != err {
			return err
		}
		return nil
	}
	return service.ErrorAuthentication
}

// VerificationSignature 验证auth
func (st *Signature) VerificationSignature(url, signature, accessKey, playload string) error {
	if access, err := st.GetUserAccess(accessKey); nil == err {
		if st.validationSign {
			// 校验参数完整性
			if sg, err := st.MakeSignature(url, access.AccessKey, access.SecretKey, playload); nil != err || sg != signature {
				return service.ErrorSignature
			}
		}
		return nil
	} else {
		return err
	}
}

// MakeSignature 请求签名一下
func (st *Signature) MakeSignature(url, accessKey, secretKey, playload string) (string, error) {
	return strutil.GetSHA256(url + ":" + accessKey + ":" + secretKey + ":" + playload), nil
}

// getFormPlayLoad 获取url传参的负载
func (st *Signature) getFormPlayLoad(r *http.Request) string {
	requestparameter := ""
	if nil != r.Form && len(r.Form) > 0 {
		keys := make([]string, 0) // 去掉参数为空的传值
		for key, val := range r.Form {
			if len(val) > 0 && key != service.AuthHeader_Sign && key != service.AuthHeader_AccessKey {
				keys = append(keys, key)
			}
		}

		_keysLen := len(keys)
		if _keysLen > 0 {
			sort.Strings(keys)
			for i := 0; i < len(keys); i++ {
				requestparameter += keys[i] + "=" + r.Form[keys[i]][0]
				if i < _keysLen-1 {
					requestparameter += "&"
				}
			}

		}
	}
	return requestparameter
}

// textReader textReader
type textReader struct {
	*strings.Reader
}

func (jr *textReader) Close() error {
	return nil
}
