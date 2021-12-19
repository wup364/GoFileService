<template>
  <div class="page" ref="login-page">
    <div>
      <div
        ref="loginwarp"
        class="login-warp"
        :style="{
          'margin-top': loginform_margin_top + 'px',
          'margin-left': loginform_margin_left + 'px',
        }"
        @keyup.enter="handleSubmit"
      >
        <Form :model="loginform">
          <FormItem>
            <span style="font-size: 20px; color: #666">登录系统</span>
          </FormItem>
          <FormItem>
            <Input type="text" v-model="loginform.user" placeholder="用户">
              <Icon type="ios-person-outline" slot="prepend"></Icon>
            </Input>
          </FormItem>
          <FormItem prop="password">
            <Input
              type="password"
              v-model="loginform.password"
              placeholder="密码"
            >
              <Icon type="ios-lock-outline" slot="prepend"></Icon>
            </Input>
          </FormItem>
          <FormItem>
            <Button
              type="primary"
              :loading="loading"
              @click="handleSubmit"
              style="width: 100%"
              >登录</Button
            >
          </FormItem>
        </Form>
      </div>
    </div>
    <div
      id="bg"
      style="width: 100%; height: 100%; position: absolute; top: 0px"
    ></div>
  </div>
</template>

 
<script>
import { Victor } from "../js/3party/vector";
import { $utils } from "../js/utils";
import { $apitools } from "../js/apis/apitools";
import { $userApi } from "../js/apis/user";
export default {
  name: "LogIn",
  data() {
    return {
      loading: false,
      loginform_margin_top: 100,
      loginform_margin_left: 100,
      loginform: {
        user: "",
        password: "",
      },
      message: "正在登录,请稍后...",
    };
  },
  methods: {
    // 登录
    handleSubmit() {
      if (!$utils.supportES6()) {
        alert("当前浏览器不受支持, 推荐: edge, chrome, firefox等现代浏览器");
        return;
      }
      this.loading = true;
      $userApi
        .login(this.loginform.user, this.loginform.password)
        .then((data) => {
          this.loading = false;
          $apitools.saveSession(JSON.parse(data));
          this.$router.push("/main");
        })
        .catch((err) => {
          console.error(err);
          this.loading = false;
          this.$Message.error("登录失败");
        });
    },
    //
    doAutoHeight(notListen) {
      this.$nextTick(() => {
        if (this.$refs["login-page"] && this.$refs.loginwarp) {
          this.loginform_margin_top =
            this.$refs["login-page"].clientHeight / 2 -
            this.$refs.loginwarp.clientHeight / 2 -
            50;
          this.loginform_margin_left =
            this.$refs["login-page"].clientWidth / 2 -
            this.$refs.loginwarp.clientWidth / 2;
        }
      });
      if (!notListen) {
        window.addEventListener("resize", () => {
          this.doAutoHeight(true);
        });
      }
    },
    doRenderbg() {
      this.$nextTick(() => {
        Victor("bg");
      });
    },
  },
  created() {
    this.doAutoHeight();
    this.doRenderbg();
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.login-warp {
  text-align: center;
  font-size: 40px;
  width: 350px;
  background: #fff;
  z-index: 1;
  position: absolute;
  padding: 25px 15px;
  border-radius: 5px;
}
</style>
