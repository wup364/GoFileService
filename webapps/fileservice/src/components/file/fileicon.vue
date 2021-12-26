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
        v-if="isEditMode"
        v-model="filename"
        style="width: calc(100% - 50px)"
        @on-blur="doSave"
        @on-enter="doSave"
      ></i-input>
      <a v-else style="color: #515a6e" @click="onNameClick" :title="filename">{{
        filename
      }}</a>
    </div>
  </div>
</template>

 
<script>
import { iconUrl } from "../../js/bizutil";
export default {
  name: "fileicon",
  props: ["node", "edit-mode"],
  data() {
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
      isEditMode: false,
    };
  },
  methods: {
    onBoxClick(e) {
      if (e.srcElement.tagName.toUpperCase() == "INPUT") {
        e.stopPropagation();
      }
    },
    onNameClick(e) {
      e.stopPropagation();
      this.$emit("click", this.node, e);
    },
    doSave() {
      if (this.isEditMode === true) {
        this.isEditMode = false;
        this.$emit("rename", this.node.path, this.filename);
      }
    },
    getFileIcon(path) {
      if (!this.node.isFile) {
        return "/static/img/file_icons/folder.png";
      }
      return iconUrl(path);
    },
    initvalue() {
      this.filename = this.node.path.getName();
      this.icon = this.getFileIcon(this.filename);
    },
  },
  created() {
    this.initvalue();
    this.$nextTick(() => {
      window.addEventListener("keydown", (e) => {
        e = e || window.event;
        if ((e.keyCode || e.which) == 13) {
          this.doSave();
        }
      });

      window.addEventListener("mousedown", (e) => {
        let dom = e.target;
        if (dom == this.$el.querySelector("input")) {
        } else {
          this.doSave();
        }
      });
    });
  },
  watch: {
    node: {
      deep: true,
      handler(newval, oldval) {
        this.initvalue();
      },
    },
    editMode(n) {
      this.isEditMode = n;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
