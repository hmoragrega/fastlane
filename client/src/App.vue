<template>
  <v-app>
    <v-progress-linear
        :active=!connected
        dark
        absolute
        top
        indeterminate
        color="yellow darken-2"
    ></v-progress-linear>

    <v-main>
      <router-view />
    </v-main>

    <template>
        <v-dialog
            v-model="dialog"
            fullscreen

            transition="dialog-bottom-transition"
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn
                fab
                dark
                small
                fixed
                top
                right
                color="primary"
                style="z-index: 100"
                v-bind="attrs"
                v-on="on"
            >
              <v-icon dark>
                mdi-cog-outline
              </v-icon>
            </v-btn>
          </template>
          <v-card>
            <v-toolbar
                dark
                color="primary"
            >
              <v-btn
                  icon
                  dark
                  @click="discard()"
              >
                <v-icon>mdi-close</v-icon>
              </v-btn>
              <v-toolbar-title>Settings</v-toolbar-title>
              <v-spacer></v-spacer>
              <v-toolbar-items>
                <v-btn
                    dark
                    text
                    @click=save()
                >
                  Save
                </v-btn>
              </v-toolbar-items>
            </v-toolbar>
            <v-list
                two-line
                subheader
            >
              <v-subheader>Configuration</v-subheader>
              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>Server address</v-list-item-title>
                  <v-list-item-subtitle>Hostname of the fastlane server</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
              <v-list-item>
                <v-list-item-content>
                  <v-text-field v-model="options_server_address" :rules=rules></v-text-field>
                </v-list-item-content>
              </v-list-item>
              <v-subheader>Notifications</v-subheader>
              <v-list-item>
                <v-list-item-action>
                  <v-checkbox v-model="notifications" disabled></v-checkbox>
                </v-list-item-action>
                <v-list-item-content>
                  <v-list-item-title>Review can be merged</v-list-item-title>
                  <v-list-item-subtitle>Notify me about reviews ready to be merged.</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-card>
        </v-dialog>
    </template>

    <v-overlay :value="!connected"></v-overlay>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import {mapMutations, mapState} from "vuex";

export default Vue.extend({
  name: "App",

  computed: {
    options_server_address: {
      get() {
        return this.$store.state.server_address
      },
      set(value) {
        this.update_server_address(value)
      }
    },
    ...mapState({
      connected: state => state.socket.isConnected,
      server_address: state => state.server_address,
    }),
  },

  methods: {
    save() {
      if (!isValidWSUrl(this.options_server_address || '')) {
        return
      }
      this.dialog = false;
      this.$store.dispatch("persist_server_address")
    },
    discard() {
      if (!isValidWSUrl(this.options_server_address || '')) {
        return
      }
      this.dialog = false;
      this.$store.dispatch("persist_server_address")
    }, ...mapMutations([
        "update_server_address",
    ])
  },

  data: () => ({
    dialog: false,
    notifications: true,
    rules: [
      value => !!value || 'Required.',
      value => isValidWSUrl(value || '') || 'Invalid hostname',
    ],
  }),
});

function isValidWSUrl(string) {
  let url;

  try {
    url = new URL("ws://"+string);
  } catch (_) {
    return false;
  }

  return true;
}

</script>
