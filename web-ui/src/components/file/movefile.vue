<template>
  <Modal
    v-model="isShowDailog"
    :closable="false"
    :mask-closable="false"
    :width="moveSettings.width"
  >
    <p slot="header" style="color: #f60; text-align: center">
      <span>移动文件</span>
    </p>
    <div style="text-align: left">
      <div>
        <table>
          <tr>
            <td width="30%">移动路径:</td>
            <td>{{ operationData.nowSrcPath }}</td>
          </tr>
          <tr>
            <td>目标位置:</td>
            <td>{{ operationData.nowDstPath }}</td>
          </tr>
          <tr v-show="moveError.IsError">
            <td>出现错误:</td>
            <td>{{ moveError.Error }}</td>
          </tr>
        </table>
      </div>
    </div>
    <div slot="footer">
      <div v-show="moveError.IsError">
        <Button long :style="{ margin: '5px 2px' }" @click="doIgnore"
          >跳过</Button
        >
        <Button
          long
          :style="{ margin: '5px 2px' }"
          v-show="moveError.IsExist"
          @click="doReplace"
          >覆盖</Button
        >
        <div style="margin: 8px 0px; color: #57a3f3">
          <Checkbox v-model="moveSettings.ignore">自动跳过出错文件</Checkbox>
          <Checkbox v-model="moveSettings.replace" v-show="moveError.IsExist"
            >自动覆盖重复文件</Checkbox
          >
        </div>
      </div>
      <Button type="error" long :style="{ margin: '5px 2px' }" @click="doStop"
        >终止</Button
      >
    </div>
  </Modal>
</template>

 
<script>
import { $fileopts } from "../../js/apis/fileopts";
export default {
  name: "movefile",
  props: {
    showDailog: {
      type: Boolean,
      default: false,
    },
    srcPaths: {
      type: Array,
      default: () => {
        return [];
      },
    },
    destPath: { type: String, default: "/" },
    moveSettings: {
      type: Object,
      default: () => {
        return { ignore: false, replace: false, width: 450 };
      },
    },
  },
  data() {
    return {
      isShowDailog: false,
      operations: {
        stop: "discontinue",
        ignore: "ignore",
        ignoreall: "ignoreall",
        replace: "replace",
        replaceall: "replaceall",
      },
      operationData: {
        multiCount: 0,
        opCount: 0,
        token: "",
        srcPaths: [],
        destPath: "",
        nowSrcPath: "",
        nowDstPath: "",
      },
      moveError: {
        IsError: false,
        Error: "",
        IsExist: false,
      },
    };
  },
  created() {},
  methods: {
    doMove() {
      if (
        this.operationData.srcPaths &&
        this.operationData.srcPaths.length > 0
      ) {
        let tempSrc =
          this.operationData.srcPaths[this.operationData.srcPaths.length - 1];
        let tempDst =
          this.operationData.destPath +
          "/" +
          tempSrc.path.getName().parsePath();
        $fileopts
          .MoveAsync(
            tempSrc.path,
            tempDst,
            this.moveSettings.replace,
            this.moveSettings.ignore
          )
          .then((data) => {
            if (this.operationData.srcPaths.length > 1) {
              this.operationData.srcPaths = this.operationData.srcPaths.slice(
                0,
                this.operationData.srcPaths.length - 1
              );
            } else {
              this.operationData.srcPaths = [];
            }
            this.operationData.token = data;
          })
          .catch((err) => {
            this.$Message.error(err.toString());
            this.$emit("on-error");
          });
      }
    },
    doRefreshPs() {
      if (!this.operationData.token || this.operationData.token == "") {
        return;
      }
      $fileopts
        .AsyncExecToken("MoveFile", this.operationData.token)
        .then((data) => {
          /*
						{
							"CountIndex":7,
							"ErrorString":"",
							"Src":"/files/Mount01/glibc-ports-2.15.tar.gz",
							"Dst":"/files/.cache/Mount01/glibc-ports-2.15.tar.gz",
							"IsSrcExist":false,
							"IsDstExist":false,
							"IsReplace":false,
							"IsReplaceAll":false,
							"IsIgnore":false,
							"IsIgnoreAll":false,
							"IsComplete":false,
							"IsDiscontinue":false
						}
					*/
          // console.log( data )
          if (data.CountIndex > 0) {
            this.operationData.opCount = data.CountIndex;
          }
          if (data.IsComplete) {
            this.operationData.token = "";
            if (data.IsDiscontinue) {
              this.isShowDailog = false;
              this.$Message.error("移动已终止");
              this.$emit("on-stop");
            } else {
              if (
                !this.operationData.srcPaths ||
                this.operationData.srcPaths.length == 0
              ) {
                this.isShowDailog = false;
                if (data.ErrorString && data.ErrorString.length > 0) {
                  this.$Message.error(data.ErrorString);
                } else {
                  this.$Message.success("移动完成");
                }
                this.$emit("on-end");
              } else {
                this.operationData.multiCount += this.operationData.opCount;
                this.doMove();
              }
            }
            return;
          }
          this.operationData.nowSrcPath = data.Src;
          this.operationData.nowDstPath = data.Dst;
          //
          this.moveError.IsError =
            data.ErrorString && data.ErrorString.length > 0 ? true : false;
          this.moveError.Error = this.parseError(data);
          this.moveError.IsExist = data.IsDstExist;
          setTimeout(() => {
            this.doRefreshPs();
          }, 100);
        })
        .catch((err) => {
          this.$Message.error(err.toString());
        });
    },
    doStop() {
      if (this.operationData.token && this.operationData.token.length > 0) {
        $fileopts
          .AsyncExecToken("MoveFile", this.operationData.token, {
            operation: this.operations.stop,
          })
          .then((data) => {
            this.operationData.opCount = 0;
            this.operationData.multiCount = 0;
            this.operationData.srcPaths = [];
            this.operationData.destPath = [];
          })
          .catch((err) => {
            this.$Message.error(err.toString());
          });
      }
    },
    doIgnore() {
      if (!this.operationData.token || this.operationData.token == "") {
        return;
      }
      $fileopts
        .AsyncExecToken("MoveFile", this.operationData.token, {
          operation: this.moveSettings.ignore
            ? this.operations.ignoreall
            : this.operations.ignore,
        })
        .then((data) => {
          // console.log(data);
        })
        .catch((err) => {
          this.$Message.error(err.toString());
        });
    },
    doReplace() {
      if (!this.operationData.token || this.operationData.token == "") {
        return;
      }
      $fileopts
        .AsyncExecToken("MoveFile", this.operationData.token, {
          operation: this.moveSettings.replace
            ? this.operations.replaceall
            : this.operations.replace,
        })
        .then((data) => {
          // console.log(data);
        })
        .catch((err) => {
          this.$Message.error(err.toString());
        });
    },
    parseError(data) {
      if (data && data.ErrorString) {
        if (data.IsDstExist) {
          return "目标位置已存在: " + data.Dst;
        } else if (!data.IsSrcExist) {
          return "源目录不存在: " + data.Src;
        } else {
          return data.ErrorString;
        }
      }
      return "";
    },
  },
  watch: {
    "operationData.token"(n, o) {
      if (n && n != "") {
        this.doRefreshPs();
      }
    },
    showDailog(n, o) {
      if (n) {
        this.moveError.IsError = false;
        this.moveError.Error = "";
        this.moveError.IsExist = false;
        this.moveSettings.ignore = false;
        this.moveSettings.replace = false;
        this.operationData.opCount = 0;
        this.operationData.multiCount = 0;
        this.operationData.showDailog = n;
        this.operationData.srcPaths = this.srcPaths;
        this.operationData.destPath = this.destPath;
        this.doMove();
      }
      this.isShowDailog = n;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
