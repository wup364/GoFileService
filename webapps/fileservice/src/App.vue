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
  },
};
</script>

<style>
</style>
