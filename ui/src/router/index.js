import Vue from 'vue'
import VueRouter from 'vue-router'
import Frame from '../views/Frame.vue'
import Login from '../views/Login.vue'
import Home from '../views/Home.vue'
import Account from '../views/Account.vue'
import Tenant from '../views/Tenant.vue'
import Rtty from '../views/Rtty.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/',
    component:Frame,
    children:[
      {
        path: '/home',
        name: 'home',
        component: Home
      },
      {
        path: '/account',
        name: 'account',
        component: Account
      },
      {
        path: '/tenant',
        name: 'tenant',
        component: Tenant
      },
      {
        path: '/rtty/:devid',
        name: 'Rtty',
        component: Rtty,
        props: true
      }
    ]
  }
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  if (to.matched.length > 0 && to.matched[0].path === '/rtty/:devid') {
    const devid = to.params['devid'];
    Vue.axios.get(`/authorized/${devid}`).then(r => {
      if (r.data.authorized)
        next();
      else
        router.push('/login');
    });
    return;
  }

  if (to.path !== '/login') {
    Vue.axios.get('/alive').then(() => {
      next();
    }).catch(() => {
      router.push('/login');
    });
  } else {
    next();
  }
})

export default router
