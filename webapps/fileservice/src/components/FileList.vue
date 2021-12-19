<template>
  <div class="page" ref="main-page">
    <div class="file-list" ref="fslist-body" style="height: 100%">
      <div style="padding: 5px 20px 5px 5px; border-bottom: 1px solid #dcdee2">
        <div>
          <Button
            type="text"
            :icon="fsOperationButtons.refresh.icon"
            v-show="fsOperationButtons.refresh.show"
            @click="fsOperationButtons.refresh.handler"
            >{{ fsOperationButtons.refresh.name }}</Button
          >
          <Button
            type="text"
            :icon="fsOperationButtons.upload.icon"
            v-show="fsOperationButtons.upload.show"
            @click="fsOperationButtons.upload.handler"
            >{{ fsOperationButtons.upload.name }}</Button
          >
          <Button
            type="text"
            :icon="fsOperationButtons.newFolder.icon"
            v-show="fsOperationButtons.newFolder.show"
            @click="fsOperationButtons.newFolder.handler"
            >{{ fsOperationButtons.newFolder.name }}</Button
          >
          <Button
            type="text"
            :icon="fsOperationButtons.download.icon"
            v-show="fsOperationButtons.download.show"
            @click="fsOperationButtons.download.handler"
            >{{ fsOperationButtons.download.name }}</Button
          >
          <Button
            type="text"
            :icon="fsOperationButtons.rename.icon"
            v-show="fsOperationButtons.rename.show"
            @click="fsOperationButtons.rename.handler"
            >{{ fsOperationButtons.rename.name }}</Button
          >
          <Button
            type="text"
            :icon="fsOperationButtons.move.icon"
            v-show="fsOperationButtons.move.show"
            @click="fsOperationButtons.move.handler"
            >{{ fsOperationButtons.move.name }}</Button
          >
          <Button
            type="text"
            :icon="fsOperationButtons.copy.icon"
            v-show="fsOperationButtons.copy.show"
            @click="fsOperationButtons.copy.handler"
            >{{ fsOperationButtons.copy.name }}</Button
          >
          <Button
            type="text"
            :icon="fsOperationButtons.delete.icon"
            v-show="fsOperationButtons.delete.show"
            @click="fsOperationButtons.delete.handler"
            >{{ fsOperationButtons.delete.name }}</Button
          >
        </div>
      </div>
      <div
        style="
          height: 40px;
          line-height: 40px;
          padding: 0px 5px 0px 20px;
          overflow: hidden;
        "
      >
        <fileAddress
          :rootname="fsAddress.rootname"
          :path="fsAddress.loadPath"
          @click="goToPath"
          style="float: left; max-width: calc(100% - 250px)"
        ></fileAddress>
        <i-input
          v-model="searchText"
          placeholder="Enter file name"
          @on-enter="loadDatePage(1)"
          style="float: right; max-width: 240px; padding: 4px 0px"
        ></i-input>
      </div>
      <Table
        ref="fsSelection"
        v-minus-height="130"
        :loading="loading"
        :columns="fsColumns"
        :data="fsData"
        :highlight-row="true"
        @on-row-click="onRowClick"
        @on-selection-change="onSelectionChange"
      ></Table>
      <Page
        show-total
        :total="fsData_Total"
        size="small"
        :current="fsData_CurrentPageIndex"
        :page-size="fsData_PageSize"
        style="
          line-height: 45px;
          background: #fff;
          float: right;
          padding: 0px 45px;
        "
        @on-change="loadDatePage"
      ></Page>
      <!-- 上传文件 -->
      <uploadfile
        :show-drawer="fsUpload.showDrawer"
        :parent="fsAddress.loadPath"
        drag-ref="fslist-body"
        @on-end="
          $Message.info('上传完毕');
          doRefresh();
        "
        @on-start="$Message.info('开始上传')"
        @on-close="fsUpload.showDrawer = false"
      ></uploadfile>
      <!-- 复制文件 -->
      <copyfile
        :src-paths="fsCopyFile.srcPaths"
        :dest-path="fsCopyFile.destPath"
        :copy-settings="fsCopyFile.copySettings"
        :show-dailog="fsCopyFile.showDailog"
        @on-error="onHiddenCopy"
        @on-stop="onHiddenCopy"
        @on-end="onHiddenCopy"
      ></copyfile>
      <!-- 移动文件 -->
      <movefile
        :src-paths="fsMoveFile.srcPaths"
        :dest-path="fsMoveFile.destPath"
        :move-settings="fsMoveFile.moveSettings"
        :show-dailog="fsMoveFile.showDailog"
        @on-error="onHiddenMove"
        @on-stop="onHiddenMove"
        @on-end="onHiddenMove"
      ></movefile>
      <!-- 选择目录 -->
      <fileSelector
        :select-muti="fsSelector.selectMuti"
        :select-file="fsSelector.selectFile"
        :select-dir="fsSelector.selectDir"
        :start-path="fsSelector.startPath"
        :show-dailog="fsSelector.showDailog"
        @on-select="onSelectedFile"
        @on-cancel="onSelectCancel"
      ></fileSelector>
      <right-click-menu
        bind-ref="fsSelection"
        :menus="fsOperationButtons"
      ></right-click-menu>
    </div>
  </div>
</template>

<script>
import fileAddress from "./file/address.vue";
import uploadfile from "./file/uploadfile.vue";
import copyfile from "./file/copyfile.vue";
import movefile from "./file/movefile.vue";
import fileSelector from "./file/fileselector.vue";
import fileicon from "./file/fileicon.vue";

import { $utils } from "../js/utils";
import { $apitools } from "../js/apis/apitools";
import { $fileopts } from "../js/apis/fileopts";
import { $filepms } from "../js/apis/filepermission";
export default {
  name: "FileList",
  components: {
    fileAddress,
    uploadfile,
    copyfile,
    movefile,
    fileSelector,
    fileicon,
  },
  data: function () {
    let _ = this;
    return {
      searchText: "",
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
        loadPath: "",
      },
      fsUpload: {
        showDrawer: false,
      },
      loading: false,
      fsColumns: [
        {
          type: "selection",
          width: 60,
          align: "center",
        },
        {
          title: "文件名称",
          key: "path",
          render: function (h, params) {
            return h(fileicon, {
              props: {
                node: params.row,
                isEditor: params.row.showRename ? true : false,
              },
              on: {
                click: _.doOpen,
                doRename: function (path, name) {
                  _.onRenameAfter(params.index, params.row, name);
                },
              },
            });
          },
        },
        {
          title: "修改时间",
          key: "mtime",
          maxWidth: 160,
          render: function (h, params) {
            return h("span", $utils.long2Time(params.row.mtime));
          },
        },
        {
          title: "文件大小",
          key: "size",
          maxWidth: 160,
          render: function (h, params) {
            if (params.row.isFile) {
              return h("span", $utils.formatSize(params.row.size));
            } else {
              h("span", "");
            }
          },
        },
        {
          title: "操作权限",
          key: "PermissionText",
          maxWidth: 280,
        },
      ],
      fsData: [],
      fsData_All: [],
      fsData_Total: 0,
      fsData_PageSize: 40,
      fsData_CurrentPageIndex: 1,
      permissionsMap: {},
      fsOperationButtons: {
        refresh: {
          name: "刷新",
          show: true,
          handler: function () {
            _.doRefresh();
          },
          icon: "ivu-icon ivu-icon-md-refresh-circle",
          divided: false,
        },
        upload: {
          name: "上传",
          show: false,
          handler: function () {
            _.fsUpload.showDrawer = true;
          },
          icon: "ivu-icon ivu-icon-md-cloud-upload",
          divided: false,
        },
        newFolder: {
          name: "文件夹",
          show: false,
          handler: function () {
            _.onNewFolder();
          },
          icon: "ivu-icon ivu-icon-md-add-circle",
          divided: false,
        },
        download: {
          name: "下载",
          show: false,
          handler: function () {
            _.doDownload();
          },
          icon: "ivu-icon ivu-icon-md-download",
          divided: false,
        },
        rename: {
          name: "重命名",
          show: false,
          handler: function () {
            _.onRename();
          },
          icon: "ivu-icon ivu-icon-ios-create",
          divided: false,
        },
        move: {
          name: "移动",
          show: false,
          handler: function () {
            _.doMove();
          },
          icon: "ivu-icon ivu-icon-md-move",
          divided: false,
        },
        copy: {
          name: "复制",
          show: false,
          handler: function () {
            _.doCopy();
          },
          icon: "ivu-icon ivu-icon-ios-copy",
          divided: false,
        },
        delete: {
          name: "删除",
          show: false,
          handler: function () {
            _.onDelete();
          },
          icon: "ivu-icon ivu-icon-md-trash",
          divided: false,
        },
      },
    };
  },
  computed: {},
  methods: {
    doRefresh: function () {
      this.doLoadData(this.fsAddress.loadPath);
    },
    doLoadData: function (path) {
      if (this.loading) {
        return;
      }
      let _ = this;
      _.loading = true;
      this.onSelectionChange([], {});
      this.fsData_All = [];
      this.fsData = [];
      $fileopts
        .List(path)
        .then(function (data) {
          _.fsData_All = JSON.parse(data);
          _.fsData_CurrentPageIndex = 1;
          _.loadDatePage();
          // _.loading = false;
        })
        .catch(function (err) {
          console.error(err);
          _.loading = false;
          _.$Message.error(err.toString());
        });
    },
    loadDatePage: function (index) {
      let dataAll = [];
      if (this.searchText) {
        let stemp = this.searchText.toLowerCase();
        for (let i = 0; i < this.fsData_All.length; i++) {
          if (
            this.fsData_All[i].path.getName().toLowerCase().indexOf(stemp) > -1
          ) {
            dataAll.push(this.fsData_All[i]);
          }
        }
      } else {
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
      datas.forEach(function (row) {
        paths.push(row.path);
      });
      let _ = this;
      $filepms
        .GetUserPermissionSum($apitools.getSession().userID, paths)
        .then(function (pms) {
          _.permissionsMap = pms ? JSON.parse(pms) : {};
          for (let i = 0; i < datas.length; i++) {
            datas[i].Permission =
              undefined != _.permissionsMap[datas[i].path]
                ? _.permissionsMap[datas[i].path]
                : -1;
            datas[i].PermissionText = $filepms.$TYPE
              .sum2Name(datas[i].Permission)
              .join(", ");
          }
          _.fsData = datas;
          _.onSelectionChange();
          _.loading = false;
        })
        .catch(function (err) {
          _.loading = false;
          _.$Message.error(err.toString());
        });
    },
    /* 重命名 */
    doRename: function (index, row, name) {
      row.showRename = false;
      this.$set(this.fsData, index, row);
      if (!name || row.path.getName() == name) {
        return;
      }
      let _ = this;
      $fileopts
        .Rename(row.path, name)
        .then(function () {
          _.$Message.info("操作成功");
          _.doRefresh();
        })
        .catch(function (err) {
          _.$Message.error("操作失败: " + err.toString());
          _.doRefresh();
        });
    },
    /** 新建文件夹 */
    doNewFolder: function (index, row, name) {
      if (!name) {
        this.doRefresh();
        return;
      }
      let _ = this;
      $fileopts
        .NewFolder(this.fsAddress.loadPath + "/" + name)
        .then(function () {
          _.$Message.info("操作成功");
          _.doRefresh();
        })
        .catch(function (err) {
          _.$Message.error("操作失败: " + err.toString());
          _.doRefresh();
        });
    },
    onDelete: function () {
      let _ = this;
      this.$Modal.confirm({
        title: "删除文件",
        content: "此操作不可逆, 是否继续删除?",
        onOk: function () {
          _.doDelete();
        },
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
        duration: 0,
      });
      let errs = [];
      let _ = this;
      function loop(i) {
        if (!i) {
          i = 0;
        }
        if (i == select.length) {
          del_loading();
          _.doRefresh();
          if (errs && errs.length > 0) {
            _.$Message.error({
              content: errs.join("</br>"),
              duration: 5,
            });
          }
          return;
        }
        $fileopts
          .Delete(select[i].path)
          .then(function () {
            setTimeout(function () {
              loop(++i);
            });
          })
          .catch(function (err) {
            errs.push(select[i].path + "删除失败: " + err);
            loop(++i);
          });
      }
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
        if (!select[i].isFile) {
          continue;
        }
        $fileopts
          .GetDownloadUrl(select[i].path)
          .then(function (data) {
            _.$nextTick(function () {
              let download = document.createElement("iframe");
              download.src = data;
              document.body.appendChild(download);
            });
          })
          .catch(function (err) {
            _.$Message.error(
              "操作失败: " + select[i].path + ", " + err.toString()
            );
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
        if (this.fsData[i].path == select[0].path) {
          this.$set(this.fsData, i, select[0]);
        }
      }
    },
    /* 打开文件 */
    doOpen: function (node) {
      if (!node.isFile) {
        this.goToPath(node.path);
      } else {
        let _ = this;
        $preview.doPreview(node.path).catch(function (err) {
          $fileopts
            .GetSteamUrl(node.path)
            .then(function (data) {
              _.$nextTick(function () {
                window.open(data);
              });
            })
            .catch(function (err) {
              _.$Message.error(
                "操作失败: " + node.path + ", " + err.toString()
              );
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
        path: "",
        isFile: false,
        IsNewFolder: true,
        showRename: true,
      });
      this.$set(this, "fsData", temp);
    },
    /* 选择单个文件夹 - OK */
    onSelectedFile: function (rows) {
      this.fsSelector.showDailog = false;
      if (rows && rows[0] && rows[0].path && rows[0].path.length > 0) {
        if (this.fsSelector.operationObj) {
          // this.fsCopyFile.destPath = rows[0].path;
          // this.fsCopyFile.showDailog = true;
          this.fsSelector.operationObj.destPath = rows[0].path;
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
      if (
        undefined != this.permissionsMap[this.fsAddress.loadPath] &&
        $filepms.$TYPE.sumInclude(
          this.permissionsMap[this.fsAddress.loadPath],
          $filepms.$TYPE.WRITE.value
        )
      ) {
        this.fsOperationButtons.upload.show = true;
        this.fsOperationButtons.newFolder.show = true;
      }
      // 选择一个或者以上
      if (len_selection >= 1) {
        let read = false;
        let fileCount = 0;
        for (let i = 0; i < selection.length; i++) {
          if (selection[i].isFile) {
            fileCount++;
          }
          read = $filepms.$TYPE.sumInclude(
            selection[i].Permission,
            $filepms.$TYPE.READ.value
          );
          if (!read) {
            read = false;
            break;
          }
        }
        if (read) {
          this.fsOperationButtons.download.show = fileCount == selection.length;
          this.fsOperationButtons.copy.show = true;
        }
        //
        let write = false;
        for (let i = 0; i < selection.length; i++) {
          write = $filepms.$TYPE.sumInclude(
            selection[i].Permission,
            $filepms.$TYPE.WRITE.value
          );
          if (!write) {
            write = false;
            break;
          }
        }
        if (write) {
          this.fsOperationButtons.rename.show = len_selection == 1;
          this.fsOperationButtons.move.show = true;
          this.fsOperationButtons.delete.show = true;
        }
      }
    },
  },
  mounted: function () {},
  created: function () {
    this.goToPath(this.fsAddress.rootPath);
    this.fsOperationButtons.upload.show = true;
    let _ = this;
    this.$nextTick(function () {
      $utils.areaCover({
        background: _.$refs.fsSelection.$el,
        coverfilter: ".ivu-table-row",
        onCover: function (el) {
          if (!el.querySelector(".ivu-checkbox-input").checked) {
            el.querySelector(".ivu-checkbox-input").click();
          }
        },
        onUnCover: function (el) {
          if (el.querySelector(".ivu-checkbox-input").checked) {
            el.querySelector(".ivu-checkbox-input").click();
          }
        },
      });
    });
  },
  watch: {
    "fsAddress.loadPath": function (n, o) {
      this.doLoadData(n);
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
