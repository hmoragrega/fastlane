<template>
  <v-card elevation="2">
    <v-card-title>{{ review.title }}</v-card-title>
    <v-card-subtitle>{{ review.web_url }}</v-card-subtitle>
    <v-card-text>
      {{ review.description }}
    </v-card-text>
    <v-card-actions>
      <v-container fluid class="pa-0">
        <v-btn class="ma-2" elevation="2" small outlined :href="review.web_url" target="blank">WEB</v-btn>
        <v-menu v-for="stage in pipeline.stages" v-bind:key="`jobs-${pipeline.id}-${stage.name}`"
            bottom
            offset-y
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn icon :color=stageColor(stage) v-bind:key="stage.name" v-bind="attrs" v-on="on">
              <v-icon>{{stageIcon(stage)}}</v-icon>
            </v-btn>
          </template>

          <v-list>
            <v-list-item
                v-for="(job, index) in stage.jobs"
                :key="index"
            >
              <v-btn icon :color=stageColor(job)>
                <v-icon>{{stageIcon(job)}}</v-icon>
              </v-btn>
              <v-list-item-title>{{ job.name }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-container>
    </v-card-actions>

   </v-card>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";

const ReviewMergedProps = Vue.extend({
  props: {
    review: Object,
    pipeline: Object,
    has_pipeline: Boolean,
  }
})

@Component
export default class ReviewMerged extends ReviewMergedProps {
  stageColor(stage) {
    switch (stage.status) {
      case "success": return "green";
      case "failed": return "red";
      case "running": return "yellow";
    }
    return "grey"
  }
  stageIcon(stage) {
    switch (stage.status) {
      case "success": return "mdi-check-circle-outline";
      case "failed": return "mdi-alert-circle-outline";
      case "running": return "mdi-play-circle-outline";
      case "skipped": return "mdi-skip-next-circle-outline";
      case "pending": return "mdi-pause-circle-outline";
      case "canceled": return "mdi-minus-circle-outline";
      case "manual": return "mdi-account-circle-outline";
    }
  }
}
</script>
