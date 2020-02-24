<template>
  <div>
    <div class="buttons">
      <span v-for="item in metadata.items" :key="item.url">
        <v-btn small @click="playSound(item)">{{ item.text }}</v-btn>
      </span>
    </div>
    <div class="links">
      <span v-for="item in links" :key="item.url">
        <a :href="item.link">{{ item.text }}</a>
      </span>
    </div>
  </div>
</template>

<style>
.buttons .v-btn {
  margin: 4px;
  text-transform: unset;
}

.links a {
  margin: 4px;
  background: white;
}
</style>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";

interface Metadata {
  items: MetadataItem[];
}

interface MetadataItem {
  url: string;
  text: string;
}

@Component
export default class HomeComponent extends Vue {
  readonly metadataURL = "/metadata.json";

  readonly links = [
    { text: "Twitter", link: "https://twitter.com/g9v9g_mirei" },
    {
      text: "YouTube",
      link: "https://www.youtube.com/channel/UCeShTCVgZyq2lsBW9QwIJcw"
    },
    { text: "OpenREC", link: "https://www.openrec.tv/user/23_gundomirei" },
    {
      text: "Inspired by 宇志海ボタン",
      link: "http://ushiumi.ichiya-boshi.net/"
    },
    { text: "ロアボタン", link: "https://yuzukiroa.love/" },
    {
      text: "〇〇ボタン一覧",
      link:
        "https://wikiwiki.jp/nijisanji/%E2%97%8B%E2%97%8B%E3%83%9C%E3%82%BF%E3%83%B3"
    },
    { text: "要望などはこちら", link: "https://twitter.com/mirei_button" }
  ];

  metadata: Metadata = { items: [] };

  mounted() {
    fetch(this.metadataURL, {}).then(async response => {
      const json = await response.json();
      this.metadata = json as Metadata;
      this.metadata.items = this.metadata.items.sort((a, b) => {
        if (a.text > b.text) {
          return 1;
        } else if (a.text < b.text) {
          return -1;
        } else {
          return 0;
        }
      });
    });
  }

  playSound(item: MetadataItem) {
    new Audio(item.url).play();
  }
}
</script>
