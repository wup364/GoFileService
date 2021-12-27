<template>
  <div class="page" ref="pw-audio-page">
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
    <div ref="aplayer" id="aplayer"></div>
  </div>
</template>

<script>
import "aplayer/dist/APlayer.min.css";
import APlayer from "aplayer";
import _ from "../../js/3party/jsmediatags/jsmediatags.min.js";
import { $filepreview } from "../../js/apis/filepreview";
import { $utils } from "../../js/utils";
export default {
  name: "pw-audio",
  data() {
    return {};
  },
  methods: {
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
              name: node.path.getName(false),
              artist: $utils.formatSize(node.size),
              url: $filepreview.buildStreamURL(token, node.path.getName()),
              cover: "/static/img/preview/default_cover.png",
            });
          }
        }
      }
      return res;
    },
    // 获取封面
    getCover(audio, cb) {
      let _ = this;
      jsmediatags.read(audio.url, {
        onSuccess(datas) {
          // console.log(datas, datas.tags);
          let tags = datas.tags;

          audio.artist = tags.artist ? tags.artist : "";
          if (tags.picture) {
            audio.cover =
              "data:" +
              tags.picture.format +
              ";base64," +
              _.arrayBufferToBase64(tags.picture.data);
          } else {
            audio.cover = "/img/preview/default_cover.png";
          }
          if (cb) {
            cb();
          }
        },
        onError(err) {
          // console.log(err);
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
    //
    initAPlayer(mlist) {
      let _ = this;
      let ap = new APlayer({
        audio: mlist,
        listMaxHeight: 0,
        container: document.getElementById("aplayer") || _.$refs["aplayer"].$el,
      });
      ap.on("listshow", (e) => {
        document.getElementById("aplayer").style.height = "100%";
      });
      ap.on("listhide", (e) => {
        document.getElementById("aplayer").style.height = "66px";
      });
      ap.on("listswitch", (e) => {
        let mitem = mlist[e.index];
        mitem.artist = undefined;
        this.$emit("on-loading", { title: mitem.name, loadding: false });
        if (!mitem["_cover_"]) {
          this.getCover(mitem, () => {
            mitem["_cover_"] = true;
            ap.template.pic.style.backgroundImage = mitem.cover
              ? "url('" + mitem.cover + "')"
              : "";
            ap.template.author.innerHTML = mitem.artist ? mitem.artist : "";
            document.getElementById("bg").style.backgroundImage =
              "url('" + mitem.cover + "')";
          });
        } else {
          if (mitem.cover) {
            document.getElementById("bg").style.backgroundImage =
              "url('" + mitem.cover + "')";
          }
        }
      });
      return ap;
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
            let ap = this.initAPlayer(mlist.list);
            if (mlist.list.length > 0) {
              ap.list.switch(mlist.cindex);
              ap.play();
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
  height: 100%;
  opacity: 0.8;
}
</style >
<style >
#aplayer .aplayer-list {
  overflow-y: auto !important;
  height: calc(100% - 70px) !important;
}
</style>
