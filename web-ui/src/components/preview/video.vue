<template>
  <div class="page" ref="pw-video-page">
    <div ref="dplayer" id="dplayer" style="height: 100%; width: 100%"></div>
  </div>
</template>

<script>
import "../../js/3party/dplayer/DPlayer.min.js";
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
            let subtitleFilePath = "";
            let fileName = datas.path.getName(false);
            if (datas && datas.peerDatas) {
              let supportSuffixe = ["vtt" /*, 'ass', 'srt', 'webvtt'*/];
              for (let i = 0; i < datas.peerDatas.length; i++) {
                let tpath = datas.peerDatas[i].path;
                if (
                  tpath &&
                  supportSuffixe.indexOf(
                    tpath.getSuffixed(false).toLowerCase()
                  ) > -1 &&
                  tpath.indexOf(fileName) > -1
                ) {
                  subtitleFilePath = datas.peerDatas[i].path;
                  break;
                }
              }
            }
            let loaddingText = "正在加载: " + fileName;
            this.$emit("on-loading", {
              title: loaddingText,
              loaddingText: loaddingText,
            });
            let filenames = [datas.path.getName()];
            if (subtitleFilePath) {
              filenames.push(subtitleFilePath.getName());
            }
            let tokendatas = $filepreview.sync.samedirtoken(token, filenames);
            let dpOpts = {
              container:
                document.getElementById("dplayer") || this.$refs["dplayer"].$el,
              video: {
                url: tokendatas[datas.path.getName()].tokenURL,
                // pic: 'demo.png',
                // thumbnails: 'thumbnails.jpg',
              },
              contextmenu: [],
            };
            // 字幕
            if (subtitleFilePath) {
              dpOpts.subtitle = {
                url: tokendatas[subtitleFilePath.getName()].tokenURL,
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
