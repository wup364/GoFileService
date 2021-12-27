<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
import _ from "./js/vue/vue.components.js";
import { $apitools } from "./js/apis/apitools";
export default {
  name: "App",
  created() {
    $apitools.setAPIHost("http://127.0.0.1:8080");
    this.$nextTick(this.stopDrop);
  },
  methods: {
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
};
</script>
<style>
#app {
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding: 0px;
  margin: 0px;
}
.page {
  width: 100%;
  height: 100%;
  padding: 0px;
  margin: 0px;
}
.ivu-layout {
  background: #fff;
  overflow: hidden;
  height: 100%;
  width: 100%;
  margin: 0px;
  padding: 0px;
}
</style>
