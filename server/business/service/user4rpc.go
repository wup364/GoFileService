// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wup364/pakku/ipakku"
)

const (
	// UserType_Admin 用户类型
	UserType_Admin = 0
	// UserType_Normal 用户类型
	UserType_Normal = 1
	// AuthHeader_AccessKey 用于验证签名的key
	AuthHeader_AccessKey = "X-Ack"
	// AuthHeader_Sign 客户端签名结果的key
	AuthHeader_Sign = "X-Sign"
	// Cachelib_UserAccessToken 用户access信息缓存库
	Cachelib_UserAccessToken = "User4RPC:AccessToken"
)

// ErrorUserIDIsNil ErrorUserIDIsNil
var ErrorUserIDIsNil = errors.New("the userID is empty")

// ErrorUserNotExist ErrorUserNotExist
var ErrorUserNotExist = errors.New("user does not exist")

// ErrorUserNameIsNil ErrorUserNameIsNil
var ErrorUserNameIsNil = errors.New("the userName is empty")

// ErrorAuthentication 认证失败
var ErrorAuthentication = errors.New("authentication failed")

// ErrorSignature 签名错误
var ErrorSignature = errors.New("request content signature error")

// User4RPC 用户管理接口
type User4RPC interface {
	UserManage
	UserAuth4Rpc
}

// UserManage access接口
type UserManage interface {
	Clear() error
	AddUser(user *UserInfo) error                  // 添加用户
	DelUser(userID string) error                   // 根据userID删除用户
	CheckPwd(userID, pwd string) bool              // 校验密码是否一致
	UpdatePWD(userID, pwd string) error            // 修改用户密码
	ListAllUsers() ([]UserInfoDto, error)          // 列出所有用户数据, 无分页
	QueryUser(userID string) (*UserInfoDto, error) // 根据用户ID查询详细信息
	UpdateUserName(userID, userName string) error  // 修改用户名字
}

// UserAuth4Rpc access接口
type UserAuth4Rpc interface {
	// GetAuthFilterFunc 获取过滤器实现
	GetAuthFilterFunc() ipakku.FilterFunc
	// AskAccess 获取access
	AskAccess(userID, pwd string) (*UserAccessDto, error)
	// GetSecretKey 获取 userAccess
	GetUserAccess(accessKey string) (*UserAccessDto, error)
	// GetAccessKey4Request 从http中获取accesskey
	GetAccessKey4Request(r *http.Request) string
	// RefreshAccessKey 刷新 access
	RefreshAccessKey(accessKey string) error
	// DestroyAccess 销毁 access
	DestroyAccess(accessKey string) error
}

// UserInfo 用户表存储的结构
type UserInfo struct {
	UserType int
	UserID   string
	UserName string
	UserPWD  string
	CtTime   time.Time
}

// UserAccess access内容
type UserAccess struct {
	UserID    string
	UserType  int
	UserName  string
	AccessKey string
	SecretKey string
	Props     map[string]string
}

// UserInfoDto UserInfo传输对象
type UserInfoDto struct {
	UserType int       `json:"userType"`
	UserID   string    `json:"userID"`
	UserName string    `json:"userName"`
	CtTime   time.Time `json:"ctTime"`
}

// UserAccessDto UserAccess传输对象
type UserAccessDto struct {
	UserType  int    `json:"userType"`
	UserID    string `json:"userID"`
	UserName  string `json:"userName"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// Clone Clone
func (ua *UserAccess) Clone(val interface{}) error {
	if uat, ok := val.(*UserAccess); ok {
		uat.UserID = ua.UserID
		uat.UserType = ua.UserType
		uat.UserName = ua.UserName
		uat.AccessKey = ua.AccessKey
		uat.SecretKey = ua.SecretKey
		uat.Props = ua.Props
		return nil
	}
	return fmt.Errorf("can't support clone %T ", val)
}

// ToDto 转传输对象
func (ua *UserAccess) ToDto() *UserAccessDto {
	return &UserAccessDto{
		UserType:  ua.UserType,
		UserID:    ua.UserID,
		UserName:  ua.UserName,
		AccessKey: ua.AccessKey,
		SecretKey: ua.SecretKey,
	}
}

// ToDto 转传输对象
func (ui *UserInfo) ToDto() *UserInfoDto {
	return &UserInfoDto{
		UserType: ui.UserType,
		UserID:   ui.UserID,
		UserName: ui.UserName,
		CtTime:   ui.CtTime,
	}
}

// ToJSON 转传JSON
func (ua *UserAccessDto) ToJSON() string {
	if bt, err := json.Marshal(ua); nil != err {
		return ""
	} else {
		return string(bt)
	}
}

// ToJSON 转传JSON
func (ui *UserInfoDto) ToJSON() string {
	if bt, err := json.Marshal(ui); nil != err {
		return ""
	} else {
		return string(bt)
	}
}
