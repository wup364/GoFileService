<template>
  <div v-watch-height="watchHeight" class="page" ref="pw-audio-page">
    <div
      id="bg"
      style="
        position: fixed;
        background-size: contain;
        background-position: center center;
        filter: blur(1px);
        width: 100%;
        height: 100%;
      "
    ></div>
    <aplayer
      v-if="currentMusic"
      :repeat="repeatMethod"
      :list="musicList"
      :music="currentMusic"
      :list-max-height="listMaxHeight"
      v-on:update:music="onMusicUpdate"
      ref="aplayer"
      id="aplayer"
    />
  </div>
</template>

<script>
import "../../js/3party/jsmediatags/jsmediatags.min.js";
import { $filepreview } from "../../js/apis/filepreview";
import { $utils } from "../../js/utils";
import Aplayer from "vue-aplayer";
Aplayer.disableVersionBadge = true;

export default {
  components: {
    Aplayer,
  },
  name: "pw-audio",
  data() {
    return {
      musicList: [],
      currentMusic: undefined,
      repeatMethod: "list", // music|list|none
      listMaxHeight: "100%",
    };
  },
  methods: {
    //
    watchHeight(ch, ph) {
      try {
        this.listMaxHeight =
          ph -
          document.getElementsByClassName("aplayer-body")[0].clientHeight +
          "px";
      } catch (error) {
        this.listMaxHeight = ph - 66 + "px";
      }
    },
    // 切换播放列表
    onMusicUpdate(music) {
      if (!music._cover_) {
        this.getCover(music, () => {
          music._cover_ = true;
          document.getElementById("bg").style.backgroundImage =
            "url('" + music.pic + "')";
        });
      } else if (music.pic) {
        document.getElementById("bg").style.backgroundImage =
          "url('" + music.pic + "')";
      }
    },

    // 加载数据列表
    buildMusicList(token, datas) {
      let res = {
        cindex: 0,
        list: [],
      };
      if (datas) {
        let cpath = datas.path;
        if (datas.peerDatas && datas.peerDatas.length > 0) {
          for (let i = 0; i < datas.peerDatas.length; i++) {
            let node = datas.peerDatas[i];
            if (!node.isFile) {
              continue;
            }
            let px = node.path.getSuffixed(false);
            if (!$filepreview.isSupport("audio", px)) {
              continue;
            }
            if (node.path == cpath) {
              res.cindex = res.list.length;
            }
            res.list.push({
              title: node.path.getName(false),
              artist: $utils.formatSize(node.size),
              src: $filepreview.buildStreamURL(token, node.path.getName()),
              pic: "/static/img/preview/default_cover.png",
              lrc: "",
            });
          }
        }
      }
      return res;
    },
    // 获取封面
    getCover(audio, cb) {
      let _ = this;
      jsmediatags.read(audio.src, {
        onSuccess(datas) {
          let tags = datas.tags;
          audio.artist = tags.artist ? tags.artist : "";
          if (tags.picture) {
            audio.pic =
              "data:" +
              tags.picture.format +
              ";base64," +
              _.arrayBufferToBase64(tags.picture.data);
          } else {
            audio.pic = "/static/img/preview/default_cover.png";
          }
          if (cb) {
            cb();
          }
        },
        onError(err) {
          if (cb) {
            cb(err);
          }
        },
      });
    },
    // Buffer数据转base64字符
    arrayBufferToBase64(buffer) {
      let binary = "";
      let bytes = new Uint8Array(buffer);
      let len = bytes.byteLength;
      for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
      }
      return window.btoa(binary);
    },
    initDatas(token) {
      try {
        $filepreview
          .samedirFiles(token)
          .then((datas) => {
            //
            let loaddingText = "正在加载: " + datas.path.getName(false);
            this.$emit("on-loading", {
              title: loaddingText,
              loaddingText: loaddingText,
            });
            // 初始胡aplayer
            let mlist = this.buildMusicList(token, datas);
            if (mlist.list.length > 0) {
              this.musicList = mlist.list;
              this.currentMusic = mlist.list[mlist.cindex];
              this.onMusicUpdate(this.currentMusic);
              this.$nextTick(() => {
                let ap = this.$refs.aplayer;
                ap.play();
              });
            }
            // 刷新顶部title
            this.$emit("on-loading", {
              title: datas.path.getName(false),
              loadding: false,
            });
          })
          .catch((err) => {
            this.$emit("on-loading", {
              title: "加载失败",
              loaddingText: err.toString(),
            });
            console.error(err);
          });
      } catch (err) {
        this.$emit("on-loading", {
          title: "加载失败",
          loaddingText: err.toString(),
        });
        console.error(err);
      }
    },
  },
  created() {
    this.initDatas(this.$route.query.token);
  },
  watch: {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
body {
  background-repeat: no-repeat;
  background-size: cover;
  overflow-y: auto;
}
#aplayer {
  max-width: 1000px;
  margin: auto;
  opacity: 0.7;
}
</style >
