<template>
  <breadcrumb separator=">">
    <BreadcrumbItem>
      <a @click="address_GoToRoot()">{{ showrootname }}</a>
    </BreadcrumbItem>
    <BreadcrumbItem
      v-for="(item, index) in paths"
      v-bind:key="index"
      v-show="item && (index >= paths.length - 2 || index <= max)"
    >
      <a @click="address_GoToPath(item, index)">{{ item }}</a>
    </BreadcrumbItem>
  </breadcrumb>
</template>

 
<script>
export default {
  name: "fileaddress",
  props: ["path", "root", "rootname", "depth"],
  data() {
    return {
      paths: [],
      max: 6,
      showrootname: "",
    };
  },
  created() {
    this.max = this.depth
      ? this.depth - 2 > 0
        ? this.depth - 2
        : 2
      : this.max;
    this.buildPaths();
  },
  methods: {
    buildPaths() {
      this.showrootname = this.rootname ? this.rootname : "/";
      this.paths = this.path.split("/");
    },
    address_GoToRoot() {
      this.$emit("click", this.root ? this.root : "/");
    },
    address_GoToPath(item, index) {
      let path = "";
      for (let i = 0; i <= index; i++) {
        if (this.paths[i]) {
          path += "/" + this.paths[i];
        }
      }
      this.$emit("click", path);
    },
  },
  watch: {
    path(v1, v2) {
      this.buildPaths();
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
