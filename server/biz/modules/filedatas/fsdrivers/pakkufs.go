// Copyright (C) 2022 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// pakkufs挂载操作驱动

package fsdrivers

import (
	"errors"
	"fileservice/biz/modules/filedatas/ifiledatas"
	"io"
	"opensdk"
	"pakku/utils/fileutil"
	"pakku/utils/logs"
	"pakku/utils/strutil"
	"path/filepath"
	"strings"
)

// PakkuFsDriver 本地文件挂载操作驱动
type PakkuFsDriver struct {
	mtn *ifiledatas.MountNode
	mtm ifiledatas.DIRMount
	sdk opensdk.IOpenApi
}

// GetDriverType 当驱动注册时调用
func (driver PakkuFsDriver) GetDriverType() string {
	return "PAKKUFS"
}

// InstanceDriver 当驱动实例化时调用
func (driver PakkuFsDriver) InstanceDriver(dirMount ifiledatas.DIRMount, mtnode *ifiledatas.MountNode) (ifiledatas.FileDriver, error) {
	authInfo := strings.Split(mtnode.Addr, "@")
	if len(authInfo) != 2 {
		return nil, errors.New("PAKKUFS连接地址不正确, 如: OPENAPI@192.168.2.201/mountdir")
	}
	rpcAddr := authInfo[1]
	pathIndex := strings.Index(rpcAddr, "/")
	if pathIndex > -1 {
		rpcAddr = authInfo[1][:pathIndex]
		mtnode.Addr = strutil.Parse2UnixPath(authInfo[1][pathIndex:])
	} else {
		mtnode.Addr = "/"
	}
	logs.Infoln("PAKKUFS: ", rpcAddr, authInfo[0], mtnode.Addr)
	// 实例化
	instance := &PakkuFsDriver{
		mtm: dirMount,
		mtn: mtnode,
		sdk: opensdk.NewOpenApi(opensdk.Conn4RPC{
			RPCAddr: rpcAddr,
			User: &opensdk.User{
				User:   authInfo[0],
				Passwd: mtnode.Passwd,
			},
		}),
	}
	// endpoint
	if val, ok := mtnode.Props["endpoints"]; ok {
		if val, ok := val.(map[string]interface{}); ok {
			endpoints := make(map[string]string)
			for k, v := range val {
				endpoints[k] = v.(string)
			}
			instance.sdk.SetDataNodeDNS(endpoints)
			logs.Infoln("EndPoints: ", endpoints)
		}
	}
	return instance, nil
}

// IsExist 文件是否存在
func (driver *PakkuFsDriver) IsExist(relativePath string) bool {
	if absPath, _, err := driver.getAbsolutePath(driver.mtn, relativePath); nil != err {
		logs.Errorln(err)
	} else {
		if ok, err := driver.sdk.IsExist(absPath); nil != err {
			logs.Errorln(err)
		} else {
			return ok
		}
	}
	return false
}

// IsDir IsDir
func (driver *PakkuFsDriver) IsDir(relativePath string) bool {
	if absPath, _, err := driver.getAbsolutePath(driver.mtn, relativePath); nil != err {
		logs.Errorln(err)
	} else {
		if ok, err := driver.sdk.IsDir(absPath); nil == err {
			return ok
		} else {
			logs.Errorln(err)
		}
	}
	return false
}

// IsFile IsFile
func (driver *PakkuFsDriver) IsFile(relativePath string) bool {
	if absPath, _, err := driver.getAbsolutePath(driver.mtn, relativePath); nil != err {
		logs.Errorln(err)
	} else {
		if ok, err := driver.sdk.IsFile(absPath); nil == err {
			return ok
		} else {
			logs.Errorln(err)
		}
	}
	return false
}

// GetNode GetNode
func (driver *PakkuFsDriver) GetNode(src string) *ifiledatas.Node {
	if absPath, _, err := driver.getAbsolutePath(driver.mtn, src); nil != err {
		logs.Errorln(err)
	} else {
		if node, err := driver.sdk.GetNode(absPath); nil == err && node != nil {
			return &ifiledatas.Node{
				Path:   src,
				Mtime:  node.Mtime,
				IsFile: node.Flag == 1,
				IsDir:  node.Flag == 0,
				Size:   node.Size,
			}
		}
	}
	return nil
}

// GetNodes GetNodes
func (driver *PakkuFsDriver) GetNodes(src []string, ignoreNotIsExist bool) (result []ifiledatas.Node, err error) {
	for i := 0; i < len(src); i++ {
		if absPath, _, err := driver.getAbsolutePath(driver.mtn, src[i]); nil != err {
			return result, err
		} else {
			src[i] = absPath
		}
	}
	if nodes, err := driver.sdk.GetNodes(src); nil == err && len(nodes) > 0 {
		result = make([]ifiledatas.Node, len(nodes))
		for i := 0; i < len(nodes); i++ {
			result[i] = ifiledatas.Node{
				// Path:   src, // TODO Path
				Mtime:  nodes[i].Mtime,
				IsFile: nodes[i].Flag == 1,
				IsDir:  nodes[i].Flag == 0,
				Size:   nodes[i].Size,
			}
		}
	} else {
		return result, err
	}
	return result, err
}

// GetDirNodeList 获取node列表
func (driver *PakkuFsDriver) GetDirNodeList(src string, limit, offset int) (result []ifiledatas.Node, err error) {
	src = strutil.Parse2UnixPath(src)
	if absPath, _, err := driver.getAbsolutePath(driver.mtn, src); err == nil {
		if dto, err := driver.sdk.GetDirNodeList(absPath, limit, offset); nil == err && nil != dto && dto.Total > 0 {
			result = make([]ifiledatas.Node, len(dto.Datas))
			for i := 0; i < len(dto.Datas); i++ {
				result[i] = ifiledatas.Node{
					Path:   src + "/" + dto.Datas[i].Name,
					Mtime:  dto.Datas[i].Mtime,
					IsFile: dto.Datas[i].Flag == 1,
					IsDir:  dto.Datas[i].Flag == 0,
					Size:   dto.Datas[i].Size,
				}
			}
		} else {
			return result, err
		}
	} else {
		return result, err
	}
	return result, err
}

// GetDirList 获取路径列表, 返回相对路径
func (driver *PakkuFsDriver) GetDirList(relativePath string, limit, offset int) []string {
	if absPath, _, err := driver.getAbsolutePath(driver.mtn, relativePath); err == nil {
		if ls, err := driver.sdk.GetDirNameList(absPath, limit, offset); nil == err {
			return ls.Datas
		}
	}
	return make([]string, 0)
}

// GetFileSize GetFileSize
func (driver *PakkuFsDriver) GetFileSize(relativePath string) int64 {
	if absPath, _, err := driver.getAbsolutePath(driver.mtn, relativePath); nil == err {
		if size, err := driver.sdk.GetFileSize(absPath); nil == err {
			return size
		}
	}
	return -1
}

// GetModifyTime GetModifyTime
func (driver *PakkuFsDriver) GetModifyTime(relativePath string) int64 {
	if absPath, _, err := driver.getAbsolutePath(driver.mtn, relativePath); nil == err {
		if node, err := driver.sdk.GetNode(absPath); nil == err && nil != node {
			return node.Mtime
		}
	}
	return -1
}

// DoMove 移动文件|夹
func (driver *PakkuFsDriver) DoMove(src string, dst string, replace bool) error {
	if driver.mtn.Path == src {
		return errors.New("Does not allow access: " + src)
	}
	absSrc, _, err := driver.getAbsolutePath(driver.mtn, src)
	if nil != err {
		return err
	} else if driver.mtn.Addr == absSrc {
		return errors.New(src + " is mount root, cannot move")
	}
	// 目标位置驱动接口
	dstMountItem := driver.mtm.GetMountNode(dst)
	switch dstMountItem.Type {
	case driver.GetDriverType():
		{ // 相同类型存储
			if absDst, _, err := driver.getAbsolutePath(dstMountItem, dst); nil != err {
				return err
			} else {
				return driver.sdk.DoMove(absSrc, absDst, replace)
			}
		}
	default:
		{ // 不支持的分区挂载类型
			return errors.New("Unsupported partition mount type: " + dstMountItem.Type)
		}
	}
}

// DoRename 重命名文件|文件夹
func (driver *PakkuFsDriver) DoRename(relativePath string, newName string) error {
	if driver.mtn.Path == relativePath {
		return errors.New("Does not allow access: " + relativePath)
	}
	if absSrc, _, err := driver.getAbsolutePath(driver.mtn, relativePath); nil != err {
		return err
	} else {
		if len(newName) > 0 {
			return driver.sdk.DoRename(absSrc, newName)
		}
	}
	return nil
}

// DoMkDir 新建文件夹
func (driver *PakkuFsDriver) DoMkDir(relativePath string) (err error) {
	if driver.mtn.Path == relativePath {
		return errors.New("Does not allow access: " + relativePath)
	}
	if absSrc, _, e := driver.getAbsolutePath(driver.mtn, relativePath); nil != e {
		return e
	} else {
		_, err = driver.sdk.DoMkDir(absSrc)
	}
	return err
}

// DoDelete 删除文件|文件夹
func (driver *PakkuFsDriver) DoDelete(relativePath string) error {
	if driver.mtn.Path == relativePath {
		return errors.New("Does not allow access: " + relativePath)
	}
	if absSrc, _, err := driver.getAbsolutePath(driver.mtn, relativePath); nil != err {
		return err
	} else {
		return driver.sdk.DoDelete(absSrc)
	}
}

// DoCopy 拷贝文件
func (driver *PakkuFsDriver) DoCopy(src, dst string, replace bool) error {
	absSrc, _, err := driver.getAbsolutePath(driver.mtn, src)
	if nil != err {
		return err
	}
	// 目标位置驱动接口
	dstMountItem := driver.mtm.GetMountNode(dst)
	switch dstMountItem.Type {
	case driver.GetDriverType():
		{ // 相同存储类型
			if absDst, _, err := driver.getAbsolutePath(dstMountItem, dst); nil != err {
				return err
			} else {
				return driver.copyNodes(absSrc, absDst, replace)
			}
		}
	default:
		{ // 不支持的分区挂载类型
			return errors.New("Unsupported partition mount type: " + dstMountItem.Type)
		}
	}
}

// copyNodes 复制文件|文件夹
func (driver *PakkuFsDriver) copyNodes(src, dest string, replace bool) (err error) {
	if len(src) == 0 || len(dest) == 0 {
		return errors.New("src or dest is nil")
	}
	var srcNode *opensdk.TNode
	if srcNode, err = driver.sdk.GetNode(src); nil == err {
		if srcNode.Flag == 0 {
			offset := 0
			var nodes *opensdk.DirNodeListDto
			for nil == err {
				if nodes, err = driver.sdk.GetDirNodeList(src, 100, offset); nil != err {
					break
				}
				if nil == nodes || nodes.Total == 0 || len(nodes.Datas) == 0 {
					break
				}
				for i := 0; i < len(nodes.Datas); i++ {
					node := nodes.Datas[i]
					if node.Flag == 0 {
						if err = driver.copyNodes(strutil.Parse2UnixPath(src+"/"+node.Name), strutil.Parse2UnixPath(dest+"/"+node.Name), replace); nil != err {
							break
						}
					} else {
						src := strutil.Parse2UnixPath(src + "/" + node.Name)
						dest := strutil.Parse2UnixPath(dest + "/" + node.Name)
						if _, err = driver.sdk.DoCopy(src, dest, replace); nil != err {
							break
						}
					}
				}
				offset += 100
				if nil != nodes && nodes.Total <= offset {
					break
				}
			}
		} else {
			_, err = driver.sdk.DoCopy(src, dest, replace)
		}
	}
	return err
}

// DoRead 读取文件, 需要手动关闭流
func (driver *PakkuFsDriver) DoRead(relativePath string, offset int64) (io.ReadCloser, error) {
	return nil, errors.New("unsupported")
}

// DoWrite 写入文件， 先写入临时位置, 然后移动到正确位置
func (driver *PakkuFsDriver) DoWrite(relativePath string, ioReader io.Reader) error {
	if ioReader == nil {
		return driver.wrapError(relativePath, "", errors.New("IO Reader is nil"))
	}
	return errors.New("unsupported")
}

// DoAskAccessToken 申请访问Token
func (driver *PakkuFsDriver) DoAskAccessToken(src string, tokenType ifiledatas.AccessTokenType, props map[string]interface{}) (*ifiledatas.AccessToken, error) {
	absSrc, _, err := driver.getAbsolutePath(driver.mtn, src)
	if nil != err {
		return nil, err
	}
	if tokenType == ifiledatas.AccessTokenType_Read {
		if ok, err := driver.sdk.IsExist(absSrc); !ok || err != nil {
			return nil, fileutil.PathNotExist("AskAccessToken", src)
		}
		if token, err := driver.sdk.DoAskReadToken(absSrc); nil == err {
			if url, err := driver.sdk.GetReadStreamURL(token.NodeNo, token.Token, token.EndPoint); nil == err {
				if len(props) > 0 {
					if val, ok := props["name"]; ok && len(val.(string)) > 0 {
						url += "?name=" + val.(string)
					}
				}
				return &ifiledatas.AccessToken{
					Token:    token.Token,
					CTime:    token.CTime,
					TokenURL: url,
				}, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else if tokenType == ifiledatas.AccessTokenType_Write {
		if token, err := driver.sdk.DoAskWriteToken(absSrc); nil == err {
			if url, err := driver.sdk.GetWriteStreamURL(token.NodeNo, token.Token, token.EndPoint); nil == err {
				return &ifiledatas.AccessToken{
					Token:    token.Token,
					CTime:    token.CTime,
					TokenURL: url,
				}, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("unsupported token type")
	}
}

// DoSubmitToken 递交令牌
func (driver *PakkuFsDriver) DoSubmitToken(token string, props map[string]interface{}) (*ifiledatas.Node, error) {
	override := false
	if len(props) > 0 {
		if val, ok := props["override"]; ok {
			if nvl, ok := val.(bool); ok {
				override = nvl
			} else if nvl, ok := val.(string); ok {
				override = strutil.String2Bool(nvl)
			}
		}
	}
	if node, err := driver.sdk.DoSubmitWriteToken(token, override); nil == err {
		return &ifiledatas.Node{
			// Path:   "", // TODO Path
			Mtime:  node.Mtime,
			IsFile: node.Flag == 1,
			IsDir:  node.Flag == 0,
			Size:   node.Size,
		}, nil

	} else {
		return nil, err
	}
}

// getAbsolutePath (mnode, 虚拟路径)(绝对位置, 挂载位置, 错误)处理路径拼接
func (driver *PakkuFsDriver) getAbsolutePath(mtn *ifiledatas.MountNode, relativePath string) (abs string, rlPath string, err error) {
	rlPath = relativePath
	if mtn.Path != "/" {
		rlPath = relativePath[len(mtn.Path):]
		if rlPath == "" {
			rlPath = "/"
		}
	}
	abs = strutil.Parse2UnixPath(mtn.Addr + rlPath)
	return
}

// wrapLocalError 重新包装本地驱动错误信息, 避免真实路径暴露
func (driver *PakkuFsDriver) wrapError(src, dst string, err error) error {
	if nil != err {
		return errors.New(driver.clearMountAddr(driver.mtm.GetMountNode(src), driver.mtm.GetMountNode(dst), err))
	}
	return nil
}

// clearMountAddr 去除挂载目录的位置信息
func (driver *PakkuFsDriver) clearMountAddr(src, dst *ifiledatas.MountNode, err error) string {
	if nil != err {
		errStr := err.Error()
		// src
		srcAddr := ""
		if nil != src {
			srcAddr = filepath.Clean(src.Addr)
			if strings.Contains(errStr, srcAddr) {
				// /root/datas/a/b -> /a/b/a/b
				rStr := src.Path
				if src.Path == "/" {
					rStr = ""
				}
				errStr = strings.Replace(errStr, srcAddr, rStr, -1)
			}
		}
		// dst
		if nil != dst {
			dstAddr := filepath.Clean(dst.Addr)
			if dstAddr != srcAddr && strings.Contains(errStr, dstAddr) {
				rStr := dst.Path
				if dst.Path == "/" {
					rStr = ""
				}
				errStr = strings.Replace(errStr, dstAddr, rStr, -1)
			}
		}
		return errStr
	}
	return ""
}
