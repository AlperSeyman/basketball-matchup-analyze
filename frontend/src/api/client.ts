import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

export const apiClient = axios.create({ baseURL: '/api' })

apiClient.interceptors.request.use((config) => {
    const authStore = useAuthStore();
    if (authStore.accessToken) {
        config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }
    return config
})

apiClient.interceptors.response.use(
    (response) => {
        return response
    },
    async (error) => {
        const originalRequest = error.config as any
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true
            const authStore = useAuthStore()
            try {
                const refreshResponse = await axios.post('/api/auth/refresh')
                authStore.accessToken = refreshResponse.data.data.access_token
            } catch (refreshError) {
                authStore.accessToken = null
                window.location.href = '/login'
                return Promise.reject(refreshError)
            }
            return apiClient.request(error.config)
        }
        return Promise.reject(error)
    }
)