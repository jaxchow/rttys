<template>
    <Row :gutter="16">
        <Col span="6">
            <div style="text-align: right;font-weight: bold;">{{ $t('Username') }}:</div>
        </Col>
        <Col span="18">
            <div>{{ username }}</div>
        </Col>
        <Col span="6">
          <div style="text-align: right;font-weight: bold;">{{ $t('Admin') }}:</div>
        </Col>
        <Col span="18" >
            <div v-if="admin==0">普通用户</div>
            <div v-if="admin==1">管理员</div>
            <div v-if="admin==2">超级管理员</div>
        </Col>
        <Col span="6">
          <div style="text-align: right;font-weight: bold;">{{ $t('Token') }}:</div>
        </Col>
        <Col span="18">
            <div>{{ token }}</div>
        </Col>
        <Col span="6">
          <div style="text-align: right;font-weight: bold;">{{ $t('Tenant') }}:</div>
        </Col>
        <Col span="18">
            <div>{{ tenant }}</div>
        </Col>
    </Row>
</template>
<script>

export default {
  name: 'Home',
  data() {
    return {
      username: '',
      admin: '',
      tenant:'',
      token:''
    }
  },
  mounted() {
    this.username = sessionStorage.getItem('rttys-username') || '';
    this.token = sessionStorage.getItem('rttys-token') || '';
    this.tenant = sessionStorage.getItem('rttys-tenant') || '';
    this.admin = sessionStorage.getItem('rttys-admin') || '';
  }
}
</script>
