<template>
  <div style="width: 100%" @click="onBoxClick">
    <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap">
      <img
        :src="icon"
        style="
          vertical-align: middle;
          margin-right: 5px;
          width: 32px;
          height: 32px;
        "
      />
      <i-input
        v-if="isEditor"
        v-model="filename"
        style="width: calc(100% - 50px)"
        @on-blur="doSave"
        @on-enter="doSave"
      ></i-input>
      <span
        v-else
        style="cursor: pointer"
        @click="onNameClick"
        :title="filename"
        >{{ filename }}</span
      >
    </div>
  </div>
</template>

 
<script>
import { iconUrl } from "../../js/bizutil";
export default {
  name: "fileicon",
  props: ["node", "isEditor"],
  data: function () {
    return {
      iconTypes: [
        "ai",
        "avi",
        "bmp",
        "catdrawing",
        "catpart",
        "catproduct",
        "cdr",
        "csv",
        "doc",
        "docx",
        "dps",
        "dpt",
        "dwg",
        "eio",
        "eml",
        "et",
        "ett",
        "exb",
        "exe",
        "file",
        "flv",
        "fold",
        "gif",
        "htm",
        "html",
        "jpeg",
        "jpg",
        "mht",
        "mhtml",
        "mid",
        "mp3",
        "mp4",
        "mpeg",
        "msg",
        "odp",
        "ods",
        "odt",
        "pdf",
        "png",
        "pps",
        "ppt",
        "pptx",
        "prt",
        "psd",
        "rar",
        "rm",
        "rmvb",
        "rtf",
        "sldprt",
        "swf",
        "tif",
        "txt",
        "url",
        "wav",
        "wma",
        "wmv",
        "wps",
        "wpt",
        "xls",
        "xlsx",
        "zip",
      ],
      icon: "",
      filename: "",
    };
  },
  methods: {
    onBoxClick: function (e) {
      if (e.srcElement.tagName.toUpperCase() == "INPUT") {
        e.stopPropagation();
      }
    },
    onNameClick: function (e) {
      e.stopPropagation();
      this.$emit("click", this.node, e);
    },
    doSave: function () {
      if (this.isEditor === true) {
        this.isEditor = false;
        this.$emit("doRename", this.node.path, this.filename);
      }
    },
    getFileIcon: function (path) {
      if (!this.node.isFile) {
        return "/static/img/file_icons/folder.png";
      }
      return iconUrl(path);
    },
    initvalue: function () {
      this.filename = this.node.path.getName();
      this.icon = this.getFileIcon(this.filename);
    },
  },
  created: function () {
    let that = this;
    this.initvalue();
    this.$nextTick(function () {
      window.addEventListener("keydown", function (e) {
        if (!e) {
          e = window.event;
        }
        if ((e.keyCode || e.which) == 13) {
          that.doSave();
        }
      });

      window.addEventListener("mousedown", function (e) {
        let dom = e.target;
        if (dom == that.$el.querySelector("input")) {
        } else {
          that.doSave();
        }
      });
    });
  },
  watch: {
    node: {
      deep: true,
      handler: function (newval, oldval) {
        this.initvalue();
      },
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
