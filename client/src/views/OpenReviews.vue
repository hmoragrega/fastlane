<template>
  <div>
    <notifications />
    <review-list />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import ReviewList from "../components/ReviewList.vue";
import Notifications from "../components/Notifications.vue";

export default Vue.extend({
  name: "OpenReviews",

  components: {
    ReviewList,
    Notifications,
  },
});
</script>
