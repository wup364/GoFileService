<template>
  <div class="page">
    <div style="padding: 0px 5px 5px 5px">
      <div>
        <Button type="text" icon="md-refresh-circle" @click="loadList"
          >刷新</Button
        >
        <Button
          type="text"
          icon="md-add-circle"
          v-show="shows.add"
          @click="onAddFPermission"
          >新增</Button
        >
        <Button
          type="text"
          icon="ios-create"
          v-show="shows.edit"
          @click="onEditFPermission"
          >修改</Button
        >
        <Button
          type="text"
          icon="md-trash"
          v-show="shows.del"
          @click="onDelFPermission"
          >删除</Button
        >
      </div>
    </div>
    <Table
      ref="fpmsTb"
      :loading="loading"
      :columns="columns"
      :data="calcPageDatas"
      :highlight-row="true"
      :height="tableHeaght"
      v-watch-height="(ch, ph) => (tableHeaght = ph - 80)"
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
    <!-- 编辑权限 -->
    <Drawer
      :title="fPermissionEdit.isAdd ? '新增权限' : '修改权限'"
      width="450px"
      v-model="fPermissionEdit.show"
    >
      <Form
        ref="useredit"
        label-position="left"
        :model="fPermissionEdit"
        :label-width="80"
      >
        <FormItem class="mg_b_10" label="用户ID: " prop="userID">
          <Select
            v-if="fPermissionEdit.isAdd"
            v-model="fPermissionEdit.userID"
            transfer
          >
            <Option
              v-for="item in userselector.datas"
              :value="item.userID"
              :key="item.userID"
              >{{ item.userName }}</Option
            >
          </Select>
          <span v-else>{{ fPermissionEdit.userID }}</span>
        </FormItem>
        <FormItem class="mg_b_10" label="授权目录: " prop="path">
          <i-input
            v-if="fPermissionEdit.isAdd"
            v-model="fPermissionEdit.path"
            placeholder="授权路径"
            readonly
          >
            <Button slot="append" @click="onSelectDir('onBeforFpmsSelectedDir')"
              >选择</Button
            >
          </i-input>
          <span v-else>{{ fPermissionEdit.path }}</span>
        </FormItem>
        <FormItem class="mg_b_10" label="授权权限: " prop="Permission">
          <!-- <i-input v-model="fPermissionEdit.Permission" placeholder=""></i-input> -->
          <Checkbox-Group v-model="fPermissionEdit.permissions">
            <Checkbox label="1">Visible</Checkbox>
            <Checkbox label="2">Read</Checkbox>
            <Checkbox label="3">Write</Checkbox>
          </Checkbox-Group>
        </FormItem>
      </Form>
      <div>
        <Button class="fr mg_l_10" type="primary" @click="onSaveFPermission"
          >确定</Button
        >
        <Button class="fr" type="default" @click="fPermissionEdit.show = false"
          >关闭</Button
        >
      </div>
    </Drawer>
    <!-- 选择目录 -->
    <fileSelector
      :select-muti="fsSelector.selectMuti"
      :select-file="fsSelector.selectFile"
      :select-dir="fsSelector.selectDir"
      :start-path="fsSelector.startPath"
      :show-dailog="fsSelector.showDailog"
      @on-select="onSelectedFile"
      @on-cancel="onSelectCancel"
    ></fileSelector>
  </div>
</template>

<script>
import { $filepms } from "../../js/apis/filepermission";
import { $userApi } from "../../js/apis/user";
import fileSelector from "../file/fileselector.vue";
export default {
  name: "FilePmsList",
  components: {
    fileSelector,
  },
  data() {
    return {
      datas: [],
      loading: false,
      tableHeaght: 500,
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
        permission: 0,
        permissions: [],
        permissionID: "",
      },
      fsSelector: {
        operationObj: false,
        showDailog: false,
        selectFile: true,
        selectDir: true,
        selectMuti: true,
        startPath: "/",
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
          key: "permission",
          render(h, c) {
            return h(
              "span",
              {},
              $filepms.$TYPE.sum2Name(c.row.permission).join(", ")
            );
          },
        },
        {
          title: "创建时间",
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
      this.$refs.fpmsTb.selectAll(false);
      this.$refs.fpmsTb.toggleSelect(index);
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
    // 选择路径
    onSelectDir(cb) {
      this.fsSelector.operationObj = cb;
      this.fsSelector.selectFile = false;
      this.fsSelector.selectDir = true;
      this.fsSelector.selectMuti = false;
      this.fsSelector.startPath = "/";
      this.fsSelector.showDailog = true;
    },
    onSelectCancel() {
      this.fsSelector.showDailog = false;
    },
    onSelectedFile(rows) {
      this.fsSelector.showDailog = false;
      if (this.fsSelector.operationObj) {
        this[this.fsSelector.operationObj](rows);
      }
    },
    // 加载权限
    doListFPermissions() {
      this.datas = [];
      this.loading = true;
      this.currentPageIndex = 1;
      $filepms
        .listFPermissions()
        .then((datas) => {
          this.$set(this, "datas", datas);
          this.loading = false;
        })
        .catch((err) => {
          this.$Message.error(err.toString());
          this.loading = false;
        });
    },
    // 添加权限
    onAddFPermission() {
      // 后期优化
      $userApi
        .listAllUsers()
        .then((datas) => {
          if (datas) {
            this.userselector.datas = datas;
          }
        })
        .catch(console.error);

      this.fPermissionEdit.isAdd = true;
      this.fPermissionEdit.userID = "";
      this.fPermissionEdit.path = "";
      this.fPermissionEdit.permission = 0;
      this.fPermissionEdit.permissions = [];
      this.fPermissionEdit.permissionID = "";
      this.fPermissionEdit.show = true;
    },
    //
    onBeforFpmsSelectedDir(rows) {
      if (rows && rows.length > 0) {
        this.fPermissionEdit.path = rows[0].path;
      }
    },
    // 修改权限
    onEditFPermission() {
      let rows = this.$refs.fpmsTb.getSelection();
      if (rows && rows.length == 1) {
        this.fPermissionEdit.isAdd = false;
        this.fPermissionEdit.userID = rows[0].userID;
        this.fPermissionEdit.path = rows[0].path;
        this.fPermissionEdit.permission = parseInt(rows[0].permission);
        this.fPermissionEdit.permissions = ((_) => {
          let permission = this.fPermissionEdit.permission;
          let res = [];
          if (permission) {
            for (let i = 0; i <= this.permissionMax; i++) {
              if (1 << i == (permission & (1 << i))) {
                res.push("" + i);
              }
            }
          }
          return res;
        })(this);
        this.fPermissionEdit.permissionID = rows[0].permissionID;
        this.fPermissionEdit.show = true;
      }
    },
    // 删除权限
    onDelFPermission() {
      let rows = this.$refs.fpmsTb.getSelection();
      if (rows) {
        try {
          for (let i = 0; i < rows.length; i++) {
            $filepms.sync.delFPermission(rows[i].permissionID);
          }
          this.$Message.info("操作成功");
        } catch (err) {
          this.$Message.error(err.toString());
        }
        this.doListFPermissions();
      }
    },
    // 保存
    onSaveFPermission() {
      if (!this.fPermissionEdit.userID) {
        this.$Message.error("用户ID不能为空");
        return;
      }
      if (!this.fPermissionEdit.path) {
        this.$Message.error("授权路径名不能为空");
        return;
      }
      this.fPermissionEdit.permission = (() => {
        let res = 0;
        if (this.fPermissionEdit.permissions.length > 0) {
          for (let i = 0; i < this.fPermissionEdit.permissions.length; i++) {
            res += 1 << this.fPermissionEdit.permissions[i];
          }
        }
        return res;
      })();
      if (this.fPermissionEdit.isAdd) {
        $filepms
          .addFPermission(
            this.fPermissionEdit.userID,
            this.fPermissionEdit.path,
            this.fPermissionEdit.permission
          )
          .then((data) => {
            this.$Message.info("操作成功");
            this.doListFPermissions();
            this.fPermissionEdit.show = false;
          })
          .catch((err) => {
            this.$Message.error(err.toString());
            this.doListFPermissions();
          });
      } else {
        $filepms
          .updateFPermission(
            this.fPermissionEdit.permissionID,
            this.fPermissionEdit.permission
          )
          .then((data) => {
            this.$Message.info("操作成功");
            this.doListFPermissions();
            this.fPermissionEdit.show = false;
          })
          .catch((err) => {
            this.$Message.error(err.toString());
            this.doListFPermissions();
          });
      }
    },
    loadList() {
      this.doListFPermissions();
    },
  },
  created() {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
