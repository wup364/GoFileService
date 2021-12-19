<template>
  <div class="page" ref="main-page">
    <div class="layout">
      <Layout>
        <Sider width="110">
          <div style="padding: 8px; text-align: center">
            <Dropdown @on-click="onDropdownClick" transfer>
              <div>
                <Icon
                  type="md-flower"
                  class="cs_p"
                  :style="{ 'font-size': '48px', color: '#c9ccd0' }"
                ></Icon>
                <a
                  href="javascript:void(0)"
                  class="d_block mg_b_5"
                  style="color: rgb(201, 204, 208); line-height: 20px"
                  >{{ currentAccount.userName }}</a
                >
              </div>
              <Dropdown-Menu slot="list">
                <Dropdown-Item name="logout">注销登录</Dropdown-Item>
              </Dropdown-Menu>
            </Dropdown>
          </div>
          <Menu
            :active-name="currentMenu"
            @on-select="(n) => (currentMenu = n)"
            theme="dark"
            width="auto"
          >
            <Menu-Item name="FileList">
              <Icon type="ios-navigate"></Icon> <span>文件</span>
            </Menu-Item>
            <Menu-Item name="SysSetting">
              <Icon type="ios-settings"></Icon> <span>设置</span>
            </Menu-Item>
          </Menu>
        </Sider>
        <Layout>
          <Content
            :style="{ background: '#fff', minHeight: '220px', height: '100%' }"
          >
            <Component v-bind:is="currentMenu"></Component>
          </Content>
        </Layout>
      </Layout>
    </div>
  </div>
</template>

<script>
import FileList from "./FileList.vue";
import SysSetting from "./SysSetting.vue";
import { $apitools } from "../js/apis/apitools";
import { $userApi } from "../js/apis/user";
export default {
  components: {
    FileList,
    SysSetting,
  },
  name: "Main",
  data() {
    return {
      loading: true,
      currentMenu: "FileList",
      currentAccount: {
        userType: "",
        userName: "",
      },
    };
  },
  methods: {
    onDropdownClick(name) {
      if (name == "logout") {
        $userApi
          .logout()
          .then((data) => {
            $apitools.destroySession();
            window.location.replace("/");
          })
          .catch((err) => {
            this.$Message.error(err.toString());
          });
      }
    },
    // 加载当前账户信息
    doInitCurrentAccountInfo() {
      let session = $apitools.getSession();
      if (session && session.userID) {
        $userApi
          .queryuser(session.userID)
          .then((data) => {
            if (data) {
              this.$set(this, "currentAccount", JSON.parse(data));
            }
          })
          .catch((err) => {
            this.$Message.error(err.toString());
          });
      }
    },
  },
  created() {
    this.doInitCurrentAccountInfo();
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
</style>
