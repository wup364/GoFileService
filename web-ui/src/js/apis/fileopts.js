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
export const $fileopts = {
	// 列表路径
	List(path) {
		return $apitools.apiGet("/file/v1/list", { "path": path ? path : '' });
	},
	// 移动文件|文件夹
	Delete(path) {
		return $apitools.apiDelete("/file/v1/del", {
			path: path ? path : '',
		});
	},
	// 重命名文件|文件夹
	Rename(path, name) {
		return $apitools.apiPost("/file/v1/rename", {
			path: path ? path : '',
			name: name ? name : '',
		});
	},
	// 新建文件夹
	NewFolder(path) {
		return $apitools.apiPost("/file/v1/newfolder", {
			path: path ? path : '',
		});
	},
	// ---------------------------------------------
	// 获取一个Token, data可以存放值
	GetStreamToken(type, datas) {
		return $apitools.apiGet("/filestream/v1/token", {
			"type": type ? type : '',
			"data": datas ? datas : ''
		});
	},
	// 获取一个下载的Url
	GetDownloadUrl(path) {
		return $fileopts.GetStreamToken('download', path).then((data) => {
			return $apitools.buildAPIURL($apitools.getSignAPIURL("/filestream/v1/read/" + data.token));
		});
	},
	// 获取一个打开的Url - 流
	GetSteamUrl(path) {
		return $fileopts.GetStreamToken('download', path).then((data) => {
			return encodeURI($apitools.buildAPIURL("/filestream/v1/read/" + data.token, { action: 'stream' }));
		});
	},
	// 获取一个上载的Url
	GetUploadUrl(path) {
		return $fileopts.GetStreamToken('upload', path).then((data) => {
			return $apitools.buildAPIURL("/filestream/v1/put/" + data.token);
		});
	},
	// ---------------------------------------------
	// 异步执行一个动作
	AsyncExec(func, params) {
		if (!params) {
			params = {};
		}
		params.func = func;
		return $apitools.apiPost("/filetask/v1/asyncexec", params);
	},
	// 异步执行Token查询 CopyFile, MoveFile
	AsyncExecToken(func, token, params) {
		if (!params) {
			params = {};
		}
		params.func = func;
		params.token = token;
		return $apitools.apiPost("/filetask/v1/asyncexectoken", params);
	},
	// 复制文件|文件夹
	CopyAsync(src, dest, replaceExist, ignoreError) {
		return this.AsyncExec("CopyFile", {
			srcPath: src ? src : '',
			dstPath: dest ? dest : '',
			replace: replaceExist ? replaceExist : false,
			ignore: ignoreError ? ignoreError : false
		});
	},
	// 移动文件|文件夹
	MoveAsync(src, dest, replaceExist, ignoreError) {
		return this.AsyncExec("MoveFile", {
			srcPath: src ? src : '',
			dstPath: dest ? dest : '',
			replace: replaceExist ? replaceExist : false,
			ignore: ignoreError ? ignoreError : false
		});
	},

};