// Copyright (C) 2021 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// token控制
package asynctask

import (
	"fileservice/biz/service"
	"pakku/ipakku"
	"pakku/utils/strutil"
)

// NewTaskToken Task token控制
func NewTaskToken(c ipakku.AppCache) *TaskToken {
	return &TaskToken{c: c}
}

// TaskToken Task token
type TaskToken struct {
	c ipakku.AppCache `@autowired:"AppCache"`
}

// AskToken 申请token
func (n *TaskToken) AskToken(val interface{}) (string, error) {
	token := strutil.GetUUID()
	if err := n.c.Set(service.CacheLib_AsyncTaskToken, token, val); nil == err {
		return token, nil
	} else {
		return "", err
	}
}

// QueryToken 查出token
func (n *TaskToken) QueryToken(token string, val interface{}) error {
	if err := n.c.Get(service.CacheLib_AsyncTaskToken, token, val); nil == err {
		return nil
	} else if err == ipakku.ErrNoCacheHit {
		return service.ErrorOprationExpires
	} else {
		return err
	}

}

// RefreshToken 重新设定值
func (n *TaskToken) RefreshToken(token string, val interface{}) (err error) {
	return n.c.Set(service.CacheLib_AsyncTaskToken, token, val)
}

// DestroyToken 销毁token
func (n *TaskToken) DestroyToken(token string, override bool) (err error) {
	if err = n.c.Del(service.CacheLib_AsyncTaskToken, token); nil != err && err == ipakku.ErrNoCacheHit {
		err = nil
	}
	return err
}
