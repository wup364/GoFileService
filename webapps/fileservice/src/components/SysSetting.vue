<template>
  <div class="page" ref="main-page">
    <div class="sys-setting" style="height: 100%; padding-top: 7px">
      <Tabs class="usermanagetab" v-model="currentTab" style="height: 100%">
        <!-- 账户信息 -->
        <Tab-Pane name="currentAccount" label="账户信息">
          <!-- <div class="mg_0_auto" style="max-width: 300px;"> -->
          <From
            ref="currentAccount"
            label-position="left"
            :model="currentAccount"
            :label-width="80"
            style="padding: 10px"
          >
            <Form-Item class="mg_b_10" label="用户名: " prop="UserName">
              <i-input
                v-model="currentAccount.UserName"
                @on-enter="currentAccount.onUserNameonEnter"
                placeholder="Enter your name"
                style="max-width: 200px"
              ></i-input>
            </Form-Item>
            <Form-Item class="mg_b_10" label="密码: " prop="UserPWD">
              <a @click="currentAccount.onUpdataPWD" href="javascript:void(0);"
                >更改密码</a
              >
            </Form-Item>
            <Form-Item class="mg_b_10" label="用户类型: " prop="usertype">
              {{ currentAccount.UserTypeDesc }}
            </Form-Item>
            <Form-Item class="mg_b_10" label="创建时间: " prop="cttime">
              {{ currentAccount.CtTime }}
            </Form-Item>
          </From>
          <!-- </div> -->
        </Tab-Pane>
        <!-- 账户管理 -->
        <Tab-Pane
          v-if="currentAccount.UserType == '1'"
          name="userManage"
          label="账户管理"
        >
          <div style="padding: 0px 5px 5px 5px">
            <div>
              <Button
                type="text"
                icon="md-add-circle"
                v-show="userManage.shows.add"
                @click="onAddUser"
                >新增</Button
              >
              <Button
                type="text"
                icon="ios-create"
                v-show="userManage.shows.edit"
                @click="onEditUser"
                >修改</Button
              >
              <Button
                type="text"
                icon="md-trash"
                v-show="userManage.shows.del"
                @click="onDelUser"
                >删除</Button
              >
            </div>
          </div>
          <Table
            ref="userManage"
            v-minus-height="80"
            :loading="userManage.loading"
            :columns="userManage.columns"
            :data="calcUserManagePageDatas"
            :highlight-row="true"
            @on-row-click="userManage.onRowClick"
            @on-selection-change="userManage.onSelectionChange"
          ></Table>
          <Page
            show-total
            :total="userManage.datas.length"
            size="small"
            :current="userManage.currentPageIndex"
            :page-size="userManage.pageSize"
            style="
              line-height: 45px;
              background: #fff;
              float: right;
              padding: 0px 45px;
            "
            @on-change="(val) => (userManage.currentPageIndex = val)"
          ></Page>
        </Tab-Pane>
        <!-- 权限管理 -->
        <Tab-Pane
          v-if="currentAccount.UserType == '1'"
          name="fPermissionManage"
          label="文件权限管理"
        >
          <div style="padding: 0px 5px 5px 5px">
            <div>
              <Button
                type="text"
                icon="md-add-circle"
                v-show="fPermissionManage.shows.add"
                @click="onAddFPermission"
                >新增</Button
              >
              <Button
                type="text"
                icon="ios-create"
                v-show="fPermissionManage.shows.edit"
                @click="onEditFPermission"
                >修改</Button
              >
              <Button
                type="text"
                icon="md-trash"
                v-show="fPermissionManage.shows.del"
                @click="onDelFPermission"
                >删除</Button
              >
            </div>
          </div>
          <Table
            ref="fPermissionManage"
            v-minus-height="80"
            :loading="fPermissionManage.loading"
            :columns="fPermissionManage.columns"
            :data="calcFPermissionManage"
            :highlight-row="true"
            @on-row-click="fPermissionManage.onRowClick"
            @on-selection-change="fPermissionManage.onSelectionChange"
          ></Table>
          <Page
            show-total
            :total="fPermissionManage.datas.length"
            size="small"
            :current="fPermissionManage.currentPageIndex"
            :page-size="fPermissionManage.pageSize"
            style="
              line-height: 45px;
              background: #fff;
              float: right;
              padding: 0px 45px;
            "
            @on-change="(val) => (fPermissionManage.currentPageIndex = val)"
          ></Page>
        </Tab-Pane>
      </Tabs>
      <!-- 编辑用户 -->
      <Drawer
        :title="userManage.userEdit.isAdd ? '新增用户' : '修改用户'"
        width="450px"
        :closable="false"
        v-model="userManage.userEdit.show"
      >
        <From
          ref="useredit"
          label-position="left"
          :model="userManage.userEdit"
          :label-width="80"
        >
          <Form-Item class="mg_b_10" label="登录ID: " prop="userID">
            <i-input
              v-if="userManage.userEdit.isAdd"
              v-model="userManage.userEdit.userID"
              placeholder="Enter your user ID"
            ></i-input>
            <span v-else>{{ userManage.userEdit.userID }}</span>
          </Form-Item>
          <Form-Item class="mg_b_10" label="用户名: " prop="UserName">
            <i-input
              v-model="userManage.userEdit.UserName"
              placeholder="Enter your name"
            ></i-input>
          </Form-Item>
          <Form-Item class="mg_b_10" label="密码: " prop="UserPWD">
            <i-input
              v-model="userManage.userEdit.UserPWD"
              placeholder="Enter your password"
              type="password"
            ></i-input>
          </Form-Item>
          <Form-Item class="mg_b_10" label="用户类型: " prop="usertype">
            <!-- {{userManage.userEdit.UserTypeDesc}} -->
            普通用户
          </Form-Item>
        </From>
        <div>
          <Button class="fr mg_l_10" type="primary" @click="onSaveUser"
            >确定</Button
          >
          <Button
            class="fr"
            type="default"
            @click="userManage.userEdit.show = false"
            >关闭</Button
          >
        </div>
      </Drawer>
      <!-- 编辑权限 -->
      <Drawer
        :title="
          fPermissionManage.fPermissionEdit.isAdd ? '新增权限' : '修改权限'
        "
        width="450px"
        v-model="fPermissionManage.fPermissionEdit.show"
      >
        <From
          ref="useredit"
          label-position="left"
          :model="fPermissionManage.fPermissionEdit"
          :label-width="80"
        >
          <Form-Item class="mg_b_10" label="用户ID: " prop="userID">
            <Select
              v-if="fPermissionManage.fPermissionEdit.isAdd"
              v-model="fPermissionManage.fPermissionEdit.userID"
              transfer
            >
              <Option
                v-for="item in fPermissionManage.userselector.datas"
                :value="item.userID"
                :key="item.userID"
                >{{ item.UserName }}</Option
              >
            </Select>
            <span v-else>{{ fPermissionManage.fPermissionEdit.userID }}</span>
          </Form-Item>
          <Form-Item class="mg_b_10" label="授权目录: " prop="path">
            <i-input
              v-if="fPermissionManage.fPermissionEdit.isAdd"
              v-model="fPermissionManage.fPermissionEdit.path"
              placeholder="授权路径"
              readonly
            >
              <Button
                slot="append"
                @click="onSelectDir('onBeforFpmsSelectedDir')"
                >选择</Button
              >
            </i-input>
            <span v-else>{{ fPermissionManage.fPermissionEdit.path }}</span>
          </Form-Item>
          <Form-Item class="mg_b_10" label="授权权限: " prop="Permission">
            <!-- <i-input v-model="fPermissionManage.fPermissionEdit.Permission" placeholder=""></i-input> -->
            <Checkbox-Group
              v-model="fPermissionManage.fPermissionEdit.Permissions"
            >
              <Checkbox label="1">Visible</Checkbox>
              <Checkbox label="2">Read</Checkbox>
              <Checkbox label="3">Write</Checkbox>
            </Checkbox-Group>
          </Form-Item>
        </From>
        <div>
          <Button class="fr mg_l_10" type="primary" @click="onSaveFPermission"
            >确定</Button
          >
          <Button
            class="fr"
            type="default"
            @click="fPermissionManage.fPermissionEdit.show = false"
            >关闭</Button
          >
        </div>
      </Drawer>
      <!-- 选择目录 -->
      <fs-selector
        :select-muti="fsSelector.selectMuti"
        :select-file="fsSelector.selectFile"
        :select-dir="fsSelector.selectDir"
        :start-path="fsSelector.startPath"
        :show-dailog="fsSelector.showDailog"
        @on-select="onSelectedFile"
        @on-cancel="onSelectCancel"
      ></fs-selector>
    </div>
  </div>
</template>

<script>
import { $utils } from "../js/utils";
import { $filepms } from "../js/apis/filepermission";
export default {
  name: "SysSetting",
  data: function () {
    let _ = this;
    return {
      UserTypes: {
        1: "管理员",
        0: "普通用户",
      },
      fsSelector: {
        operationObj: false,
        showDailog: false,
        selectFile: true,
        selectDir: true,
        selectMuti: true,
        startPath: "/",
      },
      currentTab: "currentAccount",
      // 账户管理
      userManage: {
        loading: false,
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
          UserType: "",
          UserName: "",
          UserPWD: "",
        },
        columns: [
          {
            type: "selection",
            width: 60,
            align: "center",
          },
          {
            title: "用户名",
            key: "UserName",
          },
          {
            title: "用户类型",
            key: "UserTypeDesc",
          },
          {
            title: "新建时间",
            key: "CtTime",
          },
        ],
        // 点击一行
        onRowClick: function (row, index) {
          _.$refs.userManage.selectAll(false);
          _.$refs.userManage.toggleSelect(index);
          _.userManage.onSelectionChange([row], row);
        },
        // 当Checkbox数据发生变化, 则需要刷新按钮
        onSelectionChange: function (selection, row) {
          _.userManage.shows.add = true;
          _.userManage.shows.del = false;
          _.userManage.shows.edit = false;
          let len_selection = selection.length;
          // 选择一个或者以上
          if (len_selection >= 1) {
            _.userManage.shows.del = true;
          }
          // 只选择了一个
          if (len_selection == 1) {
            _.userManage.shows.edit = true;
          }
        },
      },
      // 文件权限管理
      fPermissionManage: {
        loading: false,
        datas: [],
        currentPageIndex: 1,
        pageSize: 60,
        shows: {
          add: true,
          del: false,
          edit: false,
        },
        permissionMax: 3, // 权限最大数
        fPermissionEdit: {
          show: false,
          isAdd: false,
          userID: "",
          path: "",
          Permission: 0,
          Permissions: [],
          PermissionID: "",
        },
        // 临时用下拉-后期优化
        userselector: {
          datas: [],
        },
        columns: [
          {
            type: "selection",
            width: 60,
            align: "center",
          },
          {
            title: "用户ID",
            key: "userID",
          },
          {
            title: "文件路径",
            key: "path",
          },
          {
            title: "权限类型",
            key: "Permission",
            render: function (h, c) {
              return h(
                "span",
                {},
                $filepms.$TYPE.sum2Name(c.row.Permission).join(", ")
              );
            },
          },
          {
            title: "创建时间",
            key: "CtTime",
          },
        ],
        // 点击一行
        onRowClick: function (row, index) {
          _.$refs.fPermissionManage.selectAll(false);
          _.$refs.fPermissionManage.toggleSelect(index);
          _.fPermissionManage.onSelectionChange([row], row);
        },
        // 当Checkbox数据发生变化, 则需要刷新按钮
        onSelectionChange: function (selection, row) {
          _.fPermissionManage.shows.add = true;
          _.fPermissionManage.shows.del = false;
          _.fPermissionManage.shows.edit = false;
          let len_selection = selection.length;
          // 选择一个或者以上
          if (len_selection >= 1) {
            _.fPermissionManage.shows.del = true;
          }
          // 只选择了一个
          if (len_selection == 1) {
            _.fPermissionManage.shows.edit = true;
          }
        },
      },
      // 账户信息
      currentAccount: {
        userID: "",
        UserType: "",
        UserTypeDesc: "",
        UserName: "",
        UserPWD: "",
        CtTime: "",
        ruleValidate: {
          UserName: [{ required: true }],
        },
        onUserNameonEnter: function () {
          if (!_.currentAccount.UserName) {
            _.$Message.error("不能为空");
            return;
          }
          $userApi
            .updateUserName(_.currentAccount.userID, _.currentAccount.UserName)
            .then(function (data) {
              _.$Message.info("操作成功");
              _.$root.doInitCurrentAccountInfo();
              // _.doBuildAccountInfo();
            })
            .catch(function (err) {
              _.$Message.error("操作失败");
              _.$root.doInitCurrentAccountInfo();
              // _.doBuildAccountInfo();
            });
        },
        onUpdataPWD: function () {
          let newpwd = "";
          _.$Modal.info({
            title: "更改密码",
            render: function (h) {
              return h("i-input", {
                props: {
                  type: "password",
                  value: "",
                  autofocus: true,
                  placeholder: "Please enter your password...",
                },
                on: {
                  input: function (val) {
                    newpwd = val;
                  },
                },
              });
            },
            onOk: function () {
              $userApi
                .updateUserPwd(_.currentAccount.userID, newpwd)
                .then(function (data) {
                  _.$Message.info("操作成功");
                })
                .catch(function (err) {
                  _.$Message.error("操作失败");
                });
            },
          });
        },
      },
    };
  },
  computed: {
    calcUserManagePageDatas: function () {
      let _end = this.userManage.currentPageIndex * this.userManage.pageSize;
      let _start =
        (this.userManage.currentPageIndex - 1) * this.userManage.pageSize;
      return this.userManage.datas.slice(_start > 0 ? _start : 0, _end);
    },
    calcFPermissionManage: function () {
      let _end =
        this.fPermissionManage.currentPageIndex *
        this.fPermissionManage.pageSize;
      let _start =
        (this.fPermissionManage.currentPageIndex - 1) *
        this.fPermissionManage.pageSize;
      return this.fPermissionManage.datas.slice(_start > 0 ? _start : 0, _end);
    },
  },
  methods: {
    // 选择路径
    onSelectDir: function (cb) {
      this.fsSelector.operationObj = cb;
      this.fsSelector.selectFile = false;
      this.fsSelector.selectDir = true;
      this.fsSelector.selectMuti = false;
      this.fsSelector.startPath = "/";
      this.fsSelector.showDailog = true;
    },
    onSelectCancel: function () {
      this.fsSelector.showDailog = false;
    },
    onSelectedFile: function (rows) {
      this.fsSelector.showDailog = false;
      if (this.fsSelector.operationObj) {
        this[this.fsSelector.operationObj](rows);
      }
    },
    // 加载账户信息
    doBuildAccountInfo: function (data) {
      if (data) {
        data.UserTypeDesc = this.UserTypes[data.UserType]
          ? this.UserTypes[data.UserType]
          : "未知类型";
      }
    },
    // doListAllUsers
    doListAllUsers: function () {
      let _ = this;
      _.userManage.datas = [];
      _.userManage.loading = true;
      _.userManage.currentPageIndex = 1;
      $userApi
        .listAllUsers()
        .then(function (datas) {
          if (datas) {
            datas = JSON.parse(datas);
            for (let i = 0; i < datas.length; i++) {
              _.doBuildAccountInfo(datas[i]);
            }
          }
          _.$set(_.userManage, "datas", datas);
          _.userManage.loading = false;
        })
        .catch(function (err) {
          _.$Message.error(err.toString());
          _.userManage.loading = false;
        });
    },
    onAddUser: function () {
      this.userManage.userEdit.isAdd = true;
      this.userManage.userEdit.userID = "";
      this.userManage.userEdit.UserName = "";
      this.userManage.userEdit.UserType = "";
      this.userManage.userEdit.UserPWD = "";
      this.userManage.userEdit.show = true;
    },
    onDelUser: function () {
      let rows = this.$refs.userManage.getSelection();
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
    onEditUser: function () {
      let rows = this.$refs.userManage.getSelection();
      if (rows && rows.length == 1) {
        this.userManage.userEdit.isAdd = false;
        this.userManage.userEdit.UserPWD = "";
        this.userManage.userEdit.userID = rows[0].userID;
        this.userManage.userEdit.UserName = rows[0].UserName;
        this.userManage.userEdit.UserType = rows[0].UserType;
        this.userManage.userEdit.show = true;
      }
    },
    onSaveUser: function () {
      if (!this.userManage.userEdit.userID) {
        this.$Message.error("用户ID不能为空");
        return;
      }
      if (!this.userManage.userEdit.UserName) {
        this.$Message.error("用户名不能为空");
        return;
      }
      if (this.userManage.userEdit.isAdd) {
        let _ = this;
        $userApi
          .addUser(
            this.userManage.userEdit.userID,
            this.userManage.userEdit.UserName,
            this.userManage.userEdit.UserPWD
          )
          .then(function (data) {
            _.$Message.info("操作成功");
            _.doListAllUsers();
            _.userManage.userEdit.show = false;
          })
          .catch(function (err) {
            _.$Message.error(err.toString());
            _.doListAllUsers();
          });
      } else {
        try {
          $userApi.sync.updateUserName(
            this.userManage.userEdit.userID,
            this.userManage.userEdit.UserName
          );
          $userApi.sync.updateUserPwd(
            this.userManage.userEdit.userID,
            this.userManage.userEdit.UserPWD
          );
          this.userManage.userEdit.show = false;
          this.$Message.info("操作成功");
        } catch (err) {
          this.$Message.error(err.toString());
        }
        this.doListAllUsers();
      }
    },
    // 加载权限
    doListFPermissions: function () {
      let _ = this;
      _.fPermissionManage.datas = [];
      _.fPermissionManage.loading = true;
      _.fPermissionManage.currentPageIndex = 1;
      $filepms
        .listFPermissions()
        .then(function (datas) {
          if (datas) {
            datas = JSON.parse(datas);
          }
          _.$set(_.fPermissionManage, "datas", datas);
          _.fPermissionManage.loading = false;
        })
        .catch(function (err) {
          _.$Message.error(err.toString());
          _.fPermissionManage.loading = false;
        });
    },
    // 添加权限
    onAddFPermission: function () {
      // 后期优化
      let _ = this;
      $userApi
        .listAllUsers()
        .then(function (datas) {
          if (datas) {
            _.fPermissionManage.userselector.datas = JSON.parse(datas);
          }
        })
        .catch(console.error);

      this.fPermissionManage.fPermissionEdit.isAdd = true;
      this.fPermissionManage.fPermissionEdit.userID = "";
      this.fPermissionManage.fPermissionEdit.path = "";
      this.fPermissionManage.fPermissionEdit.Permission = 0;
      this.fPermissionManage.fPermissionEdit.Permissions = [];
      this.fPermissionManage.fPermissionEdit.PermissionID = "";
      this.fPermissionManage.fPermissionEdit.show = true;
    },
    //
    onBeforFpmsSelectedDir: function (rows) {
      if (rows && rows.length > 0) {
        this.fPermissionManage.fPermissionEdit.path = rows[0].path;
      }
    },
    // 修改权限
    onEditFPermission: function () {
      let rows = this.$refs.fPermissionManage.getSelection();
      if (rows && rows.length == 1) {
        this.fPermissionManage.fPermissionEdit.isAdd = false;
        this.fPermissionManage.fPermissionEdit.userID = rows[0].userID;
        this.fPermissionManage.fPermissionEdit.path = rows[0].path;
        this.fPermissionManage.fPermissionEdit.Permission = parseInt(
          rows[0].Permission
        );
        this.fPermissionManage.fPermissionEdit.Permissions = (function (_) {
          let permission = _.fPermissionManage.fPermissionEdit.Permission;
          let res = [];
          if (permission) {
            for (let i = 0; i <= _.fPermissionManage.permissionMax; i++) {
              if (1 << i == (permission & (1 << i))) {
                res.push("" + i);
              }
            }
          }
          return res;
        })(this);
        this.fPermissionManage.fPermissionEdit.PermissionID =
          rows[0].PermissionID;
        this.fPermissionManage.fPermissionEdit.show = true;
      }
    },
    // 删除权限
    onDelFPermission: function () {
      let rows = this.$refs.fPermissionManage.getSelection();
      if (rows) {
        try {
          for (let i = 0; i < rows.length; i++) {
            $filepms.sync.delFPermission(rows[i].PermissionID);
          }
          this.$Message.info("操作成功");
        } catch (err) {
          this.$Message.error(err.toString());
        }
        this.doListFPermissions();
      }
    },
    // 保存
    onSaveFPermission: function () {
      if (!this.fPermissionManage.fPermissionEdit.userID) {
        this.$Message.error("用户ID不能为空");
        return;
      }
      if (!this.fPermissionManage.fPermissionEdit.path) {
        this.$Message.error("授权路径名不能为空");
        return;
      }
      let _ = this;
      this.fPermissionManage.fPermissionEdit.Permission = (function () {
        let res = 0;
        if (_.fPermissionManage.fPermissionEdit.Permissions.length > 0) {
          for (
            let i = 0;
            i < _.fPermissionManage.fPermissionEdit.Permissions.length;
            i++
          ) {
            res += 1 << _.fPermissionManage.fPermissionEdit.Permissions[i];
          }
        }
        return res;
      })();
      if (this.fPermissionManage.fPermissionEdit.isAdd) {
        $filepms
          .addFPermission(
            this.fPermissionManage.fPermissionEdit.userID,
            this.fPermissionManage.fPermissionEdit.path,
            this.fPermissionManage.fPermissionEdit.Permission
          )
          .then(function (data) {
            _.$Message.info("操作成功");
            _.doListFPermissions();
            _.fPermissionManage.fPermissionEdit.show = false;
          })
          .catch(function (err) {
            _.$Message.error(err.toString());
            _.doListFPermissions();
          });
      } else {
        $filepms
          .updateFPermission(
            this.fPermissionManage.fPermissionEdit.PermissionID,
            this.fPermissionManage.fPermissionEdit.Permission
          )
          .then(function (data) {
            _.$Message.info("操作成功");
            _.doListFPermissions();
            _.fPermissionManage.fPermissionEdit.show = false;
          })
          .catch(function (err) {
            _.$Message.error(err.toString());
            _.doListFPermissions();
          });
      }
    },
  },
  created: function () {
    this.doBuildAccountInfo(this.$root.$data.currentAccount);
    this.$set(
      this,
      "currentAccount",
      $utils.extendAttrs(this.currentAccount, this.$root.$data.currentAccount)
    );
  },
  watch: {
    currentTab: function (n, o) {
      if (n == "userManage") {
        this.doListAllUsers();
      } else if (n == "fPermissionManage") {
        this.doListFPermissions();
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
.sys-setting .usermanagetab .ivu-tabs-content {
  height: calc(100% - 33px - 16px);
}
</style>
