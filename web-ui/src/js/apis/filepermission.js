// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

"use strict";
import { $apitools } from "./apitools";
export const $filepms = {
	// 权限类型
	$TYPE: {
		// VISIBLECHILD 子节点可见导致的父节点可见, 父节点其实并没有权限
		VISIBLECHILD: { name: '子节点可见', value: 0 },
		// VISIBLE 可见
		VISIBLE: { name: '可见', value: 1 },
		// READ 只读
		READ: { name: '只读', value: 2 },
		// WRITE 可写
		WRITE: { name: '可写', value: 3 },
	},
	// listFPermissions
	listFPermissions() {
		return $apitools.apiGet("/filepms/v1/listfpermissions", {})
	},
	// GetUserPermissionSum
	GetUserPermissionSum(userid, paths) {
		return $apitools.apiPost("/filepms/v1/getuserpermissionsum", {
			userid: userid ? userid : '',
			paths: paths ? JSON.stringify(paths) : ''
		})
	},
	// listUserFPermissions
	listUserFPermissions(userid) {
		return $apitools.apiGet("/filepms/v1/listuserfpermissions", { userid: userid ? userid : '' })
	},
	// addFPermission
	addFPermission(userid, path, permission) {
		return $apitools.apiPost("/filepms/v1/addfpermission", {
			'userid': userid ? userid : '',
			'path': path ? path : '',
			'permission': permission ? permission : '',
		});
	},
	// delFPermission
	delFPermission(permissionid) {
		return $apitools.apiDelete("/filepms/v1/delfpermission", {
			'permissionid': permissionid ? permissionid : '',
		});
	},
	// updateFPermission
	updateFPermission(permissionid, permission) {
		return $apitools.apiPost("/filepms/v1/updatefpermission", {
			"permissionid": permissionid ? permissionid : '',
			"permission": permission ? permission : '',
		});
	},
	sync: {
		// delFPermission
		delFPermission(permissionid) {
			return $apitools.apiDeleteSync("/filepms/v1/delfpermission", {
				'permissionid': permissionid ? permissionid : '',
			});
		},
	}
};

// 权限类型计算
$filepms.$TYPE.__proto__ = {
	// 是否包含某权限
	sumInclude(sum, p) {
		sum = Number(sum); p = Number(p);
		if (undefined == sum || sum < 0 ||
			undefined == p || p < 0) {
			return false;
		}
		return 1 << p == (sum & (1 << p));
	},
	// 转换为描述文字
	sum2Name(sum) {
		let res = [];
		if (undefined != sum && sum > 0) {
			for (let key in $filepms.$TYPE) {
				if (key == $filepms.$TYPE.VISIBLECHILD.name ||
					!$filepms.$TYPE.hasOwnProperty(key)) {
					continue;
				}
				if ($filepms.$TYPE.sumInclude(sum, $filepms.$TYPE[key].value)) {
					res.push($filepms.$TYPE[key].name);
				}
			}
		}
		if (res.length == 0) {
			if ($filepms.$TYPE.sumInclude(sum, $filepms.$TYPE.VISIBLECHILD.value)) {
				res.push($filepms.$TYPE.VISIBLECHILD.name);
			}
		}
		return res;
	},
	// 计算权限结果值
	list2Sum(pms) {
		let res = -1;
		if (pms && pms.length > 0) {
			for (let i = 0; i < pms.length; i++) {
				res += 1 << Number(pms[i]);
			}
		}
		return res;
	},
};


