// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 权限计算校验-暂时先粗暴实现

package filepermission

import (
	"fileservice/biz/service"
	"path"
	"strings"

	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/utypes"
)

// NewPermissionCheck NewPermissionCheck
func NewPermissionCheck(pmst *PermissionStory) *PermissionCheck {
	return &PermissionCheck{
		pmst: pmst,
		upms: utypes.NewSafeMap(),
	}
}

// PermissionCheck 权限校验
type PermissionCheck struct {
	pmst *PermissionStory
	upms *utypes.SafeMap
}

// InitPermission 初始化权限
func (pmc *PermissionCheck) InitPermission() error {
	pmc.initPms2Memory()
	return nil
}

// HashPermission 是否拥有某个权限
func (pmc *PermissionCheck) HashPermission(userID, path string, permission int64) bool {
	userPermission := pmc.GetUserPermissionSum(userID, path)
	if userPermission < service.FPM_VisibleChild {
		return false
	}
	// 子节点可见导致的父节点可见, 父节点其实并没有权限
	if permission == service.FPM_VisibleChild && userPermission >= service.FPM_VisibleChild {
		return true
	}
	return 1<<permission == userPermission&(1<<permission)
}

// GetUserPermissionSum 获取用户对路径拥有的权限数值
// -1 为无权限 > -1 为有权限
func (pmc *PermissionCheck) GetUserPermissionSum(userID, path string) int64 {
	var userPms *utypes.SafeMap
	if userPms = pmc.getUserPermission(userID); userPms == nil || userPms.Size() == 0 {
		return -1
	}
	permissionSum := int64(-1)
	if val, ok := userPms.Get(path); !ok {
		// 可能只是没有权限记录, 或者是个文件, 我们则需要尝试上级目录是否拥有>VISIBLECHILD的权限
		parent := path
		for {
			if parent = parent[:strings.LastIndex(parent, "/")]; len(parent) == 0 {
				parent = "/"
			}
			// 如果从上级找到了权限, 则需要是>VISIBLECHILD的权限
			// 组最近的权限设定, 即便上级再上级有权限也不管
			if val, ok := userPms.Get(parent); ok && val.(int64) > service.FPM_VisibleChild {
				return val.(int64)
			}
			// 如果没有找到, 则继续网上找, 如果到/还没有, 则无权限
			if parent == "/" {
				return -1
			}
		}
	} else {
		permissionSum = val.(int64)
	}
	// 是个文件夹&有记录
	return permissionSum
}

// initPms2Memory 加载权限结构到内存
func (pmc *PermissionCheck) initPms2Memory() {
	list, err := pmc.pmst.ListFPermissions()
	if nil != err {
		logs.Panicln(err)
	}
	pmc.upms.Clear()
	// 初始化数据
	for i := 0; i < len(list); i++ {
		var userpms *utypes.SafeMap
		if val, ok := pmc.upms.Get(list[i].UserID); !ok {
			userpms = utypes.NewSafeMap()
			pmc.upms.Put(list[i].UserID, userpms)
		} else {
			userpms = val.(*utypes.SafeMap)
		}
		// 当前用户&当前路径
		list[i].Path = path.Clean(list[i].Path)
		userpms.Put(list[i].Path, list[i].Permission)
		// 检查上级目录, 生成 VISIBLECHILD 权限
		parent := list[i].Path
		for {
			if parent = parent[:strings.LastIndex(parent, "/")]; len(parent) == 0 {
				if _, ok := userpms.Get("/"); !ok {
					userpms.Put("/", int64(service.FPM_VisibleChild))
				}
				break
			}

			if _, ok := userpms.Get(parent); !ok {
				userpms.Put(parent, int64(service.FPM_VisibleChild))
			}
		}
	}
}

// getUserPermission 获取用户权限结构
func (pmc *PermissionCheck) getUserPermission(userID string) *utypes.SafeMap {
	if val, ok := pmc.upms.Get(userID); ok {
		return val.(*utypes.SafeMap)
	}
	return nil
}
