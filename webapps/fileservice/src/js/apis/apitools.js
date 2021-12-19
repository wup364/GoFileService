// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
import { $utils } from "../utils";
import _ from "../prototypes"
import sha256 from 'crypto-js/sha256';

// API请求包装
"use strict";
export const $apitools = {
	_Const: {
		ack: "access_object",
		host: "access_host",
		hosturi: localStorage.getItem("access_host"),
	},
	/**
	 * 设置API主机地址
	 */
	setAPIHost: function (uri) {
		if (uri && uri.length > 0) {
			localStorage.setItem(this._Const.host, uri);
			this._Const.hosturi = uri;
		} else {
			localStorage.removeItem(this._Const.host);
		}
	},
	getAPIHost: function () {
		return this._Const.hosturi ? this._Const.hosturi : '';
	},
	/**
	 * 生成一个完整的url
	 * 如: /a/b/b --> https://127.0.0.1/a/b/c
	 */
	buildAPIURL: function (url) {
		return this._Const.hosturi ? this._Const.hosturi + url : window.location.origin + url;
	},
	// 构建签名url
	getSignAPIURL: function (url, params) {
		if (!params) { params = {}; }
		let session = this.getSession();
		if (!session || !session.secretKey || !session.accessKey) {
			this.apiResponseStautsHandler({ code: 401 }); return;
		}
		// 去除无效字段
		let paramsmap = new Map();
		for (let key in params) {
			if (params[key] == undefined || params[key] == null || params[key].length == 0) {
				continue;
			}
			paramsmap.set(key, params[key]);
		}
		// 构建请求负载
		let payloads = []; let payloads_encode = [];
		if (paramsmap.size > 0) {
			let keys = paramsmap.keys().sort();
			for (let i = 0; i < keys.length; i++) {
				let val = paramsmap.get(keys[i]);
				payloads.push(keys[i] + "=" + val);
				payloads_encode.push(keys[i] + "=" + encodeURIComponent(val));
			}
		}
		if (payloads.length > 0) {
			return url + '?' + payloads_encode.join("&") + '&x-Ack=' + session.accessKey + "&X-Sign=" + sha256(payloads.join("&") + session.secretKey)
		} else {
			return url + '?x-Ack=' + session.accessKey + '&X-Sign=' + sha256(session.secretKey)
		}
	},
	/**
	 * Api自动签名请求
	 */
	apiRequest: function (opt) {
		let session = this.getSession();
		opt.session = session;
		if (!session || !session.secretKey || !session.accessKey) {
			this.apiResponseStautsHandler({ code: 401 }); return;
		}
		if (this._Const.hosturi && this._Const.hosturi.length > 0) {
			if (!opt.uri.startWith('http://') && !opt.uri.startWith('https://')) {
				opt.uri = this._Const.hosturi + opt.uri;
			}
		}
		let request = $utils.AjaxRequest(opt);
		let signPayload = request.payload.payload;
		if (signPayload && signPayload.length > 0) {
			signPayload += opt.session.secretKey;
		} else {
			signPayload = opt.session.secretKey;
		}
		request.setHeader("X-Ack", session.accessKey);
		request.setHeader("X-Sign", sha256(signPayload));
		return request;
	},
	/**
	 * 相应结构格式化
	 */
	apiResponseFormat: function (xhr) {
		let result = {
			code: 0,
			data: '',
		};
		// 请求错误&没有响应数据
		// if( xhr.status !== 200 ){
		// 	result.code = xhr.status;
		// 	result.data = xhr.statusText;
		// }else{
		let _obj = {};
		if (undefined != xhr.responseText) {
			if ($utils.isString(xhr.responseText)) {
				try {
					_obj = JSON.parse(xhr.responseText);
				} catch (err) {
					_obj = xhr.responseText;
				}
			} else {
				_obj = xhr.responseText;
			}
		}
		// 没有结构化的数据返回
		if (undefined == _obj.code && undefined == _obj.data) {
			result.code = xhr.status;
			result.data = _obj;
		} else {
			result.code = undefined == _obj.code ? xhr.status : _obj.code;
			result.data = undefined == _obj.data ? '' : _obj.data;
		}
		// }
		this.apiResponseStautsHandler(result);
		// console.log("apiResponseFormat: ", result, xhr);
		return result;
	},
	// Api 状态返回翻译和处理
	apiResponseStautsHandler: function (res) {
		if (res) {
			if (res.code == 401) {
				res.data = "登录过期,请刷新页面";
				(top.location || window.location).href = "/";
			}
			else if (res.code == 403) {
				res.data = "请求被禁止,可能是权限不足";
			}
		}
	},
	/**
	 * APi-Get请求
	 */
	apiGet: function (uri, datas) {
		return new Promise((resolve, reject) => {
			this.apiRequest({
				uri: uri,
				datas: datas,
			}).do((xhr, opt) => {
				if (xhr.readyState === 4) {
					let res = this.apiResponseFormat(xhr);
					if (res.code === 200) {
						resolve(res.data);
					} else {
						reject(res.data);
					}
				}
			});
		})
	},
	/**
	 * APi-Post请求
	 */
	apiPost: function (uri, datas) {
		return new Promise((resolve, reject) => {
			this.apiRequest({
				uri: uri,
				method: "POST",
				datas: datas,
			}).do((xhr, opt) => {
				if (xhr.readyState === 4) {
					let res = this.apiResponseFormat(xhr);
					if (res.code === 200) {
						resolve(res.data);
					} else {
						reject(res.data);
					}
				}
			});
		})
	},
	apiPostSync: function (uri, datas) {
		return this.apiRequest({
			uri: uri,
			method: "POST",
			datas: datas,
			async: false,
		}).do(function (xhr, opt) {
			if (xhr.readyState === 4) {
				let res = this.apiResponseFormat(xhr);
				if (res.code === 200) {
					return res.data;
				} else {
					throw res.data;
				}
			}
		});
	},
	/**
	 * APi-Delete请求
	 */
	apiDelete: function (uri, datas) {
		return new Promise(function (resolve, reject) {
			this.apiRequest({
				uri: uri,
				method: "DELETE",
				datas: datas,
			}).do(function (xhr, opt) {
				if (xhr.readyState === 4) {
					let res = this.apiResponseFormat(xhr);
					if (res.code === 200) {
						resolve(res.data);
					} else {
						reject(res.data);
					}
				}
			});
		})
	},
	apiDeleteSync: function (uri, datas) {
		return this.apiRequest({
			uri: uri,
			method: "DELETE",
			datas: datas,
			async: false,
		}).do(function (xhr, opt) {
			if (xhr.readyState === 4) {
				let res = this.apiResponseFormat(xhr);
				if (res.code === 200) {
					return res.data;
				} else {
					throw res.data;
				}
			}
		});
	},
	/**
	 * 保存会话到cookie中
	 */
	saveSession: function (accessObj) {
		localStorage.setItem(this._Const.ack, "");
		if (accessObj && accessObj.accessKey) {
			localStorage.setItem(this._Const.ack, JSON.stringify(accessObj));
			$utils.setCookie("X-Ack", accessObj.accessKey);
		}
	},
	/**
	 * 获取会话信息
	 * {UserId, accessKey, secretKey}
	 */
	getSession: function () {
		try {
			return JSON.parse(localStorage.getItem(this._Const.ack));
		} catch (err) {
			return {};
		}
	},
	/**
	 * 注销会话
	 */
	destroySession: function () {
		localStorage.setItem(this._Const.ack, '');
		$utils.setCookie("X-Ack", '');
	},
};

/*this.saveSession({
	UserId: "system",
	accessKey: "0000000000000000000000000000",
	secretKey: "0000000000000000000000000000",
});*/