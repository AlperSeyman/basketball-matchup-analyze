import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'

export const useAuthStore = defineStore('auth', {
    state: () => {
        return {
            accessToken: null as string | null
        }
    },
    actions: {
        async login(email: string, password: string) {
            const response = await apiClient.post('/auth/login', { email, password })
            this.accessToken = response.data.data.access_token
        },
        async logout() {
            try {
                await apiClient.post('/auth/logout')
            } catch (error) {
                // ignore server errors
            }
            this.accessToken = null
        }
    },
});