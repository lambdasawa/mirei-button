<template>
  <div class="home">
    <h1>郡道美玲ボタン</h1>
    <HomeComponent />
  </div>
</template>

<style>
body {
  background-image: url("../assets/background.jpg");
  background-repeat: repeat;
}

h1 {
  text-align: center;
  color: rgb(162, 0, 99);
  text-shadow: 1px 0 #000;
}
</style>

<script>
import HomeComponent from "@/components/HomeComponent";

export default {
  name: "Home",
  components: { HomeComponent }
};
</script>
