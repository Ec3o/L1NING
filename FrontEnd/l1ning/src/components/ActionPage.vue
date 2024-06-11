<template>
  <div>
    <header>
      <h1 class="cool-font">L1N1NG Terminal</h1>
    </header>
    <div>
      <button @click="startKill">开始杀单词</button>
      <div v-if="loading">
        <p>正在处理...</p>
        <progress :value="progress" max="100">{{ progress }}%</progress>
      </div>
      <div v-if="result">
        <p>结果: {{ result }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref,  onBeforeUnmount } from 'vue';

export default {
  name: 'ActionPage',
  setup() {
    const loading = ref(false);
    const progress = ref(0);
    const result = ref(null);
    let socket = null;
    let heartbeatInterval = null;

    const startKill = () => {
      loading.value = true;
      progress.value = 0;
      result.value = null;

      setupWebSocket();
    };

    const setupWebSocket = () => {
      socket = new WebSocket('ws://localhost:8090/api/WordKiller');

      socket.onopen = () => {
        console.log('WebSocket connection opened');
        startHeartbeat();
      };

      socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === 'pong') {
          console.log('Received pong');
          return;
        }
        progress.value = data.progress;
        if (data.score) {
          result.value = data.score;
          loading.value = false;
          stopHeartbeat();
          socket.close();
        }
      };

      socket.onclose = (event) => {
        console.log('WebSocket connection closed:', event);
        stopHeartbeat();
        if (!event.wasClean) {
          console.log('Connection closed unexpectedly:', event.reason);
        }
      };

      socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        loading.value = false;
        stopHeartbeat();
      };
    };

    const startHeartbeat = () => {
      heartbeatInterval = setInterval(() => {
        if (socket && socket.readyState === WebSocket.OPEN) {
          socket.send(JSON.stringify({ type: 'heartbeat' }));
        }
      }, 10000); // 每10秒发送一次心跳
    };

    const stopHeartbeat = () => {
      if (heartbeatInterval) {
        clearInterval(heartbeatInterval);
        heartbeatInterval = null;
      }
    };

    onBeforeUnmount(() => {
      if (socket) {
        socket.close();
      }
      stopHeartbeat();
    });

    return {
      loading,
      progress,
      result,
      startKill
    };
  }
};
</script>

<style scoped>
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #333;
  color: white;
  padding: 10px 20px;
}
.cool-font {
  font-family: 'Pacifico', cursive;
}
</style>
