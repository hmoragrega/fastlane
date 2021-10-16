<template>
  <v-card elevation="2">
    <v-card-title>{{ title }}</v-card-title>
    <v-card-subtitle>{{ web_url }}</v-card-subtitle>
    <v-card-text>{{ description }}</v-card-text>

    <v-card-actions>
      <v-container fluid class="pa-0">
        <v-btn class="ma-2" elevation="2" small outlined :href="web_url" target="blank">WEB</v-btn>
        <v-btn class="ma-2" elevation="2" v-if="!merge_enabled && !merged"  small outlined disabled>MERGE</v-btn>
        <v-btn class="ma-1" elevation="2" v-if="merge_enabled && !merged" small outlined color="success" v-on:click="merge(id)">MERGE</v-btn>
        <v-tooltip bottom v-for="user in approvals" v-bind:key="user.id">
          <template v-slot:activator="{ on, attrs }">
            <v-avatar class="ma-2" size="25" rounded v-bind:key="user.id" v-bind="attrs" v-on="on">
              <img :src=user.avatar_url :alt=user.name>
            </v-avatar>
          </template>
          <span>{{ user.name }}</span>
        </v-tooltip>
      </v-container>
    </v-card-actions>
   </v-card>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";

interface User {
  id: string;
  name: number;
}

const ReviewProps = Vue.extend({
  props: {
    title: String,
    id: String,
    web_url: String,
    can_be_merged: Boolean,
    merge_enabled: Boolean,
    merged: Boolean,
    approvals: Array,
    description: String,
  }
})

@Component
export default class Review extends ReviewProps {
  merge(reviewID: string): void {
    this.$store.dispatch('merge', reviewID)
  }
}
</script>
