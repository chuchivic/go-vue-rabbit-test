import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import VueNativeSock from 'vue-native-websocket'

Vue.config.productionTip = false
Vue.use(VueNativeSock, 'ws://localhost:3000', {
  store: store, 
  format: 'json',
  reconnection: true, // (Boolean) whether to reconnect automatically (false)
  reconnectionDelay: 3000, // (Number) how long to initially wait before attempting a new (1000)
})



new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
