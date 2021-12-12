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
	"errors"
	"fileservice/biz/service"
	"pakku/ipakku"
)

// AsyncTasks 文件异步操作
type AsyncTasks struct {
	actions map[string]service.AsyncTask
}

// Pakku 模块加载器接口实现, 返回模块信息&配置
func (m *AsyncTasks) Pakku() ipakku.Opts {
	return ipakku.Opts{
		Name:        "AsyncTasks",
		Version:     1.0,
		Description: "异步动作模块",
		OnReady: func(mctx ipakku.Loader) {
			m.AddTaskObject(new(CopyFile).Init(mctx))
			m.AddTaskObject(new(MoveFile).Init(mctx))
		},
		OnSetup: func() {},
		OnInit:  func() {},
	}
}

// GetTaskObject 获取asynctask实例
func (asynctasks *AsyncTasks) GetTaskObject(name string) (service.AsyncTask, error) {
	if len(name) == 0 || nil == asynctasks.actions {
		return nil, errors.New(name + " not supported")
	}
	val, ok := asynctasks.actions[name]
	if ok {
		return val, nil
	}
	return nil, errors.New(name + " not supported")
}

// AddTaskObject 注册asynctask实例
func (asynctasks *AsyncTasks) AddTaskObject(task service.AsyncTask) {
	if nil == asynctasks.actions {
		asynctasks.actions = make(map[string]service.AsyncTask, 0)
	}
	if len(task.Name()) > 0 {
		asynctasks.actions[task.Name()] = task
	}
}
