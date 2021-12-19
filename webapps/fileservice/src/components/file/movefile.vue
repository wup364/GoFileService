<template>
  <Modal
    v-model="showDailog"
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
            <td>正在处理:</td>
            <td>{{ operationData.multiCount + operationData.opCount }}</td>
          </tr>
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
export default {
  name: "movefile",
  props: ["show-dailog", "src-paths", "dest-path", "move-settings"],
  data: function () {
    return {
      operations: {
        stop: "discontinue",
        ignore: "ignore",
        ignoreall: "ignoreall",
        replace: "replace",
        replaceall: "replaceall",
      },
      moveSettings: {
        ignore: false,
        replace: false,
        width: 450,
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
  created: function () {},
  methods: {
    doMove: function () {
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
          tempSrc.Path.getName().parsePath();
        $fsApi
          .MoveAsync(
            tempSrc.Path,
            tempDst,
            this.moveSettings.replace,
            this.moveSettings.ignore
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
      $fsApi
        .AsyncExecToken("MoveFile", this.operationData.token)
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
              _.$Message.error("移动已终止");
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
                  _.$Message.success("移动完成");
                }
                _.$emit("on-end");
              } else {
                _.operationData.multiCount += _.operationData.opCount;
                _.doMove();
              }
            }
            return;
          }
          _.operationData.nowSrcPath = data.Src;
          _.operationData.nowDstPath = data.Dst;
          //
          _.moveError.IsError =
            data.ErrorString && data.ErrorString.length > 0 ? true : false;
          _.moveError.Error = _.parseError(data);
          _.moveError.IsExist = data.IsDstExist;
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
        $fsApi
          .AsyncExecToken("MoveFile", this.operationData.token, {
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
      $fsApi
        .AsyncExecToken("MoveFile", this.operationData.token, {
          operation: this.moveSettings.ignore
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
      $fsApi
        .AsyncExecToken("MoveFile", this.operationData.token, {
          operation: this.moveSettings.replace
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
        this.moveError.IsError = false;
        this.moveError.Error = "";
        this.moveError.IsExist = false;
        this.moveSettings.ignore = false;
        this.moveSettings.replace = false;
        this.operationData.opCount = 0;
        this.operationData.multiCount = 0;
        this.operationData.showDailog = this.showDailog;
        this.operationData.srcPaths = this.srcPaths;
        this.operationData.destPath = this.destPath;
        this.doMove();
      }
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
