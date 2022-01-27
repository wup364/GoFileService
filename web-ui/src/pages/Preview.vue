<template>
  <div class="page" ref="preview-page">
    <router-view @on-loading="onLoading" />
    <Spin v-show="loadding" fix>
      <Icon
        type="ios-loading"
        size="20"
        style="animation: ani-demo-spin 1s linear infinite"
      ></Icon>
      <div>{{ loaddingText }}</div>
    </Spin>
  </div>
</template>

<script>
import { $filepreview } from "../js/apis/filepreview";
export default {
  name: "Preivew",
  data() {
    return {
      loadding: true,
      loaddingText: "加载中...",
      title: "预览文件",
      playerType: "",
      playerSrc: "",
    };
  },
  computed: {},
  methods: {
    setTitle() {
      document.title = this.title ? this.title : "查看文件";
    },
    onLoading(opts) {
      if (opts) {
        Object.keys(opts).forEach((key) => {
          this.$set(this, key, opts[key]);
        });
      }
    },
  },
  created() {
    let token = this.$route.query.token;
    this.playerType = this.$route.path.getName();
    if (token && token.length > 0 && this.playerType) {
      switch (this.playerType) {
        case "audio":
          this.title = "音频播放";
          break;
        case "video":
          this.title = "视频播放";
          break;
        case "picture":
          this.title = "图片播放";
          break;
        default:
          this.title = "暂不支持的格式";
          this.loaddingText = this.title;
          break;
      }
      setInterval(() => $filepreview.status(token), 60000);
    } else {
      this.title = "链接似乎是无效的呢~";
      this.$Message.error("链接错误, 请重新打开");
      setTimeout(window.close, 2000);
    }
  },
  watch: {
    title(n, o) {
      this.setTitle();
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
@keyframes ani-demo-spin {
  from {
    transform: rotate(0deg);
  }
  50% {
    transform: rotate(180deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
