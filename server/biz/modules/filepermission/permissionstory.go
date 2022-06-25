// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 用sqlite3存放数据

package filepermission

import (
	"database/sql"
	"fileservice/biz/constants"
	"fileservice/biz/service"
	"time"

	"github.com/wup364/pakku/utils/strutil"
)

// PermissionStory 文件权限入口
type PermissionStory struct {
	db *sql.DB
}

// Initial 初始化配置
func (pms *PermissionStory) Initial(st constants.DBSetting) (err error) {
	if nil == pms.db {
		pms.db, err = sql.Open(st.DriverName, st.DataSourceName)
		if nil == err {
			if st.DriverName == "sqlite3" {
				// database is locked
				// https://github.com/mattn/go-sqlite3/issues/209
				pms.db.SetMaxOpenConns(1)
			} else {
				pms.db.SetMaxIdleConns(250)
				pms.db.SetConnMaxLifetime(time.Hour)
			}
		}
	}
	return err
}

// Install 初始化 users 表
func (pms *PermissionStory) Install() (err error) {
	var tx *sql.Tx
	if tx, err = pms.db.Begin(); err == nil {
		if _, err = tx.Exec(
			`CREATE TABLE IF NOT EXISTS filepermissions(
				permissionid VARCHAR(64) PRIMARY KEY,
				path TEXT(1000) NULL,
				userid VARCHAR(64) NULL,
				permission INTEGER(255) NULL,
				cttime DATE NULL
			);`); nil == err {
			//
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	}
	return err
}

// ListFPermissions 列出所有权限数据, 无分页
func (pms *PermissionStory) ListFPermissions() ([]service.PermissionInfo, error) {
	rows, err := pms.db.Query("SELECT permissionid, path, userid, permission, cttime FROM filepermissions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//
	res := make([]service.PermissionInfo, 0)
	for rows.Next() {
		fpms := service.PermissionInfo{}
		if err := rows.Scan(&fpms.PermissionID, &fpms.Path, &fpms.UserID, &fpms.Permission, &fpms.CtTime); err != nil {
			return nil, err
		}
		res = append(res, fpms)
	}
	return res, nil
}

// ListUserFPermissions 列出用户所有权限数据, 无分页
func (pms *PermissionStory) ListUserFPermissions(userID string) ([]service.PermissionInfo, error) {
	rows, err := pms.db.Query("SELECT permissionid, path, userid, permission, cttime FROM filepermissions WHERE userid = '" + userID + "'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//
	res := make([]service.PermissionInfo, 0)
	for rows.Next() {
		fpms := service.PermissionInfo{}
		if err := rows.Scan(&fpms.PermissionID, &fpms.Path, &fpms.UserID, &fpms.Permission, &fpms.CtTime); err != nil {
			return nil, err
		}
		res = append(res, fpms)
	}
	return res, nil
}

// QueryFPermission 根据权限ID查询详细信息
func (pms *PermissionStory) QueryFPermission(permissionID string) (*service.PermissionInfo, error) {
	rows, err := pms.db.Query("SELECT permissionid, path, userid, permission, cttime FROM filepermissions where permissionid='" + permissionID + "'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//
	if rows.Next() {
		fpms := service.PermissionInfo{}
		if err := rows.Scan(&fpms.PermissionID, &fpms.Path, &fpms.UserID, &fpms.Permission, &fpms.CtTime); err != nil {
			return nil, err
		}
		return &fpms, nil
	}
	return nil, nil
}

// AddFPermission 添加权限
func (pms *PermissionStory) AddFPermission(fpms service.PermissionInfo) (err error) {
	if len(fpms.UserID) == 0 {
		return service.ErrorUserIDIsNil
	}
	if len(fpms.Path) == 0 {
		return service.ErrorPermissionPathIsNil
	}
	if fpms.Permission <= 0 {
		return service.ErrorPermissionIsNil
	}
	// 开启事务
	var ts *sql.Tx
	if ts, err = pms.db.Begin(); err != nil {
		return err
	}
	//
	var stmt *sql.Stmt
	if stmt, err = ts.Prepare("INSERT INTO filepermissions(permissionid, path, userid, permission, cttime) values(?,?,?,?,?)"); err != nil {
		return err
	}
	if _, err = stmt.Exec(strutil.GetUUID(), fpms.Path, fpms.UserID, fpms.Permission, time.Now()); err != nil {
		ts.Rollback()
	} else {
		err = ts.Commit()
	}
	return err
}

// UpdateFPermission 修改权限
func (pms *PermissionStory) UpdateFPermission(fpms service.PermissionInfo) (err error) {
	if len(fpms.PermissionID) == 0 {
		return service.ErrorPermissionIDIsNil
	}
	if fpms.Permission <= 0 {
		return service.ErrorPermissionIsNil
	}
	// 开启事务
	var ts *sql.Tx
	if ts, err = pms.db.Begin(); err != nil {
		return err
	}
	//
	var stmt *sql.Stmt
	// permissionid, path, userid, permission, cttime
	if stmt, err = ts.Prepare("UPDATE filepermissions SET permission=? WHERE permissionid=?"); err != nil {
		return err
	}
	if _, err = stmt.Exec(fpms.Permission, fpms.PermissionID); err != nil {
		ts.Rollback()
	} else {
		err = ts.Commit()
	}
	return err
}

// DelFPermission 删除权限
func (pms *PermissionStory) DelFPermission(permissionID string) (err error) {
	if len(permissionID) == 0 {
		return service.ErrorPermissionIDIsNil
	}
	// 开启事务
	var ts *sql.Tx
	if ts, err = pms.db.Begin(); err != nil {
		return err
	}
	//
	var stmt *sql.Stmt
	if stmt, err = ts.Prepare("DELETE FROM filepermissions WHERE permissionid = ?"); err != nil {
		return err
	}
	if _, err = stmt.Exec(permissionID); err != nil {
		ts.Rollback()
	} else {
		err = ts.Commit()
	}
	return err
}
