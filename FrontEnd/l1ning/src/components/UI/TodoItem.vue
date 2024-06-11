<template>
  <div class="todo-item" @click="toggleCheck">
    <CheckCircle :checked="completed" />
    <span :class="{ completed: completed }">
      {{ todo.title }}
    </span>
  </div>
</template>

<script>
import CheckCircle from './CheckCircle.vue';

export default {
  components: {
    CheckCircle
  },
  props: {
    todo: {
      type: Object,
      required: true,
      validator(value) {
        return value.title && typeof value.title === 'string';
      }
    }
  },
  data() {
    return {
      completed: this.todo.status
    };
  },
  watch: {
    todo: {
      handler(newVal) {
        this.completed = newVal.status;
      },
      deep: true,
      immediate: true
    }
  },
  methods: {
    toggleCheck() {
      this.completed = !this.completed;
      this.$emit('update:status', this.completed);
    }
  }
};
</script>

<style scoped>
.todo-item {
  display: flex;
  align-items: center;
  cursor: pointer;
  font-size: x-large;
}

.todo-item > svg {
  margin-right: 10px; /* 调整 CheckCircle 和文字的间距 */
}

.completed {
  position: relative;
  display: inline-block;
  color: #aaa;
  overflow: hidden;
}

.completed::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 2px;
  background: black;
  top: 50%;
  left: 0;
  animation: strike-through 0.5s forwards;
}

@keyframes strike-through {
  0% {
    width: 0;
  }
  100% {
    width: 100%;
  }
}
</style>
