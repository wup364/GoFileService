<template>
  <Modal
    v-model="showDailog"
    :closable="false"
    :mask-closable="false"
    :width="copySettings.width"
  >
    <p slot="header" style="color: #f60; text-align: center">
      <span>复制文件</span>
    </p>
    <div style="text-align: left">
      <div>
        <table>
          <tr>
            <td>正在处理:</td>
            <td>{{ operationData.multiCount + operationData.opCount }}</td>
          </tr>
          <tr>
            <td width="30%">复制路径:</td>
            <td>{{ operationData.nowSrcPath }}</td>
          </tr>
          <tr>
            <td>目标位置:</td>
            <td>{{ operationData.nowDstPath }}</td>
          </tr>
          <tr v-show="copyError.IsError">
            <td>出现错误:</td>
            <td>{{ copyError.Error }}</td>
          </tr>
        </table>
      </div>
    </div>
    <div slot="footer">
      <div v-show="copyError.IsError">
        <Button long :style="{ margin: '5px 2px' }" @click="doIgnore"
          >跳过</Button
        >
        <Button
          long
          :style="{ margin: '5px 2px' }"
          v-show="copyError.IsExist"
          @click="doReplace"
          >覆盖</Button
        >
        <div style="margin: 8px 0px; color: #57a3f3">
          <Checkbox v-model="copySettings.ignore">自动跳过出错文件</Checkbox>
          <Checkbox v-model="copySettings.replace" v-show="copyError.IsExist"
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
  name: "copyfile",
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
    copySettings: {
      type: Object,
      default: () => {
        return { ignore: false, replace: false, width: 450 };
      },
    },
  },
  data: function () {
    return {
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
      copyError: {
        IsError: false,
        Error: "",
        IsExist: false,
      },
    };
  },
  methods: {
    doCopy: function () {
      let _ = this;
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
          .CopyAsync(
            tempSrc.path,
            tempDst,
            this.copySettings.replace,
            this.copySettings.ignore
          )
          .then(function (data) {
            if (_.operationData.srcPaths.length > 1) {
              _.operationData.srcPaths = _.operationData.srcPaths.slice(
                0,
                _.operationData.srcPaths.length - 1
              );
            } else {
              _.operationData.srcPaths = [];
            }
            _.operationData.token = data;
          })
          .catch(function (err) {
            _.$Message.error(err.toString());
            _.$emit("on-error");
          });
      }
    },
    doRefreshPs: function () {
      if (!this.operationData.token || this.operationData.token == "") {
        return;
      }
      let _ = this;
      $fileopts
        .AsyncExecToken("CopyFile", this.operationData.token)
        .then(function (data) {
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
          data = JSON.parse(data);
          // console.log( data )
          if (data.CountIndex > 0) {
            _.operationData.opCount = data.CountIndex;
          }
          if (data.IsComplete) {
            _.operationData.token = "";
            if (data.IsDiscontinue) {
              _.showDailog = false;
              _.$Message.error("复制已终止");
              _.$emit("on-stop");
            } else {
              if (
                !_.operationData.srcPaths ||
                _.operationData.srcPaths.length == 0
              ) {
                _.showDailog = false;
                if (data.ErrorString && data.ErrorString.length > 0) {
                  _.$Message.error(data.ErrorString);
                } else {
                  _.$Message.success("复制完成");
                }
                _.$emit("on-end");
              } else {
                _.operationData.multiCount += _.operationData.opCount;
                _.doCopy();
              }
            }
            return;
          }
          _.operationData.nowSrcPath = data.Src;
          _.operationData.nowDstPath = data.Dst;
          //
          _.copyError.IsError =
            data.ErrorString && data.ErrorString.length > 0 ? true : false;
          _.copyError.Error = _.parseError(data);
          _.copyError.IsExist = data.IsDstExist;
          setTimeout(function () {
            _.doRefreshPs();
          }, 100);
        })
        .catch(function (err) {
          _.$Message.error(err.toString());
        });
    },
    doStop: function () {
      if (this.operationData.token && this.operationData.token.length > 0) {
        let _ = this;
        $fileopts
          .AsyncExecToken("CopyFile", this.operationData.token, {
            operation: this.operations.stop,
          })
          .then(function (data) {
            _.operationData.opCount = 0;
            _.operationData.multiCount = 0;
            _.operationData.srcPaths = [];
            _.operationData.destPath = [];
          })
          .catch(function (err) {
            _.$Message.error(err.toString());
          });
      }
    },
    doIgnore: function () {
      if (!this.operationData.token || this.operationData.token == "") {
        return;
      }
      let _ = this;
      $fileopts
        .AsyncExecToken("CopyFile", this.operationData.token, {
          operation: this.copySettings.ignore
            ? this.operations.ignoreall
            : this.operations.ignore,
        })
        .then(function (data) {
          // console.log(data);
        })
        .catch(function (err) {
          _.$Message.error(err.toString());
        });
    },
    doReplace: function () {
      if (!this.operationData.token || this.operationData.token == "") {
        return;
      }
      let _ = this;
      $fileopts
        .AsyncExecToken("CopyFile", this.operationData.token, {
          operation: this.copySettings.replace
            ? this.operations.replaceall
            : this.operations.replace,
        })
        .then(function (data) {
          // console.log(data);
        })
        .catch(function (err) {
          _.$Message.error(err.toString());
        });
    },
    parseError: function (data) {
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
    "operationData.token": function (n, o) {
      if (n && n != "") {
        this.doRefreshPs();
      }
    },
    showDailog: function (n, o) {
      if (n) {
        this.copyError.IsError = false;
        this.copyError.Error = "";
        this.copyError.IsExist = false;
        this.copySettings.ignore = false;
        this.copySettings.replace = false;
        this.operationData.opCount = 0;
        this.operationData.multiCount = 0;
        this.operationData.showDailog = this.showDailog;
        this.operationData.srcPaths = this.srcPaths;
        this.operationData.destPath = this.destPath;
        this.doCopy();
      }
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
