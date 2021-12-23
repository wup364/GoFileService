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
	"errors"
	"fmt"
	"time"
)

const (
	// FPM_VisibleChild 子节点可见导致的父节点可见, 父节点其实并没有权限
	FPM_VisibleChild = iota
	// FPM_Visible 可见
	FPM_Visible
	// FPM_Read 只读
	FPM_Read
	// FPM_Write 可写
	FPM_Write
)

// ErrorConnIsNil 空连接
var ErrorConnIsNil = errors.New("the data source is empty")

// ErrorPermissionIDIsNil ErrorPermissionIDIsNil
var ErrorPermissionIDIsNil = errors.New("the permissionid is empty")

// ErrorPermissionPathIsNil ErrorPermissionPathIsNil
var ErrorPermissionPathIsNil = errors.New("the permission path is empty")

// ErrorPermissionIsNil ErrorPermissionIsNil
var ErrorPermissionIsNil = errors.New("the permission is empty")

// PermissionInfo 文件权限表存储的结构
type PermissionInfo struct {
	PermissionID string    // 权限路径
	Path         string    // 文件路径
	UserID       string    // 用户ID
	Permission   int64     // 权限值
	CtTime       time.Time // 插入时间
}

// PermissionInfoDto PermissionInfo传输对象
type PermissionInfoDto struct {
	PermissionID string    `json:"permissionID"` // 权限路径
	Path         string    `json:"path"`         // 文件路径
	UserID       string    `json:"userID"`       // 用户ID
	Permission   int64     `json:"permission"`   // 权限值
	CtTime       time.Time `json:"ctTime"`       // 插入时间
}

// Clone Clone
func (pi *PermissionInfo) Clone(val interface{}) error {
	if pit, ok := val.(*PermissionInfo); ok {
		pit.PermissionID = pi.PermissionID
		pit.Path = pi.Path
		pit.UserID = pi.UserID
		pit.Permission = pi.Permission
		pit.CtTime = pi.CtTime
		return nil
	}
	return fmt.Errorf("can't support clone %T ", val)
}

// ToDto 转传输对象
func (pi *PermissionInfo) ToDto() *PermissionInfoDto {
	return &PermissionInfoDto{
		PermissionID: pi.PermissionID,
		Path:         pi.Path,
		UserID:       pi.UserID,
		Permission:   pi.Permission,
		CtTime:       pi.CtTime,
	}
}

// FilePermission 文件权限管理接口
type FilePermission interface {
	FilePermissionCheck
	ListFPermissions() ([]PermissionInfo, error)                   // 列出所有权限数据, 无分页
	ListUserFPermissions(userID string) ([]PermissionInfo, error)  // 列出用户所有权限数据, 无分页
	QueryFPermission(permissionID string) (*PermissionInfo, error) // 根据权限ID查询详细信息
	AddFPermission(fpm PermissionInfo) error                       // 添加权限
	UpdateFPermission(fpm PermissionInfo) error                    // 修改权限
	DelFPermission(permissionID string) error                      // 根据permissionID删除权限
}

// FilePermissionCheck 文件权限校验
type FilePermissionCheck interface {
	GetUserPermissionSum(userID, path string) int64
	HashPermission(userID, path string, permission int64) bool
}
