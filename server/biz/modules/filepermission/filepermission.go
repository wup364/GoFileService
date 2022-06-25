// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件权限管理模块

package filepermission

import (
	"errors"
	"fileservice/biz/constants"
	"fileservice/biz/service"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/fileutil"
	"github.com/wup364/pakku/utils/logs"
)

// FilePermission 用户管理模块
type FilePermission struct {
	pmc  *PermissionCheck
	pms  *PermissionStory
	conf ipakku.AppConfig `@autowired:"AppConfig"`
}

// AsModule 模块加载器接口实现, 返回模块信息&配置
func (fpms *FilePermission) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "FilePermission",
		Version:     1.0,
		Description: "文件权限模块",
		OnReady: func(mctx ipakku.Loader) {
			deftDataSource := "./.datas/" + mctx.GetParam(ipakku.PARAMKEY_APPNAME).ToString("app") + ".db?cache=shared"
			confDataSource := fpms.conf.GetConfig("filepermission.sotre.datasource").ToString(deftDataSource)
			if confDataSource == deftDataSource {
				fpms.mkSqliteDIR() // 创建sqlite文件存放目录
			}
			// PermissionStory
			fpms.pms = new(PermissionStory)
			if err := fpms.pms.Initial(constants.DBSetting{
				DriverName:     fpms.conf.GetConfig("filepermission.sotre.driver").ToString("sqlite3"),
				DataSourceName: confDataSource,
			}); nil != err {
				logs.Panicln(err)
			}
		},
		OnSetup: func() {
			// 执行建库、建表
			if err := fpms.pms.Install(); nil != err {
				logs.Panicln(err)
			}
		},
		OnInit: func() {
			// PermissionCheck
			fpms.pmc = NewPermissionCheck(fpms.pms)
			if err := fpms.pmc.InitPermission(); nil != err {
				logs.Panicln(err)
			}
		},
	}
}

// HashPermission 是否拥有某个权限
func (fpms *FilePermission) HashPermission(userID, path string, permission int64) bool {
	return fpms.pmc.HashPermission(userID, path, permission)
}

// GetUserPermissionSum 获取用户对路径拥有的权限数值
// -1 为无权限 > -1 为有权限
func (fpms *FilePermission) GetUserPermissionSum(userID, path string) int64 {
	return fpms.pmc.GetUserPermissionSum(userID, path)
}

// ListFPermissions 列出所有权限数据, 无分页
func (fpms *FilePermission) ListFPermissions() ([]service.PermissionInfo, error) {
	return fpms.pms.ListFPermissions()
}

// ListUserFPermissions 列出用户所有权限数据, 无分页
func (fpms *FilePermission) ListUserFPermissions(userID string) ([]service.PermissionInfo, error) {
	return fpms.pms.ListUserFPermissions(userID)
}

// QueryFPermission 根据权限ID查询详细信息
func (fpms *FilePermission) QueryFPermission(permissionID string) (*service.PermissionInfo, error) {
	return fpms.pms.QueryFPermission(permissionID)
}

// AddFPermission 添加权限
func (fpms *FilePermission) AddFPermission(fpm service.PermissionInfo) error {
	err := fpms.pms.AddFPermission(fpm)
	fpms.pmc.initPms2Memory()
	return err
}

// UpdateFPermission 修改权限
func (fpms *FilePermission) UpdateFPermission(fpm service.PermissionInfo) error {
	pms, err := fpms.pms.QueryFPermission(fpm.PermissionID)
	if nil != err {
		return err
	}
	if nil != pms {
		if pms.UserID == constants.AdminUserID {
			return errors.New("不能修改管理员权限")
		}
	}
	err = fpms.pms.UpdateFPermission(fpm)
	fpms.pmc.initPms2Memory()
	return err
}

// DelFPermission 根据permissionID删除权限
func (fpms *FilePermission) DelFPermission(permissionID string) error {
	pms, err := fpms.pms.QueryFPermission(permissionID)
	if nil != err {
		return err
	}
	if nil != pms {
		if pms.UserID == constants.AdminUserID {
			return errors.New("不能删除管理员权限")
		}
	}
	err = fpms.pms.DelFPermission(permissionID)
	fpms.pmc.initPms2Memory()
	return err
}

// mkSqliteDIR 创建sqlite文件存放目录
func (fpms *FilePermission) mkSqliteDIR() {
	if !fileutil.IsExist("./.datas") {
		if err := fileutil.MkdirAll("./.datas"); nil != err {
			logs.Panicln(err)
		}
	}
}
