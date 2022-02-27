// Copyright (C) 2021 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

//  文件传输token控制-先支持仅本系统模式

package serviceimpl

import (
	"fileservice/biz/service"
	"pakku/ipakku"
	"pakku/utils/fileutil"
	"pakku/utils/logs"
)

// TransportToken 文件传输token
type TransportToken struct {
	f service.FileDatas `@autowired:"FileDatas"`
	c ipakku.AppCache   `@autowired:"AppCache"`
}

// Pakku 模块加载器接口实现, 返回模块信息&配置
func (n *TransportToken) Pakku() ipakku.Opts {
	return ipakku.Opts{
		Name:        "TransportToken",
		Version:     1.0,
		Description: "文件传输token控制",
		OnInit: func() {
			if err := n.c.RegLib(service.CacheLib_StreamToken, service.CacheLib_StreamToken_Exp); nil != err {
				logs.Panicln(err)
			}
		},
	}
}

// AskWriteToken 申请写入token
func (n *TransportToken) AskWriteToken(src string, props map[string]interface{}) (*service.StreamToken, error) {
	if token, err := n.f.DoAskAccessToken(src, service.AccessTokenType(service.StreamTokenType_Write), props); nil != err {
		return nil, err
	} else {
		st := &service.StreamToken{
			FilePath: src,
			Token:    token.Token,
			CTime:    token.CTime,
			TokenURL: token.TokenURL,
			FsysType: token.DriverType,
			Type:     service.StreamTokenType_Write,
		}
		if err := n.c.Set(service.CacheLib_StreamToken, token.Token, st); nil == err {
			return st, nil
		} else {
			return nil, err
		}
	}
}

// AskReadToken 申请读取token
func (n *TransportToken) AskReadToken(src string, props map[string]interface{}) (*service.StreamToken, error) {
	var node *service.FNode
	if node = n.f.GetNode(src); nil == node || !node.IsFile {
		return nil, fileutil.PathNotExist("askReadToken", src)
	}
	if token, err := n.f.DoAskAccessToken(src, service.AccessTokenType(service.StreamTokenType_Read), props); nil != err {
		return nil, err
	} else {
		st := &service.StreamToken{
			FilePath: src,
			Token:    token.Token,
			CTime:    token.CTime,
			TokenURL: token.TokenURL,
			FsysType: token.DriverType,
			Type:     service.StreamTokenType_Read,
		}
		if err := n.c.Set(service.CacheLib_StreamToken, token.Token, st); nil == err {
			return st, nil
		} else {
			return nil, err
		}
	}
}

// SubmitToken 提交token
func (n *TransportToken) SubmitToken(token string, props map[string]interface{}) (*service.FNode, error) {
	if tokenObj, err := n.QueryToken(token); nil == err {
		return n.f.DoSubmitToken(service.AccessToken{
			Path:     tokenObj.FilePath,
			Token:    tokenObj.Token,
			CTime:    tokenObj.CTime,
			TokenURL: tokenObj.TokenURL,
		}, props)
	} else {
		return nil, err
	}
}

// QueryToken 查出token
func (n *TransportToken) QueryToken(token string) (*service.StreamToken, error) {
	var val service.StreamToken
	if err := n.c.Get(service.CacheLib_StreamToken, token, &val); nil == err {
		return &val, nil
	} else if err == ipakku.ErrNoCacheHit {
		return nil, service.ErrInvalidToken
	} else {
		return nil, err
	}

}

// RefreshToken RefreshToken
func (n *TransportToken) RefreshToken(token string) (st *service.StreamToken, err error) {
	if st, err = n.QueryToken(token); nil == err {
		if err = n.c.Set(service.CacheLib_StreamToken, token, st); nil != err {
			logs.Errorln(err)
		}
	} else {
		return nil, err
	}
	return n.QueryToken(token)
}

// DestroyToken 销毁token
func (n *TransportToken) DestroyToken(token string, override bool) (err error) {
	if _, err = n.QueryToken(token); nil == err {
		if err = n.c.Del(service.CacheLib_StreamToken, token); nil != err {
			if _, err := n.QueryToken(token); err == service.ErrInvalidToken {
				return nil
			}
		}
	} else if err == service.ErrInvalidToken {
		err = nil
	}
	return err
}
