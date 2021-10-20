import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

declare const window: any;

export default new Vuex.Store({
  state: {
    socket: {
      isConnected: false,
      message: '',
      reconnectError: false,
    },
    server_address: "",
    reviews: new Array<{id: string}>(),
    merged: [], // list of merged reviews
    notifications: new Array<{id: string}>(),
  },
  actions: {
    merge: function(context, reviewID) {
      Vue.prototype.$socket.sendObj({name: "MERGE", data: reviewID})
    },
    persist_server_address: function (context) {
      console.log("changing server address in store to ", context.state.server_address)
      localStorage.setItem("server_address", context.state.server_address)
      location.reload(); // hacky, but works
    }
  },
  modules: {},
  mutations: {
    close_notification(state, id) {
      state.notifications = state.notifications.filter(n => n.id !== id);
    },
    update_server_address(state, server_address) {
      state.server_address = server_address
    },
    SOCKET_ONOPEN (state, event)  {
      Vue.prototype.$socket = event.currentTarget
      state.socket.isConnected = true
    },
    SOCKET_ONCLOSE (state)  {
      state.socket.isConnected = false
    },
    SOCKET_ONERROR (state, event)  {
      console.error(state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message)  {
      switch (message.name) {
        case "REVIEWS":
          state.reviews = message.data;
          break;
        case "REVIEWS-MERGED":
          state.merged = message.data;
          break;
        case "NOTIFICATION":
          state.notifications.unshift(message.data);
          break;
        case "SYSTEM-NOTIFICATION":
          if (window.ipc !== undefined) {
            window.ipc.send('synchronous-message', message.data);
          }
          break;
      }
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT(state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },
  },
});
