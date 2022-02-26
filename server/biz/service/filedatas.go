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
	"io"
)

// AccessTokenType AccessTokenType
type AccessTokenType int

// FileDatas 文件数据块管理模块
type FileDatas interface {
	IsDir(src string) bool
	IsFile(src string) bool
	IsExist(src string) bool

	GetNode(src string) *FNode
	GetFileSize(src string) int64
	GetDirList(src string, limit int, offset int) []string
	GetNodes(src []string, ignoreNotIsExist bool) ([]FNode, error)
	GetDirNodeList(src string, limit int, offset int) ([]FNode, error)

	DoMkDir(src string) error
	DoDelete(src string) error
	DoRename(src string, dst string) error
	DoCopy(src, dst string, replace bool) error
	DoMove(src, dst string, replace bool) error

	DoWrite(src string, ioReader io.Reader) error
	DoRead(src string, offset int64) (io.ReadCloser, error)
	//
	DoAskAccessToken(src string, tokenType AccessTokenType, props map[string]interface{}) (*AccessToken, error)
}

// FNode 文件|夹基础属性(filedatas)
type FNode struct {
	Path   string
	Mtime  int64
	IsFile bool
	IsDir  bool
	Size   int64
}

// ToDto 转传输对象
func (f *FNode) ToDto() FNodeDto {
	return FNodeDto{
		Path:   f.Path,
		Mtime:  f.Mtime,
		IsFile: f.IsFile,
		IsDir:  f.IsDir,
		Size:   f.Size,
	}
}

// FNodeDto 传输对象
type FNodeDto struct {
	Path   string `json:"path"`
	Mtime  int64  `json:"mtime"`
	IsFile bool   `json:"isFile"`
	IsDir  bool   `json:"isDir"`
	Size   int64  `json:"size"`
}

// ToJSON 转传JSON
func (dto *FNodeDto) ToJSON() string {
	if bt, err := json.Marshal(dto); nil != err {
		return ""
	} else {
		return string(bt)
	}
}

// AccessToken token信息
type AccessToken struct {
	Token    string
	TokenURL string
	CTime    int64
}
