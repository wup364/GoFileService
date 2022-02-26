<template>
  <div class="page" ref="pw-picture-page">
    <div class="container">
      <Stick :list="list" @onScrollEnd="doLoadMore">
        <template slot-scope="scope">
          <div class="card cs_p" @click="onShowBig(scope.data)">
            <img :src="scope.data.cover" v-if="scope.data.cover" />
            <h2>{{ scope.data.title }}</h2>
            <p>{{ scope.data.intro }}</p>
          </div>
        </template>
      </Stick>
      <div class="tools">
        <div v-show="isLoading">正在加载...</div>
        <div v-if="isEnd">没有更多数据了</div>
        <a style="color: #555" v-if="!isEnd && !isLoading" @click="doLoadMore">
          加载更多
        </a>
      </div>
    </div>
  </div>
</template>

<script>
import Stick from "vue-stick";
import { $filepreview } from "../../js/apis/filepreview";
export default {
  components: {
    Stick: Stick.component,
  },
  name: "",
  data() {
    return {
      isLoading: false,
      allDatas: { cindex: 0, list: [] },
      list: [],
      pageIndex: -1,
      pageSize: 10,
    };
  },
  computed: {
    isEnd() {
      return this.allDatas.list.length <= this.list.length;
    },
  },
  methods: {
    onShowBig(item) {
      window.open(item.cover);
    },
    doLoadMore() {
      if (this.allDatas.list.length > 0) {
        this.isLoading = true;
        let startIndex = ++this.pageIndex * this.pageSize;
        let stopIndex = (this.pageIndex + 1) * this.pageSize;
        let newList = this.allDatas.list.slice(startIndex, stopIndex);
        //
        let names = [];
        let token = this.$route.query.token;
        newList.forEach((row) => {
          if (row.intro) {
            names.push(row.intro);
          }
        });
        $filepreview
          .samedirtoken(token, names)
          .then((datas) => {
            newList.forEach((row) => {
              if (row.intro && datas[row.intro]) {
                row.cover = datas[row.intro].tokenURL;
              } else {
                row.cover = "/static/img/preview/loading.gif";
              }
            });
            this.list = this.list.concat(newList);
            this.isLoading = false;
          })
          .catch((err) => {
            console.error(err);
            this.isLoading = false;
          });
      }
    },
    // 构建list
    doBuilList(datas) {
      let res = {
        cindex: -1,
        list: [],
        count: -1,
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
            if (!$filepreview.isSupport("picture", px)) {
              continue;
            }
            if (node.path == cpath) {
              res.cindex = 0; // res.list.length;
            }
            if (res.cindex == -1) {
              continue;
            }
            res.list.push({
              cover: "",
              title: "SEQ-" + (i + 1),
              intro: node.path.getName(),
            });
          }
        }
      }
      return res;
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
            // 开始播放
            this.allDatas = this.doBuilList(datas);
            this.doLoadMore();
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
  beforeCreate() {
    document.querySelector("body").setAttribute("style", "height:unset;");
  },
  created() {
    this.initDatas(this.$route.query.token);
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.page {
  background: #333;
  overflow: auto;
}
.container {
  position: relative;
  min-height: 1000px;
  margin: 20px 40px 40px 40px;
}
.tools {
  padding: 20px;
  text-align: center;
  background: rgb(148, 147, 147);
}

.card {
  background: #fff;
}

.card img {
  display: block;
  width: 100%;
  background: #aaa;
}

.card h2 {
  margin: 0;
  padding: 5px 15px;
  font-size: 14px;
}

.card p {
  margin: 0;
  padding: 10px 15px;
  font-size: 14px;
}
</style>
