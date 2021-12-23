<template>
  <Drawer
    title="上传文件"
    width="450px"
    v-model="isShowDrawer"
    @on-close="$emit('on-close')"
  >
    <div class="ivu-upload">
      <div
        class="ivu-upload ivu-upload-drag"
        ref="uploadDrag"
        @click="doSelectFiels"
      >
        <input
          ref="upload_selector_file"
          type="file"
          multiple="multiple"
          class="ivu-upload-input"
        />
        <div style="padding: 20px 0px">
          <i
            class="ivu-icon ivu-icon-ios-cloud-upload"
            style="font-size: 52px; color: rgb(51, 153, 255)"
          ></i>
          <p>点击或拖拽到此处</p>
        </div>
      </div>
      <ul class="ivu-upload-list">
        <li
          v-for="(temp, index) in files"
          v-bind:key="index"
          v-if="!temp._upload.removed"
          class="ivu-upload-list-file"
        >
          <span>{{ temp.name }}</span>
          <i
            class="ivu-icon ivu-icon-ios-close ivu-upload-list-remove"
            @click="removeTask(temp._upload.index)"
          ></i>
          <div class="ivu-progress ivu-progress-normal ivu-progress-show-info">
            <div class="ivu-progress-outer">
              <div class="ivu-progress-inner">
                <div
                  v-if="!temp._upload.ps || temp._upload.ps < 100"
                  class="ivu-progress-bg"
                  :style="{
                    width: (temp._upload.ps ? temp._upload.ps : 0) + '%',
                    height: '2px',
                  }"
                ></div>
                <div
                  v-else
                  class="ivu-progress-success-bg"
                  style="width: 100%; height: 2px"
                ></div>
              </div>
            </div>
            <span class="ivu-progress-text">
              <span
                :style="{ color: temp._upload.err ? '#b42525' : '#515a6e' }"
                class="ivu-progress-text-inner"
                >{{
                  temp._upload.err ? temp._upload.err : temp._upload.ps + "%"
                }}</span
              >
            </span>
          </div>
        </li>
      </ul>
    </div>
  </Drawer>
</template>

<script>
import { $utils } from "../../js/utils";
import { $fileopts } from "../../js/apis/fileopts";
export default {
  name: "uploadfile",
  props: ["show-drawer", "parent", "drag-ref"],
  data() {
    return {
      isShowDrawer: false, // 显示抽屉
      maxuploading: 5, // 最大正在上传的个数
      countuploading: 0, // 正在上传的个数
      dindex: 0, // 当前数据下标
      files: [], // 文件
      queueend: true, // 队列可用为空
    };
  },
  methods: {
    // 上传-触发选择文件
    doSelectFiels(ev) {
      let emited = false;
      $utils.addEvent(
        this.$refs.upload_selector_file,
        "change",
        (ev_data) => {
          if (ev_data.target.files) {
            for (let i = 0; i < ev_data.target.files.length; i++) {
              let fs = ev_data.target.files[i];
              fs._upload = {
                base: this.parent,
                index: this.files.length,
                updater: false,
                ps: 0,
              };
              if (!emited) {
                emited = true;
                this.$emit("on-start");
              }
              this.files.push(fs);
              this.doStartUpload();
            }
          }
          this.$refs.upload_selector_file.value = "";
        },
        { once: true }
      );
      $utils.triggerMouseEvent(this.$refs.upload_selector_file, "click");
    },
    // 上传-触发上传动作
    doStartUpload() {
      if (this.files && this.files.length > 0) {
        if (this.countuploading >= this.maxuploading) {
          return;
        }
        if (this.dindex >= this.files.length) {
          this.queueend = true;
          return;
        } else if (this.queueend) {
          this.queueend = false;
        }
        this.countuploading++;
        let file = this.files[this.dindex++];
        if (file._upload.removed) {
          this.countuploading--;
          this.$nextTick(this.doStartUpload);
          return;
        }
        // file._upload.index = this.dindex-1;
        let opts = {
          form: {},
          header: {},
          progress: (e) => {
            // 数据源为数组, 需要直接设置数组
            file._upload.ps = Math.round((e.loaded / e.total) * 1000) / 10;
            this.$set(this.files, file._upload.index, file);
          },
          error: (e) => {
            file._upload.err = e ? e.toString() : "上传失败";
            this.$set(this.files, file._upload.index, file);
          },
          abort: (e) => {
            file._upload.err = "上传取消";
            this.$set(this.files, file._upload.index, file);
          },
          loadstart: (e) => {},
          loadend: (e) => {
            this.countuploading--;
            file._upload.ended = true;
            this.$nextTick(() => {
              this.doStartUpload();
              //this.removeTask(file._upload.index);
            });
          },
        };
        // 预备开始
        file._upload.started = true;
        $fileopts
          .GetUploadUrl(file._upload.base + "/" + file.name)
          .then((url) => {
            file._upload.updater = $utils.uploadByFormData(url, file, opts);
            file._upload.updater.start();
          })
          .catch((err) => {
            console.error(err);
            opts.error(err);
            opts.loadend();
          });
      }
    },
    // 上传 - 移除任务
    removeTask(index) {
      let file = this.files[index];
      if (file) {
        // 正在传输
        if (!file._upload.ended && file._upload.started) {
          file._upload.updater.abort();
        }
        file._upload.removed = true;
        this.$set(this.files, index, file);
      }
    },
    // 上传 - 监听拖拽
    doListenDrag(key) {
      this.$nextTick(() => {
        let obj = undefined;
        if (key) {
          if (this.$refs[key]) {
            obj = this.$refs[key];
          } else if (this.$parent && this.$parent.$refs) {
            obj = this.$parent.$refs[key];
          }
        }
        if (!obj) {
          return;
        }
        obj.ondrop = (evn) => {
          evn.preventDefault();
          let emited = false;
          let fileList = evn.dataTransfer.files;
          for (let i = 0; i < fileList.length; i++) {
            let fs = fileList[i];
            if (!fs.type && fs.size == 0) {
              continue;
            }
            fs._upload = {
              base: this.parent,
              index: this.files.length,
              updater: false,
              ps: 0,
            };
            if (!emited) {
              emited = true;
              this.$emit("on-start");
            }
            this.files.push(fs);
            this.doStartUpload();
          }
        };
      });
    },
    // 停止所有拖拽动作
    stopDrop() {
      try {
        this.$el.ondrop = function (Even) {
          Even.preventDefault();
          Even.stopPropagation();
        };
        this.$el.ondragover = function (Even) {
          Even.preventDefault();
          Even.stopPropagation();
        };
        this.$el.ondragleave = function (Even) {
          Even.preventDefault();
          Even.stopPropagation();
        };
        this.$el.ondragenter = function (Even) {
          Even.preventDefault();
          Even.stopPropagation();
        };
        this.$el.oncontextmenu = function (Even) {
          Even.preventDefault();
          Even.stopPropagation();
        };
      } catch (Err_Catch) {
        console.error(Err_Catch);
      }
    },
  },
  created() {
    this.$nextTick(() => {
      this.stopDrop();
      this.doListenDrag("uploadDrag");
      this.doListenDrag(this.dragRef);
    });
  },
  watch: {
    countuploading(n) {
      if (n <= 0 && this.dindex >= this.files.length) {
        this.$emit("on-end");
      }
    },
    showDrawer(n) {
      this.isShowDrawer = n;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
