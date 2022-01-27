<template>
  <Form
    name="accountform"
    label-position="left"
    :label-width="80"
    style="padding: 10px"
  >
    <FormItem class="mg_b_10" label="用户名: ">
      <i-input
        v-model="user.userName"
        @on-enter="onUserNameonEnter"
        placeholder="Enter your name"
        style="max-width: 200px"
      ></i-input>
    </FormItem>
    <FormItem class="mg_b_10" label="密码: ">
      <a @click="onUpdataPWD" href="javascript:void(0);">更改密码</a>
    </FormItem>
    <FormItem class="mg_b_10" label="用户类型: ">
      {{ user.userTypeDesc }}
    </FormItem>
    <FormItem class="mg_b_10" label="创建时间: ">
      {{ user.ctTime }}
    </FormItem>
  </Form>
</template>

<script>
import { $apitools } from "../../js/apis/apitools";
import { $userApi } from "../../js/apis/user";
export default {
  name: "MyInfo",
  data() {
    return {
      user: {
        userID: "",
        userType: "",
        userTypeDesc: "",
        userName: "",
        ctTime: "",
      },
      ruleValidate: {
        userName: [{ required: true }],
      },
    };
  },
  methods: {
    onUserNameonEnter() {
      if (!this.user.userName) {
        this.$Message.error("不能为空");
        return;
      }
      $userApi
        .updateUserName(this.user.userID, this.user.userName)
        .then((data) => {
          this.$Message.info("操作成功");
          this.loadCurrentAccountInfo();
        })
        .catch((err) => {
          this.$Message.error("操作失败");
          this.loadCurrentAccountInfo();
        });
    },
    onUpdataPWD() {
      let _ = this;
      let newpwd = "";
      this.$Modal.info({
        title: "更改密码",
        render(h) {
          return h("i-input", {
            props: {
              type: "password",
              value: "",
              autofocus: true,
              placeholder: "Please enter your password...",
            },
            on: {
              input(val) {
                newpwd = val;
              },
            },
          });
        },
        onOk() {
          $userApi
            .updateUserPwd(_.user.userID, newpwd)
            .then((data) => {
              this.$Message.info("操作成功");
            })
            .catch((err) => {
              this.$Message.error("操作失败");
            });
        },
      });
    },
    // 加载当前账户信息
    loadCurrentAccountInfo() {
      let session = $apitools.getSession();
      if (session && session.userID) {
        $userApi
          .queryuser(session.userID)
          .then((data) => {
            if (data) {
              this.user.userID = data.userID;
              this.user.userType = data.userType;
              this.user.userTypeDesc = $userApi.$TYPES.parse(data.userType);
              this.user.userName = data.userName;
              this.user.ctTime = data.ctTime;
              this.$emit("user-changed", this.user);
            }
          })
          .catch((err) => {
            console.error(err);
            this.$Message.error(err.toString());
          });
      }
    },
  },
  created() {
    this.loadCurrentAccountInfo();
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
