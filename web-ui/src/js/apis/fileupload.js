// Copyright (C) 2022 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

"use strict";
// 获取图标URL
export const $fileupload = {
	opts: {
		chunkSize: 128 * 1024 * 1024,
	},

	// upload 
	upload(opts, ctrl) {
		let _opts = Object.assign({ uploadMethod: '', streamURL: '', submitURL: '', file: null }, opts);
		if (_opts.uploadMethod == "chunkUploader") {
			return this.chunkUploader(_opts.streamURL, _opts.submitURL, _opts.file, ctrl);
		} else {
			return this.formDataUploader(_opts.streamURL, _opts.file, ctrl);
		}
	},

	// chunkUploader 分片上传器
	chunkUploader(streamurl, submiturl, file, ctrl) {
		let uploader = {
			loadSize: 0,
			aborted: false,
			currentChunckNum: -1,
			currentReader: null,
			currentUploader: null,
			totalSize: file.size ? file.size : 0,
			opts: Object.assign({
				loadstart: (e) => { },
				load: (e) => { },
				loadend: (e) => { },
				error: (e) => { },
				abort: (e) => { },
				progress: (e) => { },
			}, ctrl),
			abort: () => {
				uploader.aborted = true;
				if (uploader.currentReader) {
					try {
						uploader.currentReader.abort();
					} catch (error) {
						console.error(error);
					}
				}
				if (uploader.currentUploader) {
					try {
						uploader.currentUploader.abort();
					} catch (error) {
						console.error(error);
					}
				}
			},
			start: () => {
				if (uploader.currentChunckNum === -1) {
					uploader.currentChunckNum++;
					let chunck = this.cuteFile(file, this.opts.chunkSize);
					if (chunck) {
						doUpload(chunck);
					}
				}
			}
		};

		// 分片
		let doUpload = (chunck) => {
			let chunckURL = this.buildParams(streamurl, { number: ++uploader.currentChunckNum });
			uploader.currentUploader = this.formDataUploader(chunckURL, chunck, {
				filekey: "file",
				method: uploader.opts.method,
				header: uploader.opts.header,
				form: uploader.opts.header,
				loadstart: (e) => {
					if (1 === uploader.currentChunckNum) {
						uploader.opts.loadstart(e);
					}
				},
				loadend: (e) => {
					if (!uploader.aborted) {
						let chunck = this.cuteFile(file, this.opts.chunkSize);
						if (chunck) {
							setTimeout(() => {
								doUpload(chunck);
							});
						} else {
							let doSubmit = () => {
								if (uploader.currentUploader.xhr.readyState !== 4) {
									setTimeout(doSubmit, 50);
								} else {
									try {
										this.AjaxRequest({
											method: "POST",
											uri: submiturl,
										}).do((xhr) => {
											if (xhr.readyState === 4) {
												if (xhr.status === 200) {
													uploader.opts.loadend();
												} else {
													uploader.opts.error(new Error(xhr.statusText));
												}
											}
										})
									} catch (error) {
										uploader.opts.error(error);
									}
								}
							};
							doSubmit();
						}
					} else if (!uploader.currentUploader._aborted) {
						uploader.currentUploader._aborted = true;
						uploader.opts.abort(e);
					}
				},
				error: (e) => {
					uploader.aborted = true;
					uploader.currentUploader._aborted = true;
					uploader.opts.error(e);
				},
				abort: (e) => {
					if (!uploader.currentUploader._aborted) {
						uploader.currentUploader._aborted = true;
						uploader.opts.abort(e);
					}
				},
				progress: (e) => {
					uploader.loadSize = (uploader.currentChunckNum ? (uploader.currentChunckNum - 1) * this.opts.chunkSize : 0) + e.loaded;
					if (uploader.loadSize > uploader.totalSize) {
						uploader.loadSize = uploader.totalSize;
					}
					uploader.opts.progress({
						loaded: uploader.loadSize,
						total: uploader.totalSize,
					});
				},
			});
			uploader.currentUploader.start();
		}

		return uploader;
	},
	buildParams(url, params) {
		if (params) {
			let payloads = [];
			let keys = Object.keys(params);
			for (const element of keys) {
				if (undefined !== params[element]) {
					payloads.push(element + "=" + encodeURIComponent(params[element]));
				}
			}
			url = url + "?" + payloads.join("&");
		}
		return url;
	},
	cuteFile(file, chunk_size) {
		if (!file || !chunk_size) {
			return;
		}
		file._cute = file._cute ? {
			_start: file._cute._end,
			_end: file._cute._end + chunk_size
		} : {
			_start: 0,
			_end: chunk_size
		};
		if (file.size && file._cute._start >= file.size) {
			return;
		}
		return file.slice(file._cute._start, file._cute._end);
	},
	// Form表单POST上传
	formDataUploader(url, file, ctrl) {
		let uploader = {
			ctrl: {
				filekey: file.name ? file.name : "file",
				method: "POST",
				header: {},
				form: {},
				loadstart: (e) => { },
				load: (e) => { },
				loadend: (e) => { },
				error: (e) => { },
				abort: (e) => { },
				progress: (e) => { },
			},
			xhr: new XMLHttpRequest(),
			formData: new FormData(),
			abort: () => {
				uploader.xhr.abort();
			},
			start: () => {
				try {
					uploader.xhr.open(uploader.ctrl.method, url);
					for (let key in uploader.ctrl.header) {
						uploader.xhr.setRequestHeader(key, uploader.ctrl.header[key])
					}
					uploader.xhr.overrideMimeType("application/octet-stream");
					uploader.xhr.send(uploader.formData);
					uploader.xhr.onerror = uploader.ctrl.error;
				} catch (error) {
					uploader.ctrl.error(error);
				}
			}
		};
		// 合并配置
		uploader.ctrl = Object.assign(uploader.ctrl, ctrl);
		// 添加表单
		for (let key in uploader.ctrl.form) {
			uploader.formData.append(key, uploader.ctrl.form[key]);
		}
		uploader.formData.append(uploader.ctrl.filekey, file);
		// 绑定事件
		uploader.xhr.upload.addEventListener('loadstart', uploader.ctrl.loadstart);
		uploader.xhr.upload.addEventListener('load', uploader.ctrl.load);
		uploader.xhr.upload.addEventListener("loadend", (e) => {
			let doLoadend = () => {
				if (uploader.xhr.readyState !== 4) {
					setTimeout(doLoadend, 50);
				} else {
					uploader.ctrl.loadend(e);
				}
			};
			doLoadend();
		});
		uploader.xhr.upload.addEventListener("progress", uploader.ctrl.progress);
		uploader.xhr.upload.addEventListener('error', uploader.ctrl.error);
		uploader.xhr.upload.addEventListener('abort', uploader.ctrl.abort);
		uploader.xhr.onreadystatechange = (e) => {
			let xhr = e.target;
			if (xhr.readyState == 4 && xhr.status != 200) {
				uploader.ctrl.error(new Error(xhr.responseText ? xhr.responseText : xhr.statusText));
			}
		};
		return uploader;
	},
	// Ajax请求, opt={uri, method, datas, headers, async}
	AjaxRequest(opt) {
		opt.xhrPayload = "";
		opt.async = undefined == opt.async ? true : opt.async;
		opt.datas = opt.datas ? opt.datas : {};
		opt.method = opt.method ? opt.method.toUpperCase() : "GET";
		opt.headers = opt.headers ? opt.headers : {};

		{
			// 构建请求参数
			opt.payload = this.buildPayload(opt.datas);
			opt.xhrPayload = opt.payload.payload_encode;
			// GET 通过Url传递参数
			if (opt.xhrPayload && opt.xhrPayload.length > 0) {
				if ("GET" == opt.method || "DELETE" == opt.method) {
					opt.uri = opt.uri + "?" + opt.xhrPayload;
					opt.xhrPayload = "";
				}
			}
			// Post 通过form传递参数
			if ("POST" == opt.method || "DELETE" == opt.method) {
				if (undefined == opt.headers["Content-Type"]) {
					opt.headers["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8";
				}
			}
		}
		// 
		opt.xhr = new XMLHttpRequest();
		opt.xhr.open(opt.method, opt.uri, opt.async);
		opt.xhr.onreadystatechange = () => {
			if (opt.onreadystatechange) {
				return opt.onreadystatechange(opt.xhr, opt);
			}
		};
		{
			if (opt.timeout) {
				opt.xhr.timeout = parseInt(opt.timeout);
			}
			for (let key in opt.headers) {
				opt.xhr.setRequestHeader(key, opt.headers[key]);
			}
		}
		// 
		opt.setHeader = (key, val) => {
			opt.headers[key.toString()] = val;
			opt.xhr.setRequestHeader(key, val);
		};
		opt.do = (callback) => {
			if (opt.async && callback) {
				opt.onreadystatechange = callback;
			}
			opt.xhr.send(opt.xhrPayload);
			return opt.async || !callback ? opt.xhr : callback(opt.xhr, opt);
		};
		return opt;
	},
	// HTTP请求负载构建, 去除空字段并排序字段
	buildPayload(params) {
		if (!params) { return ""; }
		// 去除无效字段
		let paramsmap = {};
		for (let key in params) {
			if (params[key] == undefined || params[key] == null || params[key].length == 0) {
				continue;
			}
			paramsmap[key] = params[key];
		}
		// 构建请求负载
		let payload, payload_encode = "";
		let keys = Object.keys(paramsmap);
		if (keys.length > 0) {
			keys = keys.sort();
			let payloads = []; let payloads_encode = [];
			for (const element of keys) {
				let val = paramsmap[element];
				payloads.push(element + "=" + val);
				payloads_encode.push(element + "=" + encodeURIComponent(val));
			}
			payload = payloads.join("&");
			payload_encode = payloads_encode.join("&");
		}
		return {
			payload: payload,
			payload_encode: payload_encode,
		};
	},
};