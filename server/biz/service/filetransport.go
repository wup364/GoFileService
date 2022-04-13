// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件操作接口

package service

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	StreamTokenType_Read     = 0
	StreamTokenType_Write    = 1
	CacheLib_StreamToken     = "FileStransport:Token"
	CacheLib_StreamToken_Exp = 60 * 30
)

// ErrInvalidToken 无效的token
var ErrInvalidToken = errors.New("invalid token")

// StreamTokenType token类型
type StreamTokenType int

// TransportToken 传输token控制
type TransportToken interface {
	AskWriteToken(src string, props map[string]string) (*StreamToken, error)
	AskReadToken(src string, props map[string]string) (*StreamToken, error)
	QueryToken(token string) (*StreamToken, error)
	RefreshToken(token string) (st *StreamToken, err error)
	DestroyToken(token string, override bool) (err error)
}

// StreamToken token信息
type StreamToken struct {
	Token    string
	TokenURL string
	FilePath string
	CTime    int64
	MTime    int64
	Type     StreamTokenType
}

// StreamTokenDto token信息
type StreamTokenDto struct {
	Token    string `json:"token"`
	TokenURL string `json:"tokenURL"`
	FilePath string `json:"filePath"`
	CTime    int64  `json:"cTime"`
}

// Clone Clone
func (ua *StreamToken) Clone(val interface{}) error {
	if st, ok := val.(*StreamToken); ok {
		st.Token = ua.Token
		st.TokenURL = ua.TokenURL
		st.FilePath = ua.FilePath
		st.CTime = ua.CTime
		st.MTime = ua.MTime
		st.Type = ua.Type
		return nil
	}
	return fmt.Errorf("can't support clone %T ", val)
}

// ToDto 转传输对象
func (s *StreamToken) ToDto() *StreamTokenDto {
	return &StreamTokenDto{
		Token:    s.Token,
		TokenURL: s.TokenURL,
		FilePath: s.FilePath,
		CTime:    s.CTime,
	}
}

// ToJSON 转传JSON
func (dto *StreamTokenDto) ToJSON() string {
	if bt, err := json.Marshal(dto); nil != err {
		return ""
	} else {
		return string(bt)
	}
}
