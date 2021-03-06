<template>
  <v-container id="trimmer">
    <div class="text-center">
      <v-icon>mdi-twitter</v-icon>
      <span>@{{ twitterStatus.status }}</span>
    </div>
    <div>
      <v-text-field type="text" v-model="rawVideoUrl" />
      <v-btn @click="setVideoUrl">読み込み</v-btn>
    </div>
    <v-row class="text-center">
      <v-col cols="12" lg="8" xs="12">
        <dir id="player-wrapper">
          <video id="player" ref="player" :src="videoUrl" controls></video>
        </dir>
      </v-col>
      <v-col cols="12" lg="4" mxs="12">
        <div class="mb-4">
          <div>
            <v-btn @click="trim">このへんを切り抜き</v-btn>
          </div>
          <div>
            <v-btn class="ma-2" @click="preview">プレビュー</v-btn>
          </div>
        </div>
        <div class="ma-4">範囲 : {{ getSpan() }}</div>
        <div>
          <span>開始位置 : {{ startPosition }}</span>
        </div>
        <div class="mb-4">
          <div>
            <v-btn class="ma-1" @click="updateStartPos(-5000)">-5s</v-btn>
            <v-btn class="ma-1" @click="updateStartPos(-1000)">-1s</v-btn>
            <v-btn class="ma-1" @click="updateStartPos(-100)">-0.1s</v-btn>
          </div>
          <div>
            <v-btn class="ma-1" @click="updateStartPos(5000)">+5s</v-btn>
            <v-btn class="ma-1" @click="updateStartPos(1000)">+1s</v-btn>
            <v-btn class="ma-1" @click="updateStartPos(100)">+0.1s</v-btn>
          </div>
        </div>
        <div>
          <span>終了位置 : {{ endPosition }}</span>
        </div>
        <div class="mb-4">
          <div>
            <v-btn class="ma-1" @click="updateEndPos(-5000)">-5s</v-btn>
            <v-btn class="ma-1" @click="updateEndPos(-1000)">-1s</v-btn>
            <v-btn class="ma-1" @click="updateEndPos(-100)">-0.1s</v-btn>
          </div>
          <div>
            <v-btn class="ma-1" @click="updateEndPos(5000)">+5s</v-btn>
            <v-btn class="ma-1" @click="updateEndPos(1000)">+1s</v-btn>
            <v-btn class="ma-1" @click="updateEndPos(100)">+0.1s</v-btn>
          </div>
        </div>
        <div class="mb-4">
          <div>
            <v-text-field type="text" v-model="text" placeholder="ボタンテキスト..." />
            <v-btn class="ma-2" @click="save">保存</v-btn>
          </div>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<style>
#trimmer .v-btn {
  text-transform: unset;
}

#player {
  width: 100%;
}
</style>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Trimmer extends Vue {
  twitterStatus: { status: string } = { status: "" };

  rawVideoUrl = "https://www.youtube.com/watch?v=p5BzZNH2mkU";
  videoUrl = "";

  text = "";

  startPosition = 0;
  endPosition = 0;

  audio: HTMLAudioElement = new Audio();

  $refs!: {
    player: HTMLVideoElement;
  };

  async mounted() {
    this.setVideoUrl();

    await fetch("/twitter/status").then(async response => {
      const json = await response.json();

      this.twitterStatus = json;
    });
  }

  setVideoUrl() {
    const videoUrl = encodeURIComponent(this.rawVideoUrl);
    this.videoUrl = `/video?url=${videoUrl}`;
  }

  getAudioUrl() {
    const start = this.startPosition;

    const duration = this.endPosition - this.startPosition;

    const url = `/sound?url=${this.rawVideoUrl}&start-ms=${start}&duration-ms=${duration}`;

    return url;
  }

  getSpan() {
    const getDurationText = (n: number): string => {
      const hour = String(Math.floor(n / 1000 / 3600)).padStart(2, "0");
      const min = String(Math.floor(((n / 1000) % 3600) / 60)).padStart(2, "0");
      const sec = String(Math.floor((n / 1000) % 60)).padStart(2, "0");
      const milli = String(Math.floor(n % 1000)).padStart(3, "0");
      return `${hour}:${min}:${sec}.${milli}`;
    };

    const start = getDurationText(this.startPosition);
    const end = getDurationText(this.endPosition);

    const duraiton = (this.endPosition - this.startPosition) / 1000;

    return `${start} ~ ${end} (${duraiton}s)`;
  }

  trim() {
    const player = this.$refs.player;
    player.pause();

    const currentTime = player.currentTime;

    this.endPosition = Math.floor(currentTime * 1000);
    this.startPosition = Math.floor(currentTime * 1000 - 5 * 1000);

    if (this.startPosition < 0) this.startPosition = 0;

    if (this.endPosition > player.duration * 1000)
      this.endPosition = player.duration * 1000;
  }

  preview() {
    this.audio.pause();

    const url = this.getAudioUrl();

    this.audio = new Audio(url);
    this.audio.play();
  }

  async save() {
    const fetchURL = "/register";
    const fetchOption: RequestInit = {
      method: "POST",
      cache: "no-cache",
      credentials: "same-origin", // include, same-origin, *omit
      headers: {
        "Content-Type": "application/json; charset=utf-8"
      },
      body: JSON.stringify({
        url: this.rawVideoUrl,
        startMS: this.startPosition,
        durationMS: this.endPosition - this.startPosition,
        text: this.text,
        tags: ""
      })
    };

    await fetch(fetchURL, fetchOption);
  }

  updateStartPos(x: number) {
    this.startPosition += x;
  }

  updateEndPos(x: number) {
    this.endPosition += x;
  }
}
</script>
