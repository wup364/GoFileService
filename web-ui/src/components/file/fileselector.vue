<template>
  <Modal
    v-model="isShowDailog"
    :title="selectFile ? '选择文件' : '选择目录'"
    :width="fsSettings.width"
    @on-ok="onOk"
    @on-cancel="onCancel"
  >
    <div style="height: 30px; line-height: 30px; padding-left: 5px">
      <fileAdress
        v-if="fsStatus.loadedPath && fsStatus.loadedPath.length > 0"
        :depth="4"
        :rootname="fsSettings.rootname"
        :path="fsStatus.loadedPath"
        @click="goToPath"
      ></fileAdress>
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
import fileicon from "./fileicon.vue";
import fileAdress from "./address.vue";
import { $utils } from "../../js/utils";
import { $fileopts } from "../../js/apis/fileopts";
export default {
  name: "fileselector",
  components: { fileAdress },
  props: [
    "show-dailog",
    "settings",
    "start-path",
    "select-muti",
    "select-file",
    "select-dir",
  ],
  data() {
    let _ = this;
    return {
      isShowDailog: false,
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
          render(h, params) {
            return h("checkbox", {
              props: {
                value:
                  params.row._checked == undefined
                    ? false
                    : params.row._checked,
              },
              on: {
                "on-change"(val) {
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
          render(h, params) {
            return h(fileicon, {
              props: {
                node: params.row,
                editMode: false,
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
          render(h, params) {
            return h("span", $utils.long2Time(params.row.CtTime));
          },
        },
      ],
      fsData: [],
      selectedDates: [],
    };
  },
  created() {},
  methods: {
    init() {
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
    goToPath(path) {
      if (this.fsStatus.loadedPath != path) {
        this.fsStatus.fsLoading = true;
        this.fsStatus.loadedPath = path;
      }
    },
    doOpenDir(node) {
      if (!node.isFile) {
        this.goToPath(node.path);
      }
    },
    onOk() {
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
    onCancel() {
      this.$emit("on-cancel");
    },
    onRowClick(row, index) {
      for (let i = 0; i < this.fsData.length; i++) {
        if (this.fsData[i]._checked && index != i) {
          this.$set(this.fsData[i], "_checked", false);
        }
      }
      this.selectedDates = [row];
      this.$set(this.fsData[index], "_checked", true);
    },
    putSelect(row) {
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
    removeSelect(row) {
      for (let i = this.selectedDates.length - 1; i >= 0; i--) {
        if (this.selectedDates[i].path == row.path) {
          this.selectedDates.remove(i);
        }
      }
    },
    onSelectionChange(selection, row) {},
  },
  watch: {
    showDailog(n, o) {
      if (n) {
        this.init();
      } else {
        this.fsStatus.loadedPath = "";
      }
      this.isShowDailog = n;
    },
    "fsStatus.loadedPath"(n, o) {
      if (!n) {
        return;
      }
      $fileopts
        .List(n)
        .then((data) => {
          this.fsData = [];
          this.selectedDates = [];
          for (let i = 0; i < data.length; i++) {
            if (
              (this.isSelectFile && data[i].isFile) ||
              (this.isSelectDir && !data[i].isFile)
            ) {
              data[i]["_checked"] = false;
              this.fsData.push(data[i]);
            }
          }
          this.fsStatus.fsLoading = false;
        })
        .catch((err) => {
          this.fsStatus.fsLoading = false;
          this.$Message.error(err.toString());
        });
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
