<template>
    <div>
      <div class="progress" id="progressBar"></div>
    </div>
  </template>
  
  <script>
  import $ from 'jquery';
  import { nextTick } from 'vue';
  
  export default {
    name: 'ProgressBar',
    props: {
      percentage: {
        type: Number,
        default: 0
      }
    },
    mounted() {
      // 动态加载 ProgressBarWars 插件
      this.loadScript('/js/ProgressBarWars.js', () => {
        nextTick(() => {
          this.initializeProgressBar();
        });
      });
    },
    watch: {
      percentage(newVal) {
        this.updateProgressBar(newVal);
      }
    },
    methods: {
      loadScript(src, callback) {
        const script = document.createElement('script');
        script.src = src;
        script.onload = callback;
        document.head.appendChild(script);
      },
      initializeProgressBar() {
        // 初始化进度条
        if (typeof $.fn.ProgressBarWars === 'function') {
          this.progressBar = $("#progressBar").ProgressBarWars({ porcentaje: this.percentage, estilo: "vader" });
        } else {
          console.error('ProgressBarWars is not loaded');
        }
      },
      updateProgressBar(percentage) {
        if (this.progressBar) {
          this.progressBar.mover(percentage);
        }
      }
    }
  }
  </script>
  
  <style scoped>
  .progress {
    margin-top: 30px;
  }
  </style>
  