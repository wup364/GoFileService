<template>
  <div class="page" ref="pw-video-page">
    <div ref="dplayer" id="dplayer" style="height: 100%; width: 100%"></div>
  </div>
</template>

<script>
import DPlayer from "dplayer";
import { $filepreview } from "../../js/apis/filepreview";
export default {
  name: "",
  data() {
    return {};
  },
  methods: {
    initDatas(token) {
      try {
        $filepreview
          .samedirFiles(token)
          .then((datas) => {
            // 获取字幕文件
            let fileName = datas.path.getName(false);
            let subtitleFilePath = "";
            if (datas && datas.peerdatas) {
              let supportSuffixe = ["vtt" /*, 'ass', 'srt', 'webvtt'*/];
              for (let i = 0; i < datas.peerdatas.length; i++) {
                let tpath = datas.peerdatas[i].path;
                if (
                  tpath &&
                  supportSuffixe.indexOf(
                    tpath.getSuffixed(false).toLowerCase()
                  ) > -1 &&
                  tpath.indexOf(fileName) > -1
                ) {
                  subtitleFilePath = datas.peerdatas[i].path;
                  break;
                }
              }
            }
            let loaddingText = "正在加载: " + fileName;
            this.$emit("on-loading", {
              title: loaddingText,
              loaddingText: loaddingText,
            });
            let dpOpts = {
              container:
                document.getElementById("dplayer") || this.$refs["dplayer"].$el,
              video: {
                url: $filepreview.buildStreamURL(token, datas.path.getName()),
                // pic: 'demo.png',
                // thumbnails: 'thumbnails.jpg',
              },
              contextmenu: [],
            };
            // 字幕
            if (subtitleFilePath) {
              dpOpts.subtitle = {
                url: $filepreview.buildStreamURL(
                  token,
                  subtitleFilePath.getName()
                ),
              };
            }
            // 开始播放
            new DPlayer(dpOpts).play();
            loaddingText = fileName;
            this.$emit("on-loading", { title: loaddingText, loadding: false });
          })
          .catch((err) => {
            this.$emit("on-loading", {
              title: "加载失败",
              loaddingText: err.toString(),
            });
          });
      } catch (err) {
        this.$emit("on-loading", {
          title: "加载失败",
          loaddingText: err.toString(),
        });
      }
    },
  },
  created() {
    this.initDatas(this.$route.query.token);
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
