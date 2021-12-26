<template>
  <div class="sys-setting" ref="syssetting-page">
    <Tabs
      v-model="currentTab"
      style="height: 100%"
      v-watch-height="(h) => (tabHeight = h - 52)"
    >
      <TabPane
        name="currentAccount"
        label="账户信息"
        :style="{ height: tabHeight + 'px' }"
      >
        <!-- 账户信息 -->
        <my-info @user-changed="(n) => $set(this, 'currentAccount', n)" />
      </TabPane>
      <TabPane
        v-if="currentAccount.userType === 0"
        name="userManage"
        label="账户管理"
        :style="{ height: tabHeight + 'px' }"
      >
        <!-- 账户管理 -->
        <user-list ref="userManage" style="width: 100%; height: 100%" />
      </TabPane>
      <TabPane
        v-if="currentAccount.userType === 0"
        name="fPermissionManage"
        label="文件权限管理"
        :style="{ height: tabHeight + 'px' }"
      >
        <!-- 权限管理 -->
        <file-pms-list ref="fPermissionManage" />
      </TabPane>
    </Tabs>
  </div>
</template>

<script>
import MyInfo from "../components/user/MyInfo.vue";
import UserList from "../components/user/UserList.vue";
import FilePmsList from "../components/permission/FPMList.vue";
export default {
  name: "SysSetting",
  components: {
    MyInfo,
    UserList,
    FilePmsList,
  },
  data() {
    let _ = this;
    return {
      tabHeight: 500,
      currentTab: "currentAccount",
      // 账户信息
      currentAccount: {
        userID: "",
        userType: "",
        userName: "",
      },
    };
  },

  methods: {},
  created() {},
  watch: {
    currentTab(n, o) {
      if (n == "userManage") {
        this.$refs["userManage"].loadList();
      } else if (n == "fPermissionManage") {
        this.$refs["fPermissionManage"].loadList();
      }
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.ivu-tabs-bar {
  margin-bottom: 8px;
}
.ivu-table-wrapper {
  border: none;
}
.ivu-table:after {
  content: none;
}
.sys-setting {
  height: 100%;
  padding-top: 7px;
}
</style>
