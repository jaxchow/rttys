<template>
<div class="layout">
        <Layout>
            <Header>
                <Menu mode="horizontal" theme="dark" active-name="1">
                    <div class="layout-logo"></div>
                    <div class="layout-nav">
                      <Dropdown placement="bottom-start" @on-click="handleUserCommand">
                        <a href="javascript:void(0)">
                          <Avatar style="background-color: #87d068" icon="ios-person" />
                          {{ username }}
                            <Icon type="ios-arrow-down"></Icon>
                        </a>
                        <template #list>
                            <DropdownMenu>
                              <DropdownItem name="logout">{{ $t('Sign out') }}</DropdownItem>
                            </DropdownMenu>
                        </template>
                    </Dropdown>
                    </div>
                </Menu>
            </Header>
            <Layout :style="{padding: '0 50px'}">
                <Content :style="{padding: '24px 0', minHeight: '280px', background: '#fff'}">
                    <Layout>
                        <Sider hide-trigger :style="{background: '#fff'}">
                            <Menu theme="light" width="auto" :open-names="['1']">
                                <MenuItem name="1-1"><a href='/home'>设备管理</a></MenuItem>
                                <MenuItem name="1-2"><a href='/tenant'>租户管理</a></MenuItem>
                                <MenuItem name="1-3"><a href='/account'>用户管理</a></MenuItem>
                            </Menu>
                        </Sider>
                        <Content :style="{padding: '24px', minHeight: '280px', background: '#fff'}">
                            <router-view></router-view>
                        </Content>
                    </Layout>
                </Content>
            </Layout>
            <Footer class="layout-footer-center">2016-2023 &copy;telkits</Footer>
        </Layout>
    </div>
</template>
<style>
.layout{
    border: 1px solid #d7dde4;
    background: #f5f7f9;
    position: relative;
    border-radius: 4px;
    overflow: hidden;
}
.layout-logo{
    width: 100px;
    height: 30px;
    background: #5b6270;
    border-radius: 3px;
    float: left;
    position: relative;
    top: 15px;
    left: 20px;
}
.layout-nav{
    width: 120px;
    margin: 0 auto;
    margin-right: 20px;
}
.layout-footer-center{
    text-align: center;
}
</style>
<script>
   export default {
    name: 'Frame',
    data() {
      return {
        username: '',
        isadmin: false
      }
    },
    methods: {
      handleUserCommand(command) {
        if (command === 'logout') {
          this.axios.get('/signout').then(() => {
            this.$router.push('/login');
          });
        }
      }
    },
    mounted() {
      this.username = sessionStorage.getItem('rttys-username') || '';

      this.axios.get('/isadmin').then(res => {
        this.isadmin = res.data.admin;
      });
    }
  }
</script>