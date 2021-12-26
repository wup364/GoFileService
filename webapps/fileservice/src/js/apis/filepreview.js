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
export const $filepreview = {
	supported: {
		audio: ['mp3', 'flac'],
		video: ['mp4', 'mov'],
		picture: ['png', 'gif', 'jpg', 'bmp', 'jpeg', 'icon'],
	},
	// 获取一个预览的Token
	askToken(path) {
		return $apitools.apiGet("/filepreview/v1/asktoken", {
			"path": path ? path : ""
		});
	},
	// Status
	status(token) {
		return $apitools.apiGet("/filepreview/v1/status/" + token);
	},
	// 获取预览文件的同级目录文件
	samedirFiles(token) {
		return $apitools.apiGet("/filepreview/v1/samedirfiles/" + token);
	},
	// 获取预览文件的文件流
	buildStreamURL(token, fileName) {
		return $apitools.buildAPIURL("/filepreview/v1/stream/" + token + "?fileName=" + (fileName ? fileName : ""));
	},
	// 预览
	doPreview(path, suffix) {
		return new Promise((resolve, reject) => {
			$filepreview.askToken(path).then((data) => {
				$filepreview.openPreview(data, suffix ? suffix.toLowerCase() : path.getSuffixed(false).toLowerCase());
				resolve();
			}).catch(reject);
		});
	},
	// 打开预览地址
	openPreview(token, suffix) {
		let type = '';
		if ($filepreview.isSupport('audio', suffix)) {
			type = 'audio';
		} else if ($filepreview.isSupport('video', suffix)) {
			type = 'video';
		} else if ($filepreview.isSupport('picture', suffix)) {
			type = 'picture';
		} else {
			throw "不支持预览该文件";
		}
		window.open("#/preview/" + type + "?token=" + token);
	},
	// 是否是支持的播放类型
	isSupport(type, suffix) {
		suffix = suffix.toLowerCase();
		return type && suffix && $filepreview.supported[type] && $filepreview.supported[type].indexOf(suffix) > -1;
	}
};
