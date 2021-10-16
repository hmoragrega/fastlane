<template>
  <v-container>
    <v-row v-for="notification in notifications" :key=notification.id>
      <v-col>
        <Notification
            v-bind:id="notification.id"
            v-bind:type="notification.type"
            v-bind:message="notification.message"
        ></Notification>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { mapState } from 'vuex'
import Notification from "./Notification.vue";

@Component({
  components: {
    Notification
  },
  computed: {
    ...mapState({
      notifications: "notifications"
    }),
  },
})
export default class Notifications extends Vue {
  /*
  merge(reviewID: string): void {
    this.connection.send(JSON.stringify({
      event: "merge",
      reviewID: reviewID,
    }))
  }*/
}
</script>
