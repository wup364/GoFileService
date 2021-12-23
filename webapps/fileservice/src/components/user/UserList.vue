<template>
  <div>
    <div style="padding: 0px 5px 5px 5px">
      <div>
        <Button type="text" icon="md-refresh-circle" @click="loadList"
          >刷新</Button
        >
        <Button
          type="text"
          icon="md-add-circle"
          v-show="shows.add"
          @click="onAddUser"
          >新增</Button
        >
        <Button
          type="text"
          icon="ios-create"
          v-show="shows.edit"
          @click="onEditUser"
          >修改</Button
        >
        <Button
          type="text"
          icon="md-trash"
          v-show="shows.del"
          @click="onDelUser"
          >删除</Button
        >
      </div>
    </div>
    <Table
      ref="umgTb"
      :height="tableHeaght"
      v-auto-height="(n) => (tableHeaght = n - 80)"
      :loading="loading"
      :columns="columns"
      :data="calcPageDatas"
      :highlight-row="true"
      @on-row-click="onRowClick"
      @on-selection-change="onSelectionChange"
    ></Table>
    <Page
      show-total
      :total="datas.length"
      size="small"
      :current="currentPageIndex"
      :page-size="pageSize"
      style="
        line-height: 45px;
        background: #fff;
        float: right;
        padding: 0px 45px;
      "
      @on-change="(val) => (currentPageIndex = val)"
    ></Page>
    <!-- 编辑用户 -->
    <Drawer
      :title="userEdit.isAdd ? '新增用户' : '修改用户'"
      width="450px"
      :closable="false"
      v-model="userEdit.show"
    >
      <Form
        ref="useredit"
        label-position="left"
        :model="userEdit"
        :label-width="80"
      >
        <FormItem class="mg_b_10" label="登录ID: " prop="userID">
          <i-input
            v-if="userEdit.isAdd"
            v-model="userEdit.userID"
            placeholder="Enter your user ID"
          ></i-input>
          <span v-else>{{ userEdit.userID }}</span>
        </FormItem>
        <FormItem class="mg_b_10" label="用户名: " prop="UserName">
          <i-input
            v-model="userEdit.userName"
            placeholder="Enter your name"
          ></i-input>
        </FormItem>
        <FormItem class="mg_b_10" label="密码: " prop="UserPWD">
          <i-input
            v-model="userEdit.userPWD"
            placeholder="Enter your password"
            type="password"
          ></i-input>
        </FormItem>
        <FormItem class="mg_b_10" label="用户类型: " prop="usertype">
          <!-- {{userEdit.UserTypeDesc}} -->
          普通用户
        </FormItem>
      </Form>
      <div>
        <Button class="fr mg_l_10" type="primary" @click="onSaveUser"
          >确定</Button
        >
        <Button class="fr" type="default" @click="userEdit.show = false"
          >关闭</Button
        >
      </div>
    </Drawer>
  </div>
</template>

<script>
import { $userApi } from "../../js/apis/user";
export default {
  name: "UserList",
  data() {
    return {
      loading: false,
      tableHeaght: 500,
      datas: [],
      currentPageIndex: 1,
      pageSize: 60,
      shows: {
        add: true,
        del: false,
        edit: false,
      },
      userEdit: {
        show: false,
        isAdd: false,
        userID: "",
        userType: "",
        userName: "",
        userPWD: "",
      },
      columns: [
        {
          type: "selection",
          width: 60,
          align: "center",
        },
        {
          title: "用户名",
          key: "userName",
        },
        {
          title: "用户类型",
          key: "userTypeDesc",
        },
        {
          title: "新建时间",
          key: "ctTime",
        },
      ],
    };
  },
  computed: {
    calcPageDatas() {
      let _end = this.currentPageIndex * this.pageSize;
      let _start = (this.currentPageIndex - 1) * this.pageSize;
      return this.datas.slice(_start > 0 ? _start : 0, _end);
    },
  },
  methods: {
    // 点击一行
    onRowClick(row, index) {
      this.$refs.umgTb.selectAll(false);
      this.$refs.umgTb.toggleSelect(index);
      this.onSelectionChange([row], row);
    },
    // 当Checkbox数据发生变化, 则需要刷新按钮
    onSelectionChange(selection, row) {
      this.shows.add = true;
      this.shows.del = false;
      this.shows.edit = false;
      let len_selection = selection.length;
      // 选择一个或者以上
      if (len_selection >= 1) {
        this.shows.del = true;
      }
      // 只选择了一个
      if (len_selection == 1) {
        this.shows.edit = true;
      }
    },
    // doListAllUsers
    doListAllUsers() {
      this.datas = [];
      this.loading = true;
      this.currentPageIndex = 1;
      $userApi
        .listAllUsers()
        .then((datas) => {
          if (datas) {
            for (let i = 0; i < datas.length; i++) {
              datas[i].userTypeDesc = $userApi.$TYPES.parse(datas[i].userType);
            }
          }
          this.$set(this, "datas", datas);
          this.loading = false;
        })
        .catch((err) => {
          this.$Message.error(err.toString());
          this.loading = false;
        });
    },
    onAddUser() {
      this.userEdit.isAdd = true;
      this.userEdit.userID = "";
      this.userEdit.userName = "";
      this.userEdit.userType = "";
      this.userEdit.userPWD = "";
      this.userEdit.show = true;
    },
    onDelUser() {
      let rows = this.$refs.umgTb.getSelection();
      if (rows) {
        try {
          for (let i = 0; i < rows.length; i++) {
            $userApi.sync.delUser(rows[i].userID);
          }
          this.$Message.info("操作成功");
        } catch (err) {
          this.$Message.error(err.toString());
        }
        this.doListAllUsers();
      }
    },
    onEditUser() {
      let rows = this.$refs.umgTb.getSelection();
      if (rows && rows.length == 1) {
        this.userEdit.isAdd = false;
        this.userEdit.userPWD = "";
        this.userEdit.userID = rows[0].userID;
        this.userEdit.userName = rows[0].userName;
        this.userEdit.userType = rows[0].userType;
        this.userEdit.show = true;
      }
    },
    onSaveUser() {
      if (!this.userEdit.userID) {
        this.$Message.error("用户ID不能为空");
        return;
      }
      if (!this.userEdit.userName) {
        this.$Message.error("用户名不能为空");
        return;
      }
      if (this.userEdit.isAdd) {
        $userApi
          .addUser(
            this.userEdit.userID,
            this.userEdit.userName,
            this.userEdit.userPWD
          )
          .then((data) => {
            this.$Message.info("操作成功");
            this.doListAllUsers();
            this.userEdit.show = false;
          })
          .catch((err) => {
            this.$Message.error(err.toString());
            this.doListAllUsers();
          });
      } else {
        try {
          $userApi.sync.updateUserName(
            this.userEdit.userID,
            this.userEdit.userName
          );
          $userApi.sync.updateUserPwd(
            this.userEdit.userID,
            this.userEdit.userPWD
          );
          this.userEdit.show = false;
          this.$Message.info("操作成功");
        } catch (err) {
          this.$Message.error(err.toString());
        }
        this.doListAllUsers();
      }
    },
    loadList() {
      this.doListAllUsers();
    },
  },
  created() {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
