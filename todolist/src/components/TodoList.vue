<template>
  <div>
    <h1>Todo List</h1>
    <input v-model="newTodo" @keyup.enter="addTodo" placeholder="Add a new task" />
    <ul>
      <li v-for="(todo, index) in todos" :key="index">
        <span :class="{ completed: todo.completed }" @click="completeTodo(index)">
          {{ todo.text }}
        </span>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  data() {
    return {
      newTodo: '',
      todos: []
    }
  },
  methods: {
    addTodo() {
      if (this.newTodo.trim() !== '') {
        this.todos.push({ text: this.newTodo, completed: false });
        this.newTodo = '';
      }
    },
    completeTodo(index) {
      this.todos[index].completed = !this.todos[index].completed;
    }
  }
}
</script>

<style scoped>
ul {
  list-style-type: none;
  padding: 0;
}

li {
  margin: 5px 0;
  cursor: pointer;
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
