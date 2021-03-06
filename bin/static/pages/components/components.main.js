// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// VUE页面
; (function (factory) {
	try {
		factory();
	} catch (e) {
		console.error(e);
	}
})(function () {
	"use strict";
	// 文件列表模板
	Vue.component('file-list', {
		template: `
		<div class='file-list' ref='fslist-body' style="height: 100%">
			<div style="padding: 5px 20px 5px 5px;border-bottom: 1px solid #dcdee2;">
				<div>
					<i-button type="text" :icon='fsOperationButtons.refresh.icon' v-show='fsOperationButtons.refresh.show' @click='fsOperationButtons.refresh.handler'>{{fsOperationButtons.refresh.name}}</i-button>
					<i-button type="text" :icon='fsOperationButtons.upload.icon' v-show='fsOperationButtons.upload.show' @click='fsOperationButtons.upload.handler'>{{fsOperationButtons.upload.name}}</i-button>
					<i-button type="text" :icon='fsOperationButtons.newFolder.icon' v-show='fsOperationButtons.newFolder.show' @click='fsOperationButtons.newFolder.handler'>{{fsOperationButtons.newFolder.name}}</i-button>
					<i-button type="text" :icon='fsOperationButtons.download.icon' v-show='fsOperationButtons.download.show' @click='fsOperationButtons.download.handler'>{{fsOperationButtons.download.name}}</i-button>
					<i-button type="text" :icon='fsOperationButtons.rename.icon' v-show='fsOperationButtons.rename.show' @click='fsOperationButtons.rename.handler'>{{fsOperationButtons.rename.name}}</i-button>
					<i-button type="text" :icon='fsOperationButtons.move.icon' v-show='fsOperationButtons.move.show' @click='fsOperationButtons.move.handler'>{{fsOperationButtons.move.name}}</i-button>
					<i-button type="text" :icon='fsOperationButtons.copy.icon' v-show='fsOperationButtons.copy.show' @click='fsOperationButtons.copy.handler'>{{fsOperationButtons.copy.name}}</i-button>
					<i-button type="text" :icon='fsOperationButtons.delete.icon' v-show='fsOperationButtons.delete.show' @click='fsOperationButtons.delete.handler'>{{fsOperationButtons.delete.name}}</i-button>
				</div>
			</div>
			<div style="height: 40px;line-height: 40px;padding: 0px 5px 0px 20px;overflow: hidden;">
				<fs-address :rootname="fsAddress.rootname" :path="fsAddress.loadPath" @click="goToPath" style='float:left;max-width: calc(100% - 250px);'></fs-address>
				<i-Input v-model='searchText' placeholder="Enter file name" @on-enter='loadDatePage(1);' style="float:right;max-width: 240px;padding: 4px 0px;"></i-Input>
			</div>
			<i-table ref="fsSelection" v-minus-height="130" :loading="loading" :columns="fsColumns" :data="fsData" highlight-row="true" @on-row-click="onRowClick" @on-selection-change="onSelectionChange"></i-table>
			<Page show-total :total="fsData_Total" size="small" :current='fsData_CurrentPageIndex' :page-size="fsData_PageSize" style="line-height: 45px;background: #FFF;float: right;padding: 0px 45px;" @on-change="loadDatePage"></Page>
			<!-- 上传文件 -->
			<fs-upload :show-drawer='fsUpload.showDrawer' :parent='fsAddress.loadPath' drag-ref='fslist-body' @on-end="$Message.info('上传完毕');doRefresh();" @on-start="$Message.info('开始上传')" @on-close='fsUpload.showDrawer=false;'></fs-upload>
			<!-- 复制文件 -->
			<fs-copyfile :src-paths="fsCopyFile.srcPaths" :dest-path="fsCopyFile.destPath" :copy-settings="fsCopyFile.copySettings" :show-dailog="fsCopyFile.showDailog" @on-error='onHiddenCopy' @on-stop='onHiddenCopy' @on-end='onHiddenCopy'></fs-copyfile>
			<!-- 移动文件 -->
			<fs-movefile :src-paths='fsMoveFile.srcPaths' :dest-path='fsMoveFile.destPath' :move-settings='fsMoveFile.moveSettings' :show-dailog='fsMoveFile.showDailog' @on-error='onHiddenMove' @on-stop='onHiddenMove' @on-end='onHiddenMove'></fs-movefile>
			<!-- 选择目录 -->
			<fs-selector :select-muti="fsSelector.selectMuti" :select-file="fsSelector.selectFile" :select-dir="fsSelector.selectDir" :start-path="fsSelector.startPath" :show-dailog="fsSelector.showDailog" @on-select="onSelectedFile" @on-cancel="onSelectCancel"></fs-selector>
			<right-click-menu bind-ref='fsSelection' :menus="fsOperationButtons"></right-click-menu>
		</div>
		`,
		data: function () {
			let _ = this;
			return {
				searchText: '',
				fsCopyFile: {
					showDailog: false,
					srcPaths: [],
					destPath: "",
					copySettings: {
						ignore: false,
						replace: false,
					},
				},
				fsMoveFile: {
					showDailog: false,
					srcPaths: [],
					destPath: "",
					moveSettings: {
						ignore: false,
						replace: false,
					},
				},
				fsSelector: {
					operationObj: false,
					showDailog: false,
					selectFile: true,
					selectDir: true,
					selectMuti: true,
					startPath: "/",
				},
				fsAddress: {
					rootPath: "/",
					rootname: "根目录",
					loadPath: ""
				},
				fsUpload: {
					showDrawer: false
				},
				loading: false,
				fsColumns: [
					{
						type: 'selection',
						width: 60,
						align: 'center',
					},
					{
						title: '文件名称',
						key: 'Path',
						render: function (h, params) {
							return h("fs-fileicon", {
								props: {
									node: params.row,
									isEditor: params.row.showRename ? true : false
								},
								on: {
									click: _.doOpen,
									doRename: function (path, name) {
										_.onRenameAfter(params.index, params.row, name);
									}
								}
							});
						}
					},
					{
						title: '修改时间',
						key: 'CtTime',
						maxWidth: 160,
						render: function (h, params) {
							return h("span", $utils.long2Time(params.row.CtTime));
						}
					},
					{
						title: '文件大小',
						key: 'FileSize',
						maxWidth: 160,
						render: function (h, params) {
							if (params.row.IsFile) {
								return h("span", $utils.formatSize(params.row.FileSize));
							} else {
								h("span", "");
							}
						}
					},
					{
						title: '操作权限',
						key: 'PermissionText',
						maxWidth: 280,
					}
				],
				fsData: [],
				fsData_All: [],
				fsData_Total: 0,
				fsData_PageSize: 40,
				fsData_CurrentPageIndex: 1,
				permissionsMap: {},
				fsOperationButtons: {
					refresh: { name: '刷新', show: true, handler: function () { _.doRefresh(); }, icon: 'ivu-icon ivu-icon-md-refresh-circle', divided: false },
					upload: { name: '上传', show: false, handler: function () { _.fsUpload.showDrawer = true; }, icon: 'ivu-icon ivu-icon-md-cloud-upload', divided: false },
					newFolder: { name: '文件夹', show: false, handler: function () { _.onNewFolder(); }, icon: 'ivu-icon ivu-icon-md-add-circle', divided: false },
					download: { name: '下载', show: false, handler: function () { _.doDownload(); }, icon: 'ivu-icon ivu-icon-md-download', divided: false },
					rename: { name: '重命名', show: false, handler: function () { _.onRename(); }, icon: 'ivu-icon ivu-icon-ios-create', divided: false },
					move: { name: '移动', show: false, handler: function () { _.doMove(); }, icon: 'ivu-icon ivu-icon-md-move', divided: false },
					copy: { name: '复制', show: false, handler: function () { _.doCopy(); }, icon: 'ivu-icon ivu-icon-ios-copy', divided: false },
					delete: { name: '删除', show: false, handler: function () { _.onDelete(); }, icon: 'ivu-icon ivu-icon-md-trash', divided: false },
				}
			}
		},
		computed: {
		},
		methods: {
			doRefresh: function () {
				this.doLoadData(this.fsAddress.loadPath);
			},
			doLoadData: function (path) {
				if (this.loading) { return; }
				let _ = this;
				_.loading = true;
				this.onSelectionChange([], {});
				this.fsData_All = []; this.fsData = [];
				$fsApi.List(path).then(function (data) {
					_.fsData_All = JSON.parse(data);
					_.fsData_CurrentPageIndex = 1;
					_.loadDatePage();
					// _.loading = false;
				}).catch(function (err) {
					_.loading = false;
					_.$Message.error(err.toString());
				});
			},
			loadDatePage: function (index) {
				let dataAll = [];
				if (this.searchText) {
					let stemp = this.searchText.toLowerCase();
					for (let i=0; i < this.fsData_All.length; i++ ) {
						if ( this.fsData_All[i].Path.getName().toLowerCase().indexOf( stemp ) > -1 ){
							dataAll.push( this.fsData_All[i] );
						}
					}
				}else{
					dataAll = this.fsData_All;
				}
				this.fsData_Total = dataAll.length;
				this.fsData_CurrentPageIndex = index ? index : 1;
				let _end = this.fsData_CurrentPageIndex * this.fsData_PageSize;
				let _start = (this.fsData_CurrentPageIndex - 1) * this.fsData_PageSize;
				let fsData = dataAll.slice(_start > 0 ? _start : 0, _end);
				this.doWarpPermission(fsData);
			},
			// 包裹权限信息 
			doWarpPermission: function (datas) {
				// if(!datas || datas.length == 0){
				//   this.loading = false; return;
				// }
				this.permissionsMap = {};
				let paths = [this.fsAddress.loadPath];
				datas.forEach(function (row) { paths.push(row.Path); });
				let _ = this;
				$fpmsApi.GetUserPermissionSum($apitools.getSession().UserID, paths).then(function (pms) {
					_.permissionsMap = pms ? JSON.parse(pms) : {};
					for (let i = 0; i < datas.length; i++) {
						datas[i].Permission = undefined != _.permissionsMap[datas[i].Path] ? _.permissionsMap[datas[i].Path] : -1;
						datas[i].PermissionText = $fpmsApi.$TYPE.sum2Name(datas[i].Permission).join(', ');
					}
					_.fsData = datas;
					_.onSelectionChange();
					_.loading = false;
				}).catch(function (err) {
					_.loading = false;
					_.$Message.error(err.toString());
				});
			},
			/* 重命名 */
			doRename: function (index, row, name) {
				row.showRename = false;
				this.$set(this.fsData, index, row);
				if (!name || row.Path.getName() == name) {
					return;
				}
				let _ = this;
				$fsApi.Rename(row.Path, name).then(function () {
					_.$Message.info("操作成功");
					_.doRefresh();
				}).catch(function (err) {
					_.$Message.error("操作失败: " + err.toString());
					_.doRefresh();
				});
			},
			/** 新建文件夹 */
			doNewFolder: function (index, row, name) {
				if (!name) {
					this.doRefresh(); return;
				}
				let _ = this;
				$fsApi.NewFolder(this.fsAddress.loadPath + "/" + name).then(function () {
					_.$Message.info("操作成功");
					_.doRefresh();
				}).catch(function (err) {
					_.$Message.error("操作失败: " + err.toString());
					_.doRefresh();
				});
			},
			onDelete: function () {
				let _ = this;
				this.$Modal.confirm({
					title: '删除文件',
					content: '此操作不可逆, 是否继续删除?',
					onOk: function () {
						_.doDelete();
					}
				});
			},
			/* 删除文件|文件夹 */
			doDelete: function () {
				let select = this.$refs.fsSelection.getSelection();
				if (!select || select.length == 0) {
					return;
				}
				let del_loading = this.$Message.loading({
					content: "正在删除...",
					duration: 0
				});
				let errs = [];
				let _ = this;
				function loop(i) {
					if (!i) { i = 0; }
					if (i == select.length) {
						del_loading();
						_.doRefresh();
						if (errs && errs.length > 0) {
							_.$Message.error({
								content: errs.join("</br>"),
								duration: 5
							});
						}
						return;
					}
					$fsApi.Delete(select[i].Path).then(function () {
						setTimeout(function () {
							loop(++i);
						});
					}).catch(function (err) {
						errs.push(select[i].Path + "删除失败: " + err);
						loop(++i);
					});
				};
				loop(0);
			},
			/* 开始复制 */
			doCopy: function () {
				// 设置源路径
				this.fsCopyFile.srcPaths = this.$refs.fsSelection.getSelection();
				// 文件选择框
				this.fsSelector.operationObj = this.fsCopyFile;
				this.fsSelector.selectFile = false;
				this.fsSelector.selectDir = true;
				this.fsSelector.selectMuti = false;
				this.fsSelector.startPath = this.fsAddress.loadPath;
				this.fsSelector.showDailog = true;
			},
			/* 开始移动 */
			doMove: function () {
				// 设置源路径
				this.fsMoveFile.srcPaths = this.$refs.fsSelection.getSelection();
				// 文件选择框
				this.fsSelector.operationObj = this.fsMoveFile;
				this.fsSelector.selectFile = false;
				this.fsSelector.selectDir = true;
				this.fsSelector.selectMuti = false;
				this.fsSelector.startPath = this.fsAddress.loadPath;
				this.fsSelector.showDailog = true;
			},
			/* 下载单个|多个文件 */
			doDownload: function () {
				let select = this.$refs.fsSelection.getSelection();
				if (!select || select.length == 0) {
					return;
				}
				let _ = this;
				for (let i = select.length - 1; i >= 0; i--) {
					if (!select[i].IsFile) { continue; }
					$fsApi.GetDownloadUrl(select[i].Path).then(function (data) {
						_.$nextTick(function () {
							let download = document.createElement("iframe");
							download.src = data;
							document.body.appendChild(download);
						});
					}).catch(function (err) {
						_.$Message.error("操作失败: " + select[i].Path + ", " + err.toString());
					});
				}
			},
			/* 重命名文件|文件夹 */
			onRename: function () {
				let select = this.$refs.fsSelection.getSelection();
				if (!select || select.length == 0) {
					return;
				}
				select[0].showRename = true;
				for (let i = 0; i < this.fsData.length; i++) {
					if (this.fsData[i].Path == select[0].Path) {
						this.$set(this.fsData, i, select[0]);
					}
				}
			},
			/* 打开文件 */
			doOpen: function (node) {
				if (!node.IsFile) {
					this.goToPath(node.Path)
				} else {
					let _ = this;
					$preview.doPreview(node.Path).catch(function (err) {
						$fsApi.GetSteamUrl(node.Path).then(function (data) {
							_.$nextTick(function () {
								window.open(data)
							});
						}).catch(function (err) {
							_.$Message.error("操作失败: " + node.Path + ", " + err.toString());
						});
					});

				}
			},
			/** 重命名|新建文件夹二合一事件 */
			onRenameAfter: function (index, row, name) {
				if (row.IsNewFolder) {
					this.doNewFolder(index, row, name);
				} else {
					this.doRename(index, row, name);
				}
			},
			/** 新建文件 */
			onNewFolder: function () {
				let temp = this.fsData;
				temp.unshift({
					"Path": "",
					"IsFile": false,
					"IsNewFolder": true,
					"showRename": true,
				});
				this.$set(this, "fsData", temp);
			},
			/* 选择单个文件夹 - OK */
			onSelectedFile: function (rows) {
				this.fsSelector.showDailog = false;
				if (rows && rows[0] && rows[0].Path && rows[0].Path.length > 0) {
					if (this.fsSelector.operationObj) {
						// this.fsCopyFile.destPath = rows[0].Path;
						// this.fsCopyFile.showDailog = true;
						this.fsSelector.operationObj.destPath = rows[0].Path;
						this.fsSelector.operationObj.showDailog = true;
					}
				}
			},
			/* 选择单个文件夹 - Cancel */
			onSelectCancel: function () {
				this.fsSelector.showDailog = false;
			},
			onHiddenCopy: function () {
				this.fsCopyFile.showDailog = false;
			},
			onHiddenMove: function () {
				this.fsMoveFile.showDailog = false;
				this.doRefresh();
			},
			goToPath: function (path) {
				if (this.fsAddress.loadPath != path) {
					this.fsAddress.loadPath = path;
				} else {
					this.doLoadData(path);
				}
			},
			// 点击一行
			onRowClick: function (row, index) {
				for (let i = 0; i < this.fsData.length; i++) {
					if (this.fsData[i]._checked && index != i) {
						this.$set(this.fsData[i], "_checked", false);
					}
				}
				this.$set(this.fsData[index], "_checked", true);
				this.onSelectionChange([row], row);
			},
			// 当Checkbox数据发生变化, 则需要刷新按钮
			onSelectionChange: function (selection, row) {
				// 处理按钮是否显示
				this.fsOperationButtons.upload.show = false;
				this.fsOperationButtons.newFolder.show = false;
				this.fsOperationButtons.download.show = false;
				this.fsOperationButtons.rename.show = false;
				this.fsOperationButtons.delete.show = false;
				this.fsOperationButtons.move.show = false;
				this.fsOperationButtons.copy.show = false;
				let len_selection = selection ? selection.length : 0;
				// 计算上级文件夹权限
				if (undefined != this.permissionsMap[this.fsAddress.loadPath] &&
					$fpmsApi.$TYPE.sumInclude(this.permissionsMap[this.fsAddress.loadPath], $fpmsApi.$TYPE.WRITE.value)) {
					this.fsOperationButtons.upload.show = true;
					this.fsOperationButtons.newFolder.show = true;
				}
				// 选择一个或者以上
				if (len_selection >= 1) {
					let read = false; let fileCount = 0;
					for (let i = 0; i < selection.length; i++) {
						if (selection[i].IsFile) { fileCount++; }
						read = $fpmsApi.$TYPE.sumInclude(selection[i].Permission, $fpmsApi.$TYPE.READ.value);
						if (!read) { read = false; break; }
					}
					if (read) {
						this.fsOperationButtons.download.show = fileCount == selection.length;
						this.fsOperationButtons.copy.show = true;
					}
					// 
					let write = false;
					for (let i = 0; i < selection.length; i++) {
						write = $fpmsApi.$TYPE.sumInclude(selection[i].Permission, $fpmsApi.$TYPE.WRITE.value);
						if (!write) { write = false; break; }
					}
					if (write) {
						this.fsOperationButtons.rename.show = len_selection == 1;
						this.fsOperationButtons.move.show = true;
						this.fsOperationButtons.delete.show = true;
					}
				}
			},
		},
		mounted: function () { },
		created: function () {
			this.goToPath(this.fsAddress.rootPath);
			this.fsOperationButtons.upload.show = true;
			let _ = this;
			this.$nextTick(function () {
				$utils.areaCover({
					background: _.$refs.fsSelection.$el,
					coverfilter: '.ivu-table-row',
					onCover: function (el) {
						if (!el.querySelector('.ivu-checkbox-input').checked) {
							el.querySelector('.ivu-checkbox-input').click();
						}
					},
					onUnCover: function (el) {
						if (el.querySelector('.ivu-checkbox-input').checked) {
							el.querySelector('.ivu-checkbox-input').click();
						}
					}
				})
			})
		},
		watch: {
			'fsAddress.loadPath': function (n, o) {
				this.doLoadData(n);
			}
		}
	});

	// 系统设置
	Vue.component('sys-setting', {
		template: `
		<div class='sys-setting' style="height: 100%;padding-top:7px;">
			<Tabs class="usermanagetab" v-model='currentTab' style='height:100%;'>
				<!-- 账户信息 -->
				<Tab-Pane name='currentAccount' label="账户信息">
				<!-- <div class="mg_0_auto" style="max-width: 300px;"> -->
					<i-Form ref="currentAccount" label-position="left" :model="currentAccount" :label-width="80" style='padding:10px;'>
						<Form-Item class="mg_b_10" label="用户名: " prop="UserName">
							<i-Input v-model="currentAccount.UserName" @on-enter='currentAccount.onUserNameonEnter' placeholder="Enter your name" style="max-width: 200px;"></i-Input>
						</Form-Item>
						<Form-Item class="mg_b_10" label="密码: " prop="UserPWD">
							<a @click='currentAccount.onUpdataPWD' href="javascript:void(0);" >更改密码</a>
						</Form-Item>
						<Form-Item class="mg_b_10" label="用户类型: " prop="usertype">
							{{currentAccount.UserTypeDesc}}
						</Form-Item>
						<Form-Item class="mg_b_10" label="创建时间: " prop="cttime">
							{{currentAccount.CtTime}}
						</Form-Item>
					</i-Form>
				<!-- </div> -->
				</Tab-Pane>
				<!-- 账户管理 -->
				<Tab-Pane v-if="currentAccount.UserType=='1'" name='userManage' label="账户管理">
				<div style="padding: 0px 5px 5px 5px;">
					<div>
						<i-button type="text" icon='md-add-circle' v-show="userManage.shows.add" @click="onAddUser">新增</i-button>
						<i-button type="text" icon='ios-create' v-show="userManage.shows.edit" @click="onEditUser">修改</i-button>
						<i-button type="text" icon='md-trash' v-show="userManage.shows.del" @click="onDelUser">删除</i-button>
					</div>
				</div>
				<i-table ref="userManage" v-minus-height="80" :loading="userManage.loading" :columns="userManage.columns" :data="calcUserManagePageDatas"
				highlight-row="true" @on-row-click="userManage.onRowClick" @on-selection-change="userManage.onSelectionChange" ></i-table>
				<Page show-total :total="userManage.datas.length" size="small" :current='userManage.currentPageIndex' :page-size="userManage.pageSize" style="line-height: 45px;background: #FFF;float: right;padding: 0px 45px;" @on-change="function(val){userManage.currentPageIndex=val;}"></Page>
				</Tab-Pane>
				<!-- 权限管理 -->
				<Tab-Pane v-if="currentAccount.UserType=='1'" name='fPermissionManage' label="文件权限管理">
				<div style="padding: 0px 5px 5px 5px;">
					<div>
						<i-button type="text" icon='md-add-circle' v-show="fPermissionManage.shows.add" @click="onAddFPermission">新增</i-button>
						<i-button type="text" icon='ios-create' v-show="fPermissionManage.shows.edit" @click="onEditFPermission">修改</i-button>
						<i-button type="text" icon='md-trash' v-show="fPermissionManage.shows.del" @click="onDelFPermission">删除</i-button>
					</div>
				</div>
				<i-table ref="fPermissionManage" v-minus-height="80" :loading="fPermissionManage.loading" :columns="fPermissionManage.columns" :data="calcFPermissionManage"
				highlight-row="true" @on-row-click="fPermissionManage.onRowClick" @on-selection-change="fPermissionManage.onSelectionChange" ></i-table>
				<Page show-total :total="fPermissionManage.datas.length" size="small" :current='fPermissionManage.currentPageIndex' :page-size="fPermissionManage.pageSize" style="line-height: 45px;background: #FFF;float: right;padding: 0px 45px;" @on-change="function(val){fPermissionManage.currentPageIndex=val;}"></Page>
				</Tab-Pane>
			</Tabs>
			<!-- 编辑用户 -->
			<Drawer :title="userManage.userEdit.isAdd?'新增用户':'修改用户'" width="450px" :closable="false" v-model="userManage.userEdit.show">
				<i-Form ref="useredit" label-position="left" :model="userManage.userEdit" :label-width="80">
				<Form-Item class="mg_b_10" label="登录ID: " prop="UserID">
					<i-Input v-if='userManage.userEdit.isAdd' v-model="userManage.userEdit.UserID" placeholder="Enter your user ID"></i-Input>
					<span v-else>{{userManage.userEdit.UserID}}</span>
				</Form-Item>
				<Form-Item class="mg_b_10" label="用户名: " prop="UserName">
					<i-Input v-model="userManage.userEdit.UserName" placeholder="Enter your name"></i-Input>
				</Form-Item>
				<Form-Item class="mg_b_10" label="密码: " prop="UserPWD">
					<i-Input v-model="userManage.userEdit.UserPWD" placeholder="Enter your password" type='password'></i-Input>
				</Form-Item>
				<Form-Item class="mg_b_10" label="用户类型: " prop="usertype">
					<!-- {{userManage.userEdit.UserTypeDesc}} --> 普通用户
				</Form-Item>
				</i-Form>
				<div>
				<i-button class="fr mg_l_10" type='primary' @click='onSaveUser'>确定</i-button>
				<i-button class="fr" type='default' @click='userManage.userEdit.show=false'>关闭</i-button>
				</div>
			</Drawer>
			<!-- 编辑权限 -->
			<Drawer :title="fPermissionManage.fPermissionEdit.isAdd?'新增权限':'修改权限'" width="450px" v-model="fPermissionManage.fPermissionEdit.show">
				<i-Form ref="useredit" label-position="left" :model="fPermissionManage.fPermissionEdit" :label-width="80">
				<Form-Item class="mg_b_10" label="用户ID: " prop="UserID">
					<i-select v-if='fPermissionManage.fPermissionEdit.isAdd' v-model="fPermissionManage.fPermissionEdit.UserID" transfer>
						<i-option v-for="item in fPermissionManage.userselector.datas" :value="item.UserID" :key="item.UserID">{{ item.UserName }}</i-option>
					</i-select>
					<span v-else>{{fPermissionManage.fPermissionEdit.UserID}}</span>
				</Form-Item>
				<Form-Item class="mg_b_10" label="授权目录: " prop="Path">
					<i-Input v-if='fPermissionManage.fPermissionEdit.isAdd' v-model="fPermissionManage.fPermissionEdit.Path" placeholder="授权路径" readonly>
						<i-button slot="append" @click="onSelectDir('onBeforFpmsSelectedDir')">选择</i-button>
					</i-Input>
					<span v-else>{{fPermissionManage.fPermissionEdit.Path}}</span>
				</Form-Item>
				<Form-Item class="mg_b_10" label="授权权限: " prop="Permission">
					<!-- <i-Input v-model="fPermissionManage.fPermissionEdit.Permission" placeholder=""></i-Input> -->
					<Checkbox-Group v-model='fPermissionManage.fPermissionEdit.Permissions'>
					<Checkbox label="1">Visible</Checkbox>
					<Checkbox label="2">Read</Checkbox>
					<Checkbox label="3">Write</Checkbox>
					</Checkbox-Group>
				</Form-Item>
				</i-Form>
				<div>
				<i-button class="fr mg_l_10" type='primary' @click='onSaveFPermission'>确定</i-button>
				<i-button class="fr" type='default' @click='fPermissionManage.fPermissionEdit.show=false'>关闭</i-button>
				</div>
			</Drawer>
			<!-- 选择目录 -->
			<fs-selector :select-muti="fsSelector.selectMuti" :select-file="fsSelector.selectFile" :select-dir="fsSelector.selectDir" :start-path="fsSelector.startPath" :show-dailog="fsSelector.showDailog" @on-select="onSelectedFile" @on-cancel="onSelectCancel"></fs-selector>
		</div>
		`,
		data: function () {
			let _ = this;
			return {
				UserTypes: {
					'1': '管理员',
					'0': '普通用户',
				},
				fsSelector: {
					operationObj: false,
					showDailog: false,
					selectFile: true,
					selectDir: true,
					selectMuti: true,
					startPath: "/",
				},
				currentTab: 'currentAccount',
				// 账户管理
				userManage: {
					loading: false,
					datas: [],
					currentPageIndex: 1,
					pageSize: 60,
					shows: {
						add: true,
						del: false,
						edit: false,
					},
					userEdit: {
						show: false,
						isAdd: false,
						UserID: '',
						UserType: '',
						UserName: '',
						UserPWD: '',
					},
					columns: [
						{
							type: 'selection',
							width: 60,
							align: 'center',
						},
						{
							title: '用户名',
							key: 'UserName',
						},
						{
							title: '用户类型',
							key: 'UserTypeDesc',
						},
						{
							title: '新建时间',
							key: 'CtTime',
						},
					],
					// 点击一行
					onRowClick: function (row, index) {
						_.$refs.userManage.selectAll(false);
						_.$refs.userManage.toggleSelect(index);
						_.userManage.onSelectionChange([row], row);
					},
					// 当Checkbox数据发生变化, 则需要刷新按钮
					onSelectionChange: function (selection, row) {
						_.userManage.shows.add = true;
						_.userManage.shows.del = false;
						_.userManage.shows.edit = false;
						let len_selection = selection.length;
						// 选择一个或者以上
						if (len_selection >= 1) {
							_.userManage.shows.del = true;
						}
						// 只选择了一个
						if (len_selection == 1) {
							_.userManage.shows.edit = true;
						}
					},
				},
				// 文件权限管理
				fPermissionManage: {
					loading: false,
					datas: [],
					currentPageIndex: 1,
					pageSize: 60,
					shows: {
						add: true,
						del: false,
						edit: false,
					},
					permissionMax: 3, // 权限最大数
					fPermissionEdit: {
						show: false,
						isAdd: false,
						UserID: '',
						Path: '',
						Permission: 0,
						Permissions: [],
						PermissionID: '',
					},
					// 临时用下拉-后期优化
					userselector: {
						datas: []
					},
					columns: [
						{
							type: 'selection',
							width: 60,
							align: 'center',
						},
						{
							title: '用户ID',
							key: 'UserID',
						},
						{
							title: '文件路径',
							key: 'Path',
						},
						{
							title: '权限类型',
							key: 'Permission',
							render: function (h, c) {
								return h('span', {}, $fpmsApi.$TYPE.sum2Name(c.row.Permission).join(', '));
							},
						},
						{
							title: '创建时间',
							key: 'CtTime',
						},
					],
					// 点击一行
					onRowClick: function (row, index) {
						_.$refs.fPermissionManage.selectAll(false);
						_.$refs.fPermissionManage.toggleSelect(index);
						_.fPermissionManage.onSelectionChange([row], row);
					},
					// 当Checkbox数据发生变化, 则需要刷新按钮
					onSelectionChange: function (selection, row) {
						_.fPermissionManage.shows.add = true;
						_.fPermissionManage.shows.del = false;
						_.fPermissionManage.shows.edit = false;
						let len_selection = selection.length;
						// 选择一个或者以上
						if (len_selection >= 1) {
							_.fPermissionManage.shows.del = true;
						}
						// 只选择了一个
						if (len_selection == 1) {
							_.fPermissionManage.shows.edit = true;
						}
					},

				},
				// 账户信息
				currentAccount: {
					UserID: '',
					UserType: '',
					UserTypeDesc: '',
					UserName: '',
					UserPWD: '',
					CtTime: '',
					ruleValidate: {
						UserName: [{ required: true },],
					},
					onUserNameonEnter: function () {
						if (!_.currentAccount.UserName) {
							_.$Message.error('不能为空'); return;
						}
						$userApi.updateUserName(_.currentAccount.UserID, _.currentAccount.UserName).then(function (data) {
							_.$Message.info("操作成功");
							_.$root.doInitCurrentAccountInfo();
							// _.doBuildAccountInfo();
						}).catch(function (err) {
							_.$Message.error("操作失败");
							_.$root.doInitCurrentAccountInfo();
							// _.doBuildAccountInfo();
						});
					},
					onUpdataPWD: function () {
						let newpwd = '';
						_.$Modal.info({
							title: '更改密码',
							render: function (h) {
								return h('Input', {
									props: {
										type: 'password',
										value: '',
										autofocus: true,
										placeholder: 'Please enter your password...'
									},
									on: {
										input: function (val) {
											newpwd = val;
										}
									}
								})
							},
							onOk: function () {
								$userApi.updateUserPwd(_.currentAccount.UserID, newpwd).then(function (data) {
									_.$Message.info('操作成功');
								}).catch(function (err) {
									_.$Message.error("操作失败");
								});
							},
						})
					},
				},
			};
		},
		computed: {
			calcUserManagePageDatas: function () {
				let _end = this.userManage.currentPageIndex * this.userManage.pageSize;
				let _start = (this.userManage.currentPageIndex - 1) * this.userManage.pageSize;
				return this.userManage.datas.slice(_start > 0 ? _start : 0, _end);
			},
			calcFPermissionManage: function () {
				let _end = this.fPermissionManage.currentPageIndex * this.fPermissionManage.pageSize;
				let _start = (this.fPermissionManage.currentPageIndex - 1) * this.fPermissionManage.pageSize;
				return this.fPermissionManage.datas.slice(_start > 0 ? _start : 0, _end);
			},
		},
		methods: {
			// 选择路径
			onSelectDir: function (cb) {
				this.fsSelector.operationObj = cb;
				this.fsSelector.selectFile = false;
				this.fsSelector.selectDir = true;
				this.fsSelector.selectMuti = false;
				this.fsSelector.startPath = "/";
				this.fsSelector.showDailog = true;
			},
			onSelectCancel: function () {
				this.fsSelector.showDailog = false;
			},
			onSelectedFile: function (rows) {
				this.fsSelector.showDailog = false;
				if (this.fsSelector.operationObj) {
					this[this.fsSelector.operationObj](rows);
				}
			},
			// 加载账户信息
			doBuildAccountInfo: function (data) {
				if (data) {
					data.UserTypeDesc = this.UserTypes[data.UserType] ? this.UserTypes[data.UserType] : '未知类型';
				}
			},
			// doListAllUsers
			doListAllUsers: function () {
				let _ = this;
				_.userManage.datas = [];
				_.userManage.loading = true;
				_.userManage.currentPageIndex = 1;
				$userApi.listAllUsers().then(function (datas) {
					if (datas) {
						datas = JSON.parse(datas);
						for (let i = 0; i < datas.length; i++) {
							_.doBuildAccountInfo(datas[i]);
						}
					}
					_.$set(_.userManage, 'datas', datas);
					_.userManage.loading = false;
				}).catch(function (err) {
					_.$Message.error(err.toString());
					_.userManage.loading = false;
				})
			},
			onAddUser: function () {
				this.userManage.userEdit.isAdd = true;
				this.userManage.userEdit.UserID = '';
				this.userManage.userEdit.UserName = '';
				this.userManage.userEdit.UserType = '';
				this.userManage.userEdit.UserPWD = '';
				this.userManage.userEdit.show = true;
			},
			onDelUser: function () {
				let rows = this.$refs.userManage.getSelection();
				if (rows) {
					try {
						for (let i = 0; i < rows.length; i++) {
							$userApi.sync.delUser(rows[i].UserID)
						}
						this.$Message.info('操作成功');
					} catch (err) {
						this.$Message.error(err.toString());
					}
					this.doListAllUsers();
				}
			},
			onEditUser: function () {
				let rows = this.$refs.userManage.getSelection();
				if (rows && rows.length == 1) {
					this.userManage.userEdit.isAdd = false;
					this.userManage.userEdit.UserPWD = '';
					this.userManage.userEdit.UserID = rows[0].UserID;
					this.userManage.userEdit.UserName = rows[0].UserName;
					this.userManage.userEdit.UserType = rows[0].UserType;
					this.userManage.userEdit.show = true;
				}
			},
			onSaveUser: function () {
				if (!this.userManage.userEdit.UserID) {
					this.$Message.error('用户ID不能为空'); return;
				}
				if (!this.userManage.userEdit.UserName) {
					this.$Message.error('用户名不能为空'); return;
				}
				if (this.userManage.userEdit.isAdd) {
					let _ = this;
					$userApi.addUser(this.userManage.userEdit.UserID, this.userManage.userEdit.UserName, this.userManage.userEdit.UserPWD).then(function (data) {
						_.$Message.info('操作成功');
						_.doListAllUsers();
						_.userManage.userEdit.show = false;
					}).catch(function (err) {
						_.$Message.error(err.toString());
						_.doListAllUsers();
					});
				} else {
					try {
						$userApi.sync.updateUserName(this.userManage.userEdit.UserID, this.userManage.userEdit.UserName);
						$userApi.sync.updateUserPwd(this.userManage.userEdit.UserID, this.userManage.userEdit.UserPWD);
						this.userManage.userEdit.show = false;
						this.$Message.info('操作成功');
					} catch (err) {
						this.$Message.error(err.toString());
					}
					this.doListAllUsers();
				}
			},
			// 加载权限
			doListFPermissions: function () {
				let _ = this;
				_.fPermissionManage.datas = [];
				_.fPermissionManage.loading = true;
				_.fPermissionManage.currentPageIndex = 1;
				$fpmsApi.listFPermissions().then(function (datas) {
					if (datas) {
						datas = JSON.parse(datas);
					}
					_.$set(_.fPermissionManage, 'datas', datas);
					_.fPermissionManage.loading = false;
				}).catch(function (err) {
					_.$Message.error(err.toString());
					_.fPermissionManage.loading = false;
				})
			},
			// 添加权限
			onAddFPermission: function () {
				// 后期优化
				let _ = this;
				$userApi.listAllUsers().then(function (datas) {
					if (datas) {
						_.fPermissionManage.userselector.datas = JSON.parse(datas);
					}
				}).catch(console.error)

				this.fPermissionManage.fPermissionEdit.isAdd = true;
				this.fPermissionManage.fPermissionEdit.UserID = '';
				this.fPermissionManage.fPermissionEdit.Path = '';
				this.fPermissionManage.fPermissionEdit.Permission = 0;
				this.fPermissionManage.fPermissionEdit.Permissions = [];
				this.fPermissionManage.fPermissionEdit.PermissionID = '';
				this.fPermissionManage.fPermissionEdit.show = true;
			},
			// 
			onBeforFpmsSelectedDir: function (rows) {
				if (rows && rows.length > 0) {
					this.fPermissionManage.fPermissionEdit.Path = rows[0].Path;
				}
			},
			// 修改权限
			onEditFPermission: function () {
				let rows = this.$refs.fPermissionManage.getSelection();
				if (rows && rows.length == 1) {
					this.fPermissionManage.fPermissionEdit.isAdd = false;
					this.fPermissionManage.fPermissionEdit.UserID = rows[0].UserID;
					this.fPermissionManage.fPermissionEdit.Path = rows[0].Path;
					this.fPermissionManage.fPermissionEdit.Permission = parseInt(rows[0].Permission);
					this.fPermissionManage.fPermissionEdit.Permissions = (function (_) {
						let permission = _.fPermissionManage.fPermissionEdit.Permission;
						let res = [];
						if (permission) {
							for (let i = 0; i <= _.fPermissionManage.permissionMax; i++) {
								if (1 << i == (permission & (1 << i))) {
									res.push("" + i);
								}
							}
						}
						return res;
					})(this);
					this.fPermissionManage.fPermissionEdit.PermissionID = rows[0].PermissionID;
					this.fPermissionManage.fPermissionEdit.show = true;
				}
			},
			// 删除权限
			onDelFPermission: function () {
				let rows = this.$refs.fPermissionManage.getSelection();
				if (rows) {
					try {
						for (let i = 0; i < rows.length; i++) {
							$fpmsApi.sync.delFPermission(rows[i].PermissionID)
						}
						this.$Message.info('操作成功');
					} catch (err) {
						this.$Message.error(err.toString());
					}
					this.doListFPermissions();
				}
			},
			// 保存
			onSaveFPermission: function () {
				if (!this.fPermissionManage.fPermissionEdit.UserID) {
					this.$Message.error('用户ID不能为空'); return;
				}
				if (!this.fPermissionManage.fPermissionEdit.Path) {
					this.$Message.error('授权路径名不能为空'); return;
				}
				let _ = this;
				this.fPermissionManage.fPermissionEdit.Permission = (function () {
					let res = 0;
					if (_.fPermissionManage.fPermissionEdit.Permissions.length > 0) {
						for (let i = 0; i < _.fPermissionManage.fPermissionEdit.Permissions.length; i++) {
							res += (1 << _.fPermissionManage.fPermissionEdit.Permissions[i]);
						}
					}
					return res;
				})();
				if (this.fPermissionManage.fPermissionEdit.isAdd) {
					$fpmsApi.addFPermission(this.fPermissionManage.fPermissionEdit.UserID, this.fPermissionManage.fPermissionEdit.Path, this.fPermissionManage.fPermissionEdit.Permission).then(function (data) {
						_.$Message.info('操作成功');
						_.doListFPermissions();
						_.fPermissionManage.fPermissionEdit.show = false;
					}).catch(function (err) {
						_.$Message.error(err.toString());
						_.doListFPermissions();
					});
				} else {
					$fpmsApi.updateFPermission(this.fPermissionManage.fPermissionEdit.PermissionID, this.fPermissionManage.fPermissionEdit.Permission).then(function (data) {
						_.$Message.info('操作成功');
						_.doListFPermissions();
						_.fPermissionManage.fPermissionEdit.show = false;
					}).catch(function (err) {
						_.$Message.error(err.toString());
						_.doListFPermissions();
					});
				}
			},
		},
		created: function () {
			this.doBuildAccountInfo(this.$root.$data.currentAccount);
			this.$set(this, 'currentAccount', $utils.extendAttrs(this.currentAccount, this.$root.$data.currentAccount));
		},
		watch: {
			currentTab: function (n, o) {
				if (n == 'userManage') {
					this.doListAllUsers();
				} else if (n == 'fPermissionManage') {
					this.doListFPermissions();
				}
			}
		}
	});

});