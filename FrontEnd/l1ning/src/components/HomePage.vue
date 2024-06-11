<template>
  <div>
    <header>
      <h1 class="cool-font">L1N1NG Terminal</h1>
      <div class="user-info">
        <img src="@/assets/user-profile.jpg" alt="User Profile" class="profile-pic" @click="toggleDropdown" />
        <span class="username" @click="toggleDropdown">Luna</span>
        <div v-if="showDropdown" class="dropdown-container">
          <div class="dropdown-arrow"></div>
          <div class="dropdown-menu">
            <ul>
              <li><a href="#">Profile</a></li>
              <li><a href="#">Settings</a></li>
              <li><a href="#">Logout</a></li>
            </ul>
          </div>
        </div>
      </div>
    </header>
    <main>
      <div class="status">
        <h2 class="cool-font">SunnySport Status</h2>
        <EchartsPieChart :chartData="chartData" />
      </div>
      <div class="todo-list">
        <h2 class="cool-font">TodoList</h2>
      <TodoList />
        <div class="tabs">
          <!-- <button class="tab-btn">Today</button>
          <button class="tab-btn">This Week</button>
          <button class="tab-btn">This Month</button> -->
        </div>
        <ul class="tasks">
          <!-- <li><input type="checkbox" /> Task 1</li>
          <li><input type="checkbox" /> Task 2</li>
          <li><input type="checkbox" /> Task 3</li>
          <li><input type="checkbox" /> Task 4</li> -->
        </ul>
        <!-- <button class="submit-btn">Submit</button> -->
      </div>
    </main>
  </div>
</template>

<script>
import axios from 'axios';
import EchartsPieChart from './EchartsPieChart.vue';
import TodoList from './TodoList.vue';
export default {
  name: 'HomePage',
  components: {
    EchartsPieChart,TodoList
  },
  data() {
    return {
      chartData: {},
      showDropdown: false, // 控制下拉菜单显示状态
      isDone: false, // 控制任务完成状态
    };
  },
  created() {
    this.fetchChartData();
  },
  methods: {
    fetchChartData() {
      axios.get('http://localhost:8090/sunnysport')
        .then(response => {
          this.chartData = response.data;
        })
        .catch(error => {
          console.error("There was an error fetching the chart data!", error);
        });
    },
    toggleDropdown() {
      this.showDropdown = !this.showDropdown; // 切换下拉菜单显示状态
    },
    toggleDone() {
      this.isDone = !this.isDone; // 切换任务完成状态
    },
  }
};
</script>

<style scoped>
body {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
  background-color: #f4f4f4;
}

.cool-font {
  font-family: 'Pacifico', cursive;
}

#app {
  width: 80%;
  margin: 0 auto;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #333;
  color: white;
  padding: 10px 20px;
}

.user-info {
  position: relative; /* 相对定位以便放置下拉菜单 */
  display: flex;
  align-items: center;
  cursor: pointer;
}

.profile-pic {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
}

.config-btn {
  background-color: #fff;
  border: none;
  padding: 5px 10px;
  margin-right: 10px;
  cursor: pointer;
}

.dropdown-container {
  position: absolute;
  top: 60px;
  right: 0;
}

.dropdown-arrow {
  position: absolute;
  top: -10px;
  right: 10px;
  width: 0;
  height: 0;
  border-left: 10px solid transparent;
  border-right: 10px solid transparent;
  border-bottom: 10px solid white;
}

.dropdown-menu {
  background-color: white;
  border: 1px solid #ddd;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  border-radius: 5px;
  padding: 10px 0;
}

.dropdown-menu ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.dropdown-menu li {
  padding: 10px 20px;
}

.dropdown-menu li a {
  text-decoration: none;
  color: #333;
}

.dropdown-menu li a:hover {
  background-color: #f4f4f4;
}

main {
  display: flex;
  justify-content: space-between;
  padding: 20px;
}

.status,
.todo-list {
  background-color: #fffdea;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  width: 48%;
}

.tabs {
  display: flex;
  justify-content: space-between;
}

.tab-btn {
  flex: 1;
  padding: 10px;
  margin: 0 5px;
  background-color: #333;
  color: white;
  border: none;
  cursor: pointer;
}

.tasks {
  list-style: none;
  padding: 0;
  margin: 20px 0;
}

.tasks li {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  font-size: x-large;
}

.tasks li input[type="checkbox"] {
  width: 1.5em;
  height: 1.5em;
  margin-right: 10px;
}

.submit-btn {
  background-color: #333;
  color: white;
  border: none;
  padding: 10px 20px;
  cursor: pointer;
}

.todo {
  margin-left: 20px;
}

.todo:before {
  content: " ";
  display: inline-block;
  line-height: normal;
  font-size: 16px;
  border: 1px solid black;
  border-radius: 4px;
  background-position: 50%;
  white-space: normal;
  width: 16px;
  height: 16px;
  margin-right: 8px;
  margin-left: -24px;
  box-sizing: border-box !important;
  position: relative;
  cursor: pointer;
  pointer-events: auto;
  top: 2px;
}

.todo:hover:before {
  border: 1px solid rgb(100, 149, 237);
}

.done {
  text-decoration: line-through;
}

.done::before {
  background-image: url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTQiIGhlaWdodD0iMTQiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHBhdGggZD0iTTEwLjU4OSAzLjkwM2wuODA4LjgwOGEuMzUuMzUgMCAwMTAgLjQ5NUw2LjE4IDEwLjQyNWEuMzUuMzUgMCAwMS0uNDk1IDBMMi43MDMgNy40NDRhLjM1LjM1IDAgMDEwLS40OTVsLjgwOC0uODA4YS4zNS4zNSAwIDAxLjQ5NSAwbDEuOTI1IDEuOTI0IDQuMTYzLTQuMTYzYS4zNS4zNSAwIDAxLjQ5NSAweiIgZmlsbD0iIzMzNzBGRiIvPjwvc3ZnPg==);
}
</style>
