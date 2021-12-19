<template>
  <Modal
    v-model="showDailog"
    :title="selectFile ? '选择文件' : '选择目录'"
    :width="fsSettings.width"
    @on-ok="onOk"
    @on-cancel="onCancel"
  >
    <div style="height: 30px; line-height: 30px; padding-left: 5px">
      <fs-address
        v-if="fsStatus.loadedPath && fsStatus.loadedPath.length > 0"
        :depth="4"
        :rootname="fsSettings.rootname"
        :path="fsStatus.loadedPath"
        @click="goToPath"
      ></fs-address>
    </div>
    <Table
      :loading="fsStatus.fsLoading"
      :columns="fsColumns"
      :data="fsData"
      :height="fsSettings.height"
      @on-row-click="onRowClick"
      @on-selection-change="onSelectionChange"
    ></Table>
  </Modal>
</template>

 
<script>
import { $utils } from "../../js/utils";
import { $fileopts } from "../../js/apis/fileopts";
export default {
  name: "fileselector",
  props: [
    "show-dailog",
    "settings",
    "start-path",
    "select-muti",
    "select-file",
    "select-dir",
  ],
  data: function () {
    let _ = this;
    return {
      isSelectFile: false,
      isSelectDir: false,
      isSelectMuti: false,
      fsStatus: {
        loadedPath: "",
        fsLoading: true,
      },
      fsSettings: {
        rootname: "/",
        width: 680,
        height: 450,
      },
      fsColumns: [
        {
          title: "#",
          width: 60,
          render: function (h, params) {
            return h("checkbox", {
              props: {
                value:
                  params.row._checked == undefined
                    ? false
                    : params.row._checked,
              },
              on: {
                "on-change": function (val) {
                  _.fsStatus.fsLoading = true;
                  if (!_.isSelectMuti) {
                    for (let i = 0; i < _.fsData.length; i++) {
                      if (i == params.index) {
                        continue;
                      }
                      _.fsData[i]["_checked"] = false;
                    }
                  }
                  _.$set(_.fsData[params.index], "_checked", val);
                  if (val) {
                    _.putSelect(params.row);
                  } else {
                    _.removeSelect(params.row);
                  }
                  _.fsStatus.fsLoading = false;
                },
              },
            });
          },
        },
        {
          title: "文件名称",
          key: "path",
          render: function (h, params) {
            return h("fs-fileicon", {
              props: {
                node: params.row,
                isEditor: false,
              },
              on: {
                click: _.doOpenDir,
              },
            });
          },
        },
        {
          title: "修改时间",
          key: "CtTime",
          width: 180,
          render: function (h, params) {
            return h("span", $utils.long2Time(params.row.CtTime));
          },
        },
      ],
      fsData: [],
      selectedDates: [],
    };
  },
  created: function () {},
  methods: {
    init: function () {
      if (this.setting) {
        this.fsSettings.width = this.settings.width
          ? this.settings.width
          : this.fsSettings.width;
        this.fsSettings.height = this.settings.height
          ? this.settings.height
          : this.fsSettings.height;
      }
      //
      this.isSelectFile = this.selectFile ? this.selectFile : false;
      this.isSelectDir = this.selectDir ? this.selectDir : false;
      this.isSelectMuti = this.selectMuti ? eval(this.selectMuti) : false;
      this.selectedDates = [];
      //
      if (this.startPath && this.startPath.length > 0) {
        this.goToPath(this.startPath);
      }
    },
    goToPath: function (path) {
      if (this.fsStatus.loadedPath != path) {
        this.fsStatus.fsLoading = true;
        this.fsStatus.loadedPath = path;
      }
    },
    doOpenDir: function (node) {
      if (!node.isFile) {
        this.goToPath(node.path);
      }
    },
    onOk: function () {
      this.fsStatus.fsLoading = true;
      if (this.fsData) {
        for (let i = 0; i < this.fsData.length; i++) {
          this.fsData[i]["_checked"] = false;
        }
      }
      this.fsStatus.fsLoading = false;
      this.$emit(
        "on-select",
        this.selectedDates && this.selectedDates.length > 0
          ? this.selectedDates
          : this.isSelectDir
          ? [
              {
                path: this.fsStatus.loadedPath,
                isFile: false,
              },
            ]
          : []
      );
    },
    onCancel: function () {
      this.$emit("on-cancel");
    },
    onRowClick: function (row, index) {
      for (let i = 0; i < this.fsData.length; i++) {
        if (this.fsData[i]._checked && index != i) {
          this.$set(this.fsData[i], "_checked", false);
        }
      }
      this.selectedDates = [row];
      this.$set(this.fsData[index], "_checked", true);
    },
    putSelect: function (row) {
      if (!this.isSelectMuti) {
        this.selectedDates = [row];
      } else {
        for (let i = 0; i < this.selectedDates.length; i++) {
          if (this.selectedDates[i].path == row.path) {
            return;
          }
        }
        this.selectedDates.push(row);
      }
    },
    removeSelect: function (row) {
      for (let i = this.selectedDates.length - 1; i >= 0; i--) {
        if (this.selectedDates[i].path == row.path) {
          this.selectedDates.remove(i);
        }
      }
    },
    onSelectionChange: function (selection, row) {},
  },
  watch: {
    showDailog: function (n, o) {
      if (n) {
        this.init();
      } else {
        this.fsStatus.loadedPath = "";
      }
    },
    "fsStatus.loadedPath": function (n, o) {
      if (!n) {
        return;
      }
      let _ = this;
      $fileopts
        .List(n)
        .then(function (data) {
          data = JSON.parse(data);
          _.fsData = [];
          _.selectedDates = [];
          for (let i = 0; i < data.length; i++) {
            if (
              (_.isSelectFile && data[i].isFile) ||
              (_.isSelectDir && !data[i].isFile)
            ) {
              data[i]["_checked"] = false;
              _.fsData.push(data[i]);
            }
          }
          _.fsStatus.fsLoading = false;
        })
        .catch(function (err) {
          _.fsStatus.fsLoading = false;
          _.$Message.error(err.toString());
        });
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
