<template>
  <div>
    <header>
      <h1 class="cool-font">L1N1NG Terminal</h1>
    </header>
    <div class="content">
      <div class="login-container">
        <h1 class="cool-font">Login</h1>
        <form @submit.prevent="login">
          <div class="form-group">
            <div class="cool-font">
              <label for="username">Username:</label>
            </div>
            <input type="text" id="username" v-model="username" required placeholder="Please Input Your Username">
          </div>
          <div class="form-group">
            <div class="cool-font">
              <label for="password">Password:</label>
            </div>
            <div class="password-container">
              <input :type="passwordFieldType" id="password" v-model="password" required placeholder="Please Input Your Password">
              <button type="button" @click="togglePasswordVisibility">
                <i :class="passwordFieldType === 'password' ? 'fas fa-eye' : 'fas fa-eye-slash'"></i>
              </button>
            </div>
          </div>
          <div class="button-container">
            <button type="submit" class="login-btn">Login</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>


import axios from "axios";

export default {
  name: 'LoginPage',
  data() {
    return {
      username: '',
      password: '',
      passwordFieldType: 'password', // 控制密码输入框类型
      error: null,
    };
  },
  methods: {
    async login() {
      try {
        const response = await axios.post('/api/login', {
          username: this.username,
          password: this.password,
        });
        const token = response.data.token;
        localStorage.setItem('auth', token); // 保存 JWT 令牌到本地存储
        this.$router.push(this.$route.query.redirect || '/status'); // 重定向到原始请求路径或首页
      } catch (err) {
        this.error = 'Invalid credentials';
      }
    },
    togglePasswordVisibility() {
      // 切换密码输入框的类型
      this.passwordFieldType = this.passwordFieldType === 'password' ? 'text' : 'password';
    },
  },
};
</script>

<style scoped>
body {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
  background-color: #f4f4f4;
}

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

.content {
  background-color: #e0e0e0;
  min-height: calc(100vh - 60px); /* 60px accounts for the header height */
  padding-top: 20px;
  display: flex;
  justify-content: center;
}

.login-container {
  width: 600px; /* 控制登录框的宽度 */
  padding: 20px; /* 控制登录框内的填充 */
  background-color: #f5f5dc; /* 登录框的背景颜色 */
  border-radius: 36px; /* 圆角 */
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); /* 阴影效果 */
  margin-top: 20px; /* 控制登录框顶部的外边距 */
  margin-bottom: 500px; /* 控制登录框底部的外边距 */
  height: 400px;
}

h1 {
  margin-bottom: 10px; /* 减少标题下方的间距 */
  text-align: center;
}

.form-group {
  margin-bottom: 10px; /* 减少表单组间的间距 */
  text-align: left;
}

label {
  margin-bottom: .5em;
  color: #333;
  display: block;
  text-align: left;
}

input {
  border: 1px solid #CCCCCC;
  padding: 8px; /* 减少输入框的内边距 */
  width: 100%;
  box-sizing: border-box;
  border-radius: 36px;
  color: #666666;
}

.password-container {
  display: flex;
  align-items: center;
}

.password-container input {
  flex: 1;
  padding-right: 30px; /* 确保空间容纳按钮 */
}

.password-container button {
  margin-left: -30px; /* 将按钮移到输入框内 */
  background: none;
  border: none;
  cursor: pointer;
  color: #666666;
}

input::placeholder {
  font-size: 10px;
}

input[type="password"]::-ms-reveal,
input[type="password"]::-ms-clear {
  display: none;
}

#email {
  width: 410px;
  border: 1px solid #CCCCCC;
  padding: 8px; /* 减少输入框的内边距 */
  box-sizing: border-box;
  border-radius: 36px;
  color: #666666;
  text-align: left;
}

.button-container {
  display: flex;
  justify-content: center; /* 居中对齐 */
}

.login-btn {
  padding: 8px 16px; /* 减少按钮的内边距 */
  color: white;
  background-color: #333;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  width: 120px;
  height: 35px;
  box-shadow: 5px 5px 0 0 #666666;
  transition: all .1s ease-out;
  text-align: center;
  margin-top: 30px;
}

.login-btn:hover {
  background-color: #b3a0a0;
}

.login-btn:active {
  transform: translate(5px, 5px);
  box-shadow: 0px 0px 0 0 #f5f5dc;
}

.captcha {
  margin-left: 50px;
  padding: 8px 16px; /* 减少按钮的内边距 */
  color: white;
  background-color: #333;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  width: 110px;
  height: 26px;
  box-shadow: 5px 5px 0 0 #666666;
  transition: all .1s ease-out;
  font-size: 10px;
}

.captcha :hover {
  background-color: #555;
}

.captcha:hover:active {
  transform: translate(5px, 5px);
  box-shadow: 0px 0px 0 0 #f5f5dc;
}
</style>
