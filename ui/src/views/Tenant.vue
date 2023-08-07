<template>
  <div style="padding:5px;">
    <Button style="margin-right: 4px;" type="primary" shape="circle" icon="ios-refresh" @click="handleRefresh" :disabled="loading">{{ $t('Refresh List') }}</Button>
    <Input style="margin-right: 4px;width:200px" v-model="filterString" search @input="handleSearch" :placeholder="$t('Please enter the filter key...')"/>
    <Button style="margin-right: 4px;" @click="showCmdForm" type="primary">创建租户</Button>
    <Button v-if="isadmin" style="margin-right: 4px;"  type="primary">{{ $t('users') }}</Button>
    <Table :loading="loading" :columns="columnsTenants" :data="userlists" style="margin-top: 10px; width: 100%" :no-data-text="$t('No Tenant')" @on-selection-change='handleSelection'>
      <template v-slot:action="{ row }">
        <Button type="warning" size="small" style="vertical-align: bottom;" @click="editTenant(row)">{{ $t('Edit') }}</Button>
        <Button type="warning" size="small" style="vertical-align: bottom;" @click="deleteTenant(row)">{{ $t('Delete') }}</Button>
      </template>
    </Table>
    <Modal v-model="cmdModal" title="创建租户" @on-ok="doCmd">
        <Form ref="cmdForm" :model="tenantData" :label-width="100" label-position="left">
          <FormItem :label="$t('name')" prop="name">
            <Input v-model="tenantData.name"/>
          </FormItem>
          <FormItem :label="$t('Owner')" prop="owner">
            <Input v-model.trim="tenantData.owner"/>
          </FormItem>
          <FormItem :label="$t('description')" prop="description">
            <Input v-model="tenantData.description" />
          </FormItem>
        </Form>
      </Modal>
  </div>
</template>

<script>
export default {
  name: 'tenant',
  data() {
    return {
      username: '',
      isadmin: false,
      filterString: '',
      loading: true,
      userlists: [],
      filteredUsers: [],
      selection: [],
      cmdModal:false,
      tenantData: {
        isModify:false,
        name: '',
        description: '',
        owner: ''
      },
      columnsTenants: [
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
          title: this.$t('name'),
          key: 'name',
          width: 200
        },
        {
          title: this.$t('owner'),
          key: 'owner',
          width: 200
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
    },
    showCmdForm() {
      this.cmdModal = true;
    },
    doCmd() {
      const data = this.tenantData
      // if (this.tenantData.isModify === true){
        this.axios.request({url:'/tenants', data, method:this.tenantData.isModify === true ? 'PUT' : 'POST'}).then((response) => {
            const resp = response.data;
            console.log(resp)
            this.getTenants()
        })
      // } else {
      //   this.axios.post('/tenants', data).then((response) => {
      //       const resp = response.data;
      //       this.getTenants()
      //   })
      // }
    },
    getTenants() {
      this.axios.get('/tenants').then(res => {
        this.loading = false;
        this.userlists = res.data.telant;
        this.selection = [];
        this.handleSearch();
      }).catch(() => {
       // this.$router.push('/login');
      });
    },
    handleRefresh() {
      this.loading = true;
      setTimeout(() => {
        this.getTenants();
      }, 500);
    },
    handleSelection(selection) {
      this.selection = selection;
    },
    editTenant(row) {
      this.tenantData = row;
      this.tenantData.isModify = true
      this.cmdModal = true;
    },
    deleteTenant(row) {
       this.axios.delete(`/tenants/${row.name}`, {
      }).then(() => {
        this.getTenants();
        this.$Message.success(this.$t('Delete success'));
      });
    }
  },
  mounted() {
    this.getTenants();
  }
}
</script>
