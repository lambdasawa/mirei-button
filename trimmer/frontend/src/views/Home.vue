<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png" />
    <Trimmer msg="Welcome to Your Vue.js App" />
  </div>
</template>

<script>
// @ is an alias to /src
import Trimmer from "@/components/Trimmer.vue";

export default {
  name: "Home",
  components: {
    Trimmer
  }
};
</script>
