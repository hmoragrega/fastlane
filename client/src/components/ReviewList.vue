<template>
  <v-container>
    <v-row v-for="x in merged" v-bind:key="`merged-${x.review.id}`">
      <v-col>
        <ReviewMerged
            v-bind:review="x.review"
            v-bind:pipeline="x.pipeline"
            v-bind:has_pipeline="x.has_pipeline"
        ></ReviewMerged>
      </v-col>
    </v-row>
    <v-row v-for="review in reviews" v-bind:key="review.id">
      <v-col>
        <Review
            v-bind:id="review.id"
            v-bind:title="review.title"
            v-bind:description="review.description"
            v-bind:web_url="review.web_url"
            v-bind:can_be_merged="review.can_be_merged"
            v-bind:approvals="review.approvals"
            v-bind:merge_enabled="review.merge_enabled"
        ></Review>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {mapState} from 'vuex'
import Review from "./Review.vue";
import ReviewMerged from "@/components/ReviewMerged.vue";

@Component({
  // Specify `components` option.
  // See Vue.js docs for all available options:
  // https://vuejs.org/v2/api/#Options-Data
  components: {
    Review,
    ReviewMerged
  },
  computed: {
    ...mapState({
      reviews: "reviews",
      merged: "merged",
    }),
  },
})
export default class ReviewList extends Vue {
  /*
  merge(reviewID: string): void {
    this.connection.send(JSON.stringify({
      event: "merge",
      reviewID: reviewID,
    }))
  }*/
}
</script>
