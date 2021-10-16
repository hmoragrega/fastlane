<template>
  <v-alert
      :type=type
      dark
      dismissible
      outlined
      dense
      @input="closeNotification"
  >
    {{ message }}
  </v-alert>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";

const NotificationProps = Vue.extend({
  props: {
    id: String,
    type: String,
    message: String,
    index: Number,
  }
})

@Component({})
export default class Notifications extends NotificationProps {
  closeNotification() {
    this.$store.commit("close_notification", this.$props.id)
  }
}
</script>
