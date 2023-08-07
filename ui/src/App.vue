<template>
  <div id="app">
        <router-view></router-view>
    </div>
</template>

<script>
  export default {
    name: 'app',
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

<style>
  html, body, #app {
    margin:0;
    padding:0;
    width:100%;
    height:100%;
  }
</style>
