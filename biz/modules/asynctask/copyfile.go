// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件异步操作, 一般用于批量操作或后台任务

package asynctask

import (
	"encoding/json"
	"errors"
	"fileservice/biz/service"
	"net/http"
	"pakku/ipakku"
	"pakku/utils/fileutil"
	"pakku/utils/logs"
	"pakku/utils/serviceutil"
	"pakku/utils/strutil"
	"time"
)

// CopyFileTokenObject 复制文件Token保存对象
type CopyFileTokenObject struct {
	ErrorString   string // 错误信息
	Src           string // 当前正在处理的源路径
	Dst           string // 当前正在处理的目标路径
	IsSrcExist    bool   // 源路径是否存在
	IsDstExist    bool   // 目标路径是否存在
	IsReplace     bool   // 是否替换, 单次中断执行指令, 读取后设为false
	IsReplaceAll  bool   // 是否替换, 单次API执行指令, 设置后后续中断时自动替换
	IsIgnore      bool   // 是否忽略错误, 单次中断执行指令, 读取后设为false
	IsIgnoreAll   bool   // 是否忽略错误, 单次API执行指令, 设置后后续中断时自动替换
	IsComplete    bool   // 是否执行完毕
	IsDiscontinue bool   // 是否已中断操作
}

// ToJSON 转传JSON
func (dto *CopyFileTokenObject) ToJSON() string {
	if bt, err := json.Marshal(dto); nil != err {
		return ""
	} else {
		return string(bt)
	}
}

// CopyFile 复制文件|文件夹
type CopyFile struct {
	c     ipakku.AppCache             `@autowired:"AppCache"`
	fm    service.FileDatas           `@autowired:"FileDatas"`
	sg    service.UserAuth4Rpc        `@autowired:"User4RPC"`
	pmc   service.FilePermissionCheck `@autowired:"FilePermission"`
	token *TaskToken
}

// Name 动作名字
func (task CopyFile) Name() string {
	return "CopyFile"
}

// Init 初始化对象
func (task *CopyFile) Init(mctx ipakku.Loader) service.AsyncTaskExecI {
	if err := mctx.AutoWired(task); nil != err {
		logs.Panicln(err)
	}
	task.token = NewTaskToken(task.c)
	return task
}

// Execute 动作执行, 返回一个tooken
func (task *CopyFile) Execute(r *http.Request) (string, error) {
	qSrcPath := r.FormValue("srcPath")
	qDstPath := r.FormValue("dstPath")
	qReplace := strutil.String2Bool(r.FormValue("replace"))
	qIgnore := strutil.String2Bool(r.FormValue("ignore"))

	if len(qSrcPath) == 0 {
		return "", errors.New("srcPath parameter not found")
	}
	if len(qDstPath) == 0 {
		return "", errors.New("dstPath parameter not found")
	}
	userID := task.GetUserID4Request(r)
	if !task.pmc.HashPermission(userID, qSrcPath, service.FPM_Read) {
		return "", service.ErrorPermissionInsufficient
	}
	if !task.pmc.HashPermission(userID, qDstPath, service.FPM_Write) {
		return "", service.ErrorPermissionInsufficient
	}
	// 异步处理, 返回一个Token用于查询进度
	token, err := task.token.AskToken(&CopyFileTokenObject{
		ErrorString:  "",
		Src:          qSrcPath,
		Dst:          qDstPath,
		IsSrcExist:   true,
		IsDstExist:   false,
		IsReplace:    false,
		IsReplaceAll: qReplace,
		IsIgnore:     false,
		IsIgnoreAll:  qIgnore,
	})
	if nil != err {
		return "", err
	}
	go func(token string) {
		// 这里面已经不属于一个会话, 使用令牌保存数据
		copyDirErr := task.doCopy(qSrcPath, qDstPath, qReplace, qIgnore, func(s_src, s_dst string, copyErr *CopyError) error {
			// 获取令牌数据, 不存在则说明已经销毁
			var tokenBody *CopyFileTokenObject
			if tokenObj, err := task.getTokenObject(token); nil != err {
				return err
			} else if tokenObj.IsDiscontinue {
				return service.ErrorDiscontinue
			} else {
				tokenBody = tokenObj
				tokenBody.IsIgnore = false
				tokenBody.IsReplace = false
				tokenBody.IsSrcExist = false
				tokenBody.IsDstExist = false
				tokenBody.ErrorString = ""
				tokenBody.Src = s_src
				tokenBody.Dst = s_dst
				// 更新缓存
				if err := task.token.RefreshToken(token, tokenBody); nil != err {
					logs.Errorln(err)
				}
			}
			// 如果遇到错误了
			if nil != copyErr {
				// 判断是否是目标位置已经存在的错误, 如果是的话需要选择是否覆盖他
				if copyErr.DstIsExist {
					// 查找之前是否设置了 替换全部错误
					if tokenBody.IsReplaceAll {
						// 先删除然后再替换, 如果覆盖操作没有出现问题
						if reCopyErr := task.doCopy(s_src, s_dst, true, false, nil); nil == reCopyErr {
							return nil
						} else {
							tokenBody.ErrorString = reCopyErr.Error()
						}
					}
					// 如果设置了自动覆盖, 但是任然出错, 则判断是否忽略错误选项
					if tokenBody.IsIgnoreAll {
						tokenBody.ErrorString = ""
						return nil // 跳过这个文件
					}
				} else {
					// 不是路径重复类错误
					// 如果是其他错误就不管了, 暂时无法处理只能选择 忽略|暂停
					// 查找之前是否设置了 忽略全部错误
					if tokenBody.IsIgnoreAll {
						return nil // 跳过这个文件
					}
				}

				// 到此说明 没有设置自动覆盖和自动忽略
				tokenBody.IsSrcExist = copyErr.SrcIsExist
				tokenBody.IsDstExist = copyErr.DstIsExist
				// 设置错误, 等待客户端获取, 等待操作
				if len(tokenBody.ErrorString) == 0 {
					tokenBody.ErrorString = copyErr.ErrorString
				}
				// 更新缓存-抛出错误
				if err := task.token.RefreshToken(token, tokenBody); nil != err {
					logs.Errorln(err)
				}
				for {
					// 循环读取最想指令
					if token, err := task.getTokenObject(token); nil != err {
						return err
					} else if token.IsDiscontinue {
						return service.ErrorDiscontinue
					} else {
						tokenBody = token
					}
					// 选择了忽略|忽略全部
					if tokenBody.IsIgnore || tokenBody.IsIgnoreAll {
						return nil
					}
					// 选择了覆盖|覆盖全部
					if tokenBody.IsReplace || tokenBody.IsReplaceAll {
						if copyErr.SrcIsExist {
							// 先删除然后再替换
							if reCopyErr := task.doCopy(s_src, s_dst, true, false, nil); nil == reCopyErr {
								return nil
							} else {
								tokenBody.ErrorString = reCopyErr.Error()
							}
						} else {
							tokenBody.ErrorString = fileutil.PathNotExist("replace", s_src).Error()
						}
						// 更新缓存-抛出错误
						if err := task.token.RefreshToken(token, tokenBody); nil != err {
							logs.Errorln(err)
						}
					}
					time.Sleep(time.Duration(100) * time.Millisecond) // 休眠100ms
				}
			}
			return nil
		})
		// 到这里如果没有错误就是成功了
		if tokenBody, err := task.getTokenObject(token); nil == err {
			if nil != copyDirErr {
				tokenBody.ErrorString = copyDirErr.Error()
				logs.Errorln(copyDirErr)
			} else {
				tokenBody.ErrorString = ""
			}
			tokenBody.IsComplete = true
			tokenBody.IsDiscontinue = service.ErrorDiscontinue.Error() == tokenBody.ErrorString
			if err := task.token.RefreshToken(token, tokenBody); nil != err {
				logs.Errorln(err)
			}
		} else {
			logs.Errorln(err)
		}
	}(token)
	return token, nil
}

// doCopy 递归copy 重复策略 replace|ignore
func (task *CopyFile) doCopy(src, dst string, replace, ignore bool, walk func(s_src, s_dst string, copyErr *CopyError) error) error {
	doWalk := func(src, dst string, err error) error {
		if nil == walk {
			return err
		}
		if nil != err {
			return walk(src, dst, &CopyError{
				SrcIsExist:  task.fm.IsExist(src),
				DstIsExist:  task.fm.IsExist(dst),
				ErrorString: err.Error(),
			})
		} else {
			return walk(src, dst, nil)
		}
	}
	if task.fm.IsFile(src) {
		if !replace && task.fm.IsExist(dst) {
			if err := doWalk(src, dst, fileutil.PathExist("copy", dst)); nil != err {
				return err
			}
		} else {
			if err := doWalk(src, dst, task.fm.DoCopy(src, dst, replace)); nil != err {
				return err
			}
		}
	} else if task.fm.IsDir(src) {
		if names := task.fm.GetDirList(src); len(names) > 0 {
			for i := 0; i < len(names); i++ {
				child := src + "/" + names[i]
				childdst := dst + "/" + names[i]
				dstexist := task.fm.IsExist(childdst)
				if task.fm.IsFile(child) {
					// 如果不是覆盖重复模式, 则判断是否忽略
					if !replace && dstexist {
						if !ignore {
							if err := doWalk(child, childdst, fileutil.PathExist("copy", childdst)); nil != err {
								return err
							}
						} else {
							doWalk(child, childdst, nil)
							continue
						}
					} else {
						// 目标位置不存在|目标存在但是允许覆盖
						if err := doWalk(child, childdst, task.fm.DoCopy(child, childdst, replace)); nil != err {
							return err
						}
					}
				} else {
					// 如果不是覆盖重复模式&目标位置是个文件(应为空或文件夹), 则判断是否忽略
					if !replace && dstexist && task.fm.IsFile(childdst) {
						if !ignore {
							if err := doWalk(child, childdst, fileutil.PathExist("copy", childdst)); nil != err {
								return err
							}
						} else {
							doWalk(child, childdst, nil)
							continue
						}
					} else {
						if !dstexist {
							if err := doWalk(child, childdst, task.fm.DoMkDir(childdst)); nil != err {
								return err
							}
						}
						if err := task.doCopy(child, childdst, replace, ignore, walk); nil != err {
							return err
						}
					}
				}
			}
		}

	} else {
		return fileutil.PathNotExist("copy", src)
	}
	return nil
}

// Status 查询动作状态, 在内部返回数据
func (task *CopyFile) Status(w http.ResponseWriter, r *http.Request) {
	qToken := r.FormValue("token")
	qOperation := r.FormValue("operation")
	tokenBody, tokenErr := task.getTokenObject(qToken)
	if nil == tokenBody || nil != tokenErr {
		serviceutil.SendBadRequest(w, service.ErrorOprationExpires.Error())
		return
	}
	task.token.RefreshToken(qToken, tokenBody)
	// 用于获取令牌信息
	if len(qOperation) == 0 {
		serviceutil.SendSuccess(w, tokenBody.ToJSON())

		// 用于操作|中断
	} else {
		switch qOperation {
		// 忽略单个 错误
		case "ignore":
			tokenBody.ErrorString = ""
			tokenBody.IsIgnore = true
		// 为后续的 错误 执行忽略
		case "ignoreall":
			tokenBody.ErrorString = ""
			tokenBody.IsIgnoreAll = true
		// 覆盖单个 已存在 错误
		case "replace":
			tokenBody.ErrorString = ""
			tokenBody.IsReplace = true
		// 每次都覆盖 已存在 错误
		case "replaceall":
			tokenBody.ErrorString = ""
			tokenBody.IsReplaceAll = true
		// 立即中断操作
		case "discontinue":
			tokenBody.ErrorString = ""
			tokenBody.IsComplete = true
			tokenBody.IsDiscontinue = true
		default:
			serviceutil.SendBadRequest(w, service.ErrorOprationFailed.Error())
			return
		}
		if err := task.token.RefreshToken(qToken, tokenBody); nil != err {
			serviceutil.SendServerError(w, err.Error())
		} else {
			serviceutil.SendSuccess(w, "")
		}
	}
}

// getTokenObject 获取文件传输Token对象
func (task *CopyFile) getTokenObject(token string) (*CopyFileTokenObject, error) {
	var tokenBody CopyFileTokenObject
	if err := task.token.QueryToken(token, &tokenBody); nil != err {
		return nil, err
	}
	return &tokenBody, nil
}

// GetUserID4Request 获取登录用户
func (task *CopyFile) GetUserID4Request(r *http.Request) string {
	if askstr := task.sg.GetAccessKey4Request(r); len(askstr) > 0 {
		if ack, err := task.sg.GetUserAccess(askstr); nil != err {
			return ack.UserID
		}
	}
	return ""
}
