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
    reviews: new Array<{id: string}>(),
    merged: [], // list of merged reviews
    notifications: new Array<{id: string}>(),
  },
  actions: {
    merge: function(context, reviewID) {
      Vue.prototype.$socket.sendObj({name: "MERGE", data: reviewID})
    }
  },
  modules: {},
  mutations: {
    close_notification(state, id) {
      console.log("deleting notification with ID", id)
      console.log("current count", state.notifications.length)
      state.notifications = state.notifications.filter(n => n.id !== id);
      console.log("count after deletion", state.notifications.length)
    },
    SOCKET_ONOPEN (state, event)  {
      Vue.prototype.$socket = event.currentTarget
      state.socket.isConnected = true
    },
    SOCKET_ONCLOSE (state, event)  {
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
          console.log("reviews merged", message.data)
          state.merged = message.data;
          break;
        case "NOTIFICATION":
          console.log("received notification", message.data)
          state.notifications.unshift(message.data);
          break;
        case "SYSTEM-NOTIFICATION":
          console.log("received system notification", message.data)
          window.ipc.send('synchronous-message', message.data);
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
