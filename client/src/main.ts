import Vue from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import vuetify from "./plugins/vuetify";
import './assets/css/styles.scss';
import VueNativeSock from "vue-native-websocket";

let serverAddress = localStorage.getItem("server_address");
console.log("got server address from local storage", serverAddress)

if (serverAddress === undefined || serverAddress === null || serverAddress.length == 0) {
    serverAddress = "127.0.0.1:3000";
    console.log("server address is empty, defaulting", serverAddress)
}

console.log("store server address before", store.state.server_address)
store.state.server_address = serverAddress
console.log("store server address now", store.state.server_address)

Vue.use(VueNativeSock, "ws://"+store.state.server_address+"/v1/ws", {
    store: store,
    format: 'json',
    reconnection: true, // (Boolean) whether to reconnect automatically (false)
    //reconnectionAttempts: 5, // (Number) number of reconnection attempts before giving up (Infinity),
    reconnectionDelay: 3000,
})

Vue.config.productionTip = false;

new Vue({
    router,
    store,
    vuetify,
    render: (h) => h(App),
}).$mount("#app");
