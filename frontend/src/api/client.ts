import axios from 'axios';

// vite proxy ayarı ile /api isteklerini 8080'e gönderiyoruz.
export const apiClient = axios.create({
    baseURL: '/api',
    headers: {
        'Content-Type': 'application/json',
    },
});

// auth token'ı her request'e auto eklemek için interceptor.
apiClient.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});