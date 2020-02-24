<template>
  <v-app>
    <v-app-bar app color="#a20063" dark>
      <div>mb-trimmer</div>

      <v-spacer />

      <v-btn icon href="/twitter/sign-in">
        <v-icon>mdi-twitter</v-icon>
      </v-btn>
      <v-btn icon href="/twitter/sign-out">
        <v-icon>mdi-twitter</v-icon>
      </v-btn>
    </v-app-bar>

    <v-content>
      <Trimmer />
    </v-content>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import Trimmer from "./components/Trimmer.vue";

export default Vue.extend({
  name: "App",

  components: {
    Trimmer
  },

  data: () => ({
    //
  })
});
</script>
