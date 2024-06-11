<template>
  <div>
    <input v-model="newTodo" @keyup.enter="addTodo" placeholder="Add a new task" />
    <ul>
      <li v-for="(todo, index) in paginatedTodos" :key="index">
        <TodoItem :todo="todo" @update:status="updateStatus(index, $event)" />
      </li>
    </ul>
    <div class="pagination">
      <button class="button" @click="prevPage" :disabled="currentPage === 1">
        <div class="button-box">
          <span class="button-elem">
            <svg viewBox="0 0 46 40" xmlns="http://www.w3.org/2000/svg">
              <path d="M46 20.038c0-.7-.3-1.5-.8-2.1l-16-17c-1.1-1-3.2-1.4-4.4-.3-1.2 1.1-1.2 3.3 0 4.4l11.3 11.9H3c-1.7 0-3 1.3-3 3s1.3 3 3 3h33.1l-11.3 11.9c-1 1-1.2 3.3 0 4.4 1.2 1.1 3.3.8 4.4-.3l16-17c.5-.5.8-1.1.8-1.9z"></path>
            </svg>
          </span>
        </div>
      </button>
      <span class="cool-font">Page {{ currentPage }} of {{ totalPages }}</span>
      <button class="button" @click="nextPage" :disabled="currentPage === totalPages">
        <div class="button-box">
          <span class="button-elem">
            <svg viewBox="0 0 46 40">
              <path d="M46 20.038c0-.7-.3-1.5-.8-2.1l-16-17c-1.1-1-3.2-1.4-4.4-.3-1.2 1.1-1.2 3.3 0 4.4l11.3 11.9H3c-1.7 0-3 1.3-3 3s1.3 3 3 3h33.1l-11.3 11.9c-1 1-1.2 3.3 0 4.4 1.2 1.1 3.3.8 4.4-.3l16-17c.5-.5.8-1.1.8-1.9z" transform="rotate(180 23 20)"></path>
            </svg>
          </span>
        </div>
      </button>
    </div>
  </div>
</template>

<script>
import TodoItem from './UI/TodoItem.vue';
import axios from "axios";

export default {
  components: {
    TodoItem
  },
  data() {
    return {
      newTodo: '',
      todos: [],
      currentPage: 1,
      pageSize: 12
    };
  },
  computed: {
    paginatedTodos() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.todos.slice(start, end);
    },
    totalPages() {
      return Math.ceil(this.todos.length / this.pageSize);
    }
  },
  methods: {
    async addTodo() {
      if (this.newTodo.trim() !== '') {
        try {
          const newTodo = {
            title: this.newTodo,
            status: false
          };
          const response = await axios.post('/api/todos', [newTodo]);

          if (response.data.status === "数据提交成功") {
            // 手动创建并添加新任务对象到 todos 数组中
            this.todos.push({
              id: this.todos.length ? this.todos[this.todos.length - 1].id + 1 : 1,
              ...newTodo
            });
            this.newTodo = '';

            if (this.todos.length % this.pageSize === 1) {
              this.currentPage = this.totalPages;
            }
          } else {
            console.error('Error adding todo: unexpected response', response);
          }
        } catch (error) {
          console.error('Error adding todo:', error);
        }
      }
    },
    async updateStatus(index, status) {
      const globalIndex = (this.currentPage - 1) * this.pageSize + index;
      const todo = this.todos[globalIndex];
      try {
        await axios.put(`/api/todos/${todo.id}`, {
          ...todo,
          status
        });
      } catch (error) {
        console.error('Error updating todo:', error);
      }
    },
    async loadTodos() {
      try {
        axios.defaults.headers.common['Authorization']=localStorage.getItem('auth')
        const response = await axios.get('/api/todos',);
        this.todos = response.data;
        this.currentPage = 1;
      } catch (error) {
        console.error('Error loading todos:', error);
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
      }
    },
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
      }
    }
  },
  mounted() {
    this.loadTodos();
  }
};
</script>

<style scoped>
ul {
  list-style-type: none;
  padding: 0;
}

.cool-font {
  font-family: 'Pacifico', cursive;
}

li {
  margin: 5px 0;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 20px;
}

.button {
  display: block;
  position: relative;
  width: 56px;
  height: 56px;
  margin: 0;
  overflow: hidden;
  outline: none;
  background-color: transparent;
  cursor: pointer;
  border: 0;
}

.button:before,
.button:after {
  content: "";
  position: absolute;
  border-radius: 50%;
  inset: 7px;
}

.button:before {
  border: 4px solid gray;
  transition: opacity 0.4s cubic-bezier(0.77, 0, 0.175, 1) 80ms,
  transform 0.5s cubic-bezier(0.455, 0.03, 0.515, 0.955) 80ms;
}

.button:after {
  border: 4px solid #000000;
  transform: scale(1.3);
  transition: opacity 0.4s cubic-bezier(0.165, 0.84, 0.44, 1),
  transform 0.5s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  opacity: 0;
}

.button:hover:before,
.button:focus:before {
  opacity: 0;
  transform: scale(0.7);
  transition: opacity 0.4s cubic-bezier(0.165, 0.84, 0.44, 1),
  transform 0.5s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.button:hover:after,
.button:focus:after {
  opacity: 1;
  transform: scale(1);
  transition: opacity 0.4s cubic-bezier(0.77, 0, 0.175, 1) 80ms,
  transform 0.5s cubic-bezier(0.455, 0.03, 0.515, 0.955) 80ms;
}

.button-box {
  display: flex;
  position: absolute;
  top: 0;
  left: 0;
}

.button-elem {
  display: block;
  width: 20px;
  height: 20px;
  margin: 17px 18px 0 18px;
  transform: rotate(180deg);
  fill: #000000;
}
</style>
