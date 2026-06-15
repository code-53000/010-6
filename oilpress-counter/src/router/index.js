import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/queue' },
  { path: '/queue', name: 'Queue', component: () => import('../views/Queue.vue'), meta: { title: '排队登记' } },
  { path: '/call', name: 'Call', component: () => import('../views/CallNumber.vue'), meta: { title: '叫号/榨油' } },
  { path: '/pickup', name: 'Pickup', component: () => import('../views/Pickup.vue'), meta: { title: '取油管理' } },
  { path: '/storage', name: 'Storage', component: () => import('../views/Storage.vue'), meta: { title: '油桶寄存' } },
  { path: '/pricing', name: 'Pricing', component: () => import('../views/Pricing.vue'), meta: { title: '价格配置' } }
]

export default createRouter({
  history: createWebHashHistory(),
  routes
})
