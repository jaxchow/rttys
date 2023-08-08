<template>
  <div style="padding:5px;">
    <Button style="margin-right: 4px;" type="primary" shape="circle" icon="ios-refresh" @click="handleRefresh" :disabled="loading">{{ $t('Refresh List') }}</Button>
    <Input style="margin-right: 4px;width:200px" v-model="filterString" search @input="handleSearch" :placeholder="$t('Please enter the filter key...')"/>
    <Button style="margin-right: 4px;" @click="showCmdForm" type="primary">添加用户</Button>
    <Table :loading="loading" :columns="columnsUsers" :data="filteredUser" style="margin-top: 10px; width: 100%" :no-data-text="$t('No acounnt')" @on-selection-change='handleSelection'>
      <template v-slot:admin="{ row }">
        <span v-if="row.admin ==0">普通用户</span>
        <span v-if="row.admin ==1">租户管理员</span>
        <span v-if="row.admin ==2">超级管理员</span>
      </template>
      <template v-slot:action="{ row }">
        <Button size="small" style="vertical-align: bottom;" @click="editUser(row)">{{ $t('Edit') }}</Button>
        <Button type="warning" size="small" style="vertical-align: bottom;" @click="deleteUsers(row)">{{ $t('Delete') }}</Button>
      </template>
    </Table>
    <Modal v-model="cmdModal" title="添加用户" @on-ok="doCmd">
      <Form ref="cmdForm" :model="cmdData" :label-width="100" label-position="left">
        <FormItem :label="$t('Role')" prop='admin' >
          <Select v-model="cmdData.admin" style="width:200px">
            <Option value=0 key=0>普通用户</Option>
            <Option value=1 key=1>租户管理员</Option>
          </Select>
        </FormItem>
        <FormItem :label="$t('Username')" prop="username">
          <Input v-model="cmdData.username"/>
        </FormItem>
        <FormItem :label="$t('Password')" prop="password">
          <Input v-model="cmdData.password" type="password" password />
        </FormItem>
        <FormItem :label="$t('Tenant')" prop="tenant">
          <Input v-model="cmdData.tenant" />
        </FormItem>
        <FormItem :label="$t('Description')" prop="description">
          <Input v-model.trim="cmdData.description"/>
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<script>
export default {
  name: 'User',
  data() {
    return {
      username: '',
      isadmin: 0,
      filterString: '',
      loading: true,
      userlists: [],
      filteredUsers: [],
      selection: [],
      cmdModal: false,
      cmdData: {
        username: '',
        isModify:false,
        description: '',
        admin:2,
        tenant: sessionStorage.getItem('rttys-tenant')
      },
      columnsUsers: [
        {
          title: '#',
          type: 'index',
          width: 100
        },
        {
          type: 'selection',
          width: 60
        },
        {
          title: this.$t('username'),
          key: 'username',
          width: 200
        },
        {
          title: this.$t('admin'),
          // key: 'admin' ,
          slot: 'admin',
          width: 200
        },
        {
          title: this.$t('tenant'),
           key: 'tenant',
          width: 200
        },
        {
          title: this.$t('Token'),
          key: 'token'
        },
        {
          title: this.$t('Description'),
          key: 'description',
          tooltip: true
        },
        {
          slot: 'action',
          width: 200
        }
      ]
    }
  },
  methods: {
    handleUserCommand(command) {
      if (command === 'logout') {
        this.axios.get('/signout').then(() => {
         // this.$router.push('/login');
        });
      }
    },
    handleSearch() {
      this.filteredUser = this.userlists.filter((d) => {
        const filterString = this.filterString.toLowerCase();
        return d.username.toLowerCase().indexOf(filterString) > -1;
      });
      // console.log('this.filteredUser', this.filteredUser)
    },
    editUser(row) {
      this.cmdData = row;
      this.cmdModal = true;
      this.cmdData.isModify = true
    },
    getUsers() {
      this.axios.get('/tenants/users').then(res => {
        this.loading = false;
        this.userlists = res.data.users;
        this.selection = [];
        this.handleSearch();
      }).catch(() => {
       // this.$router.push('/login');
      });
    },
    handleRefresh() {
      this.loading = true;
      setTimeout(() => {
        this.getUsers();
      }, 500);
    },
    handleSelection(selection) {
      this.selection = selection;
    },
    deleteUsers(row) {
      this.$Modal.confirm({
        content: '是否确认删除数据',
        onCancel: () => {},
        onOk: () => {
          this.axios.delete(`/tenants/users/${row.username}`, {
            // devices: offlines.map(s => s.id)
          }).then(() => {
            this.$Message.success(this.$t('Delete success'));
            this.getUsers();
          });
        }
       })
    },
    doCmd() {
      this.cmdData.admin = parseInt(this.cmdData.admin, 10)
      const data = this.cmdData
      if (this.cmdData.isModify === true) {
         this.axios.put('/tenants/users', data).then((response) => {
          const resp = response.data;
          console.log(resp)
          this.getUsers();
        })
      } else {
        this.axios.post('/tenants/users', data).then((response) => {
          const resp = response.data;
          console.log(resp)
          this.getUsers();
        })
      }
    },
    showCmdForm() {
      this.cmdModal = true;
    }
  },
  mounted() {
    this.isadmin = sessionStorage.getItem('rttys-admin')
    this.getUsers();
  }
}
</script>
