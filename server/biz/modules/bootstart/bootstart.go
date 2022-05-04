// Copyright (C) 2021 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

//

package bootstart

import (
	"fileservice/biz/constants"
	"fileservice/biz/service"
	"pakku/ipakku"
	"pakku/utils/logs"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// BootStart bootstart
type BootStart struct {
	au  service.User4RPC       `@autowired:"User4RPC"`
	pms service.FilePermission `@autowired:"FilePermission"`
}

// AsModule Template接口实现
func (b *BootStart) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "BootStart",
		Version:     1.0,
		Description: "BootStart",
		OnSetup: func() {
			b.createDefaultUser()
			b.createDefaultPms()
		},
	}
}

// createDefaultUser 创建默认账号
func (bs *BootStart) createDefaultUser() {
	if err := bs.au.AddUser(&service.UserInfo{
		UserID:   constants.AdminUserID,
		UserName: "管理员",
		UserType: 0,
	}); nil != err {
		logs.Panicln(err)
	}
}

// createDefaultPms 创建默认权限
func (bs *BootStart) createDefaultPms() {
	if err := bs.pms.AddFPermission(service.PermissionInfo{
		Path:       "/",
		UserID:     constants.AdminUserID,
		Permission: (1 << service.FPM_Visible) + (1 << service.FPM_Read) + (1 << service.FPM_Write),
	}); nil != err {
		logs.Panicln(err)
	}
}
