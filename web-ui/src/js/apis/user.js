// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
import { $utils } from "../utils";
import { $apitools } from "./apitools";

"use strict";
export const $userApi = {
	$TYPES: {
		0: "管理员",
		1: "普通用户",
	},
	// 登录
	login(user, pwd) {
		return new Promise((resolve, reject) => {
			$utils.AjaxRequest({
				method: "POST",
				uri: $apitools.buildAPIURL("/user/v1/checkpwd"),
				datas: {
					"userid": user,
					"pwd": pwd,
				},
			}).do((xhr, opt) => {
				if (xhr.readyState === 4) {
					var res = $apitools.apiResponseFormat(xhr);
					if (res.code === 200) {
						resolve(res.data);
					} else {
						reject(res.data);
					}
				}
			});
		});
	},
	// logout
	logout() {
		return $apitools.apiPost("/user/v1/logout")
	},
	// QueryUser
	queryuser(userid) {
		return $apitools.apiGet("/user/v1/queryuser", {
			"userid": userid
		})
	},
	// UpdateUserName
	updateUserName(userid, username) {
		return $apitools.apiPost("/user/v1/updateusername", {
			"userid": userid,
			"username": username,
		})
	},
	// UpdateUserPwd
	updateUserPwd(userid, userpwd) {
		return $apitools.apiPost("/user/v1/updateuserpwd", {
			"userid": userid,
			"userpwd": userpwd,
		})
	},
	// ListAllUsers
	listAllUsers(userid, userpwd) {
		return $apitools.apiGet("/user/v1/listallusers", {})
	},
	// AddUser
	addUser(userid, username, userpwd) {
		return $apitools.apiPost("/user/v1/adduser", {
			'userid': userid,
			'username': username,
			'userpwd': userpwd ? userpwd : '',
		})
	},
	// DelUser
	delUser(userid) {
		return $apitools.apiDelete("/user/v1/deluser", {
			'userid': userid,
		})
	},
	sync: {
		// UpdateUserName
		updateUserName(userid, username) {
			return $apitools.apiPostSync("/user/v1/updateusername", {
				"userid": userid,
				"username": username,
			})
		},
		// UpdateUserPwd
		updateUserPwd(userid, userpwd) {
			return $apitools.apiPostSync("/user/v1/updateuserpwd", {
				"userid": userid,
				"userpwd": userpwd,
			})
		},
		// DelUser
		delUser(userid) {
			return $apitools.apiDeleteSync("/user/v1/deluser", {
				'userid': userid,
			})
		},
	}
};

// 翻译用户类型
$userApi.$TYPES.__proto__ = {
	parse(userType) {
		return $userApi.$TYPES[userType]
			? $userApi.$TYPES[userType]
			: "未知类型";
	}
};