import axios from 'axios';

// 创建axios实例
const instance = axios.create({
    // baseURL: '/api', // 删除基础URL设置
});

// 添加请求拦截器
instance.interceptors.request.use(
    config => {
        // 在发送请求之前做些什么，比如添加token
        const token = localStorage.getItem('auth');
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    error => {
        // 对请求错误做些什么
        return Promise.reject(error);
    }
);

export default instance;
