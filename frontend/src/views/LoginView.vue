<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { NCard, NForm, NFormItem, NInput, NButton, NAlert } from 'naive-ui'

const router = useRouter()
const email = ref('')
const password = ref('')
const errorMessage = ref('')
const authStore = useAuthStore()

const handleLogin = async () => {
    errorMessage.value = ''
    try {
        await authStore.login(email.value, password.value)
        router.push({ name: 'profile' })
    } catch (error: any) {
        errorMessage.value = error.response?.data?.error ?? 'Login failed.'
    }
}
</script>

<template>
  <NCard title="Login" style="max-width: 400px; margin: 40px auto;">
    <form @submit.prevent="handleLogin">
      <NFormItem label="Email">
        <NInput v-model:value="email" />
      </NFormItem>
      <NFormItem label="Password">
        <NInput v-model:value="password" type="password" show-password-on="click" />
      </NFormItem>
      <NAlert v-if="errorMessage" type="error" style="margin-bottom: 16px;">{{ errorMessage }}</NAlert>
      <NButton type="primary" attr-type="submit" block>Log In</NButton>
    </form>
    <p style="margin-top: 16px;">Don't have an account? <RouterLink :to="{ name: 'register' }">Register</RouterLink></p>
    <p><RouterLink :to="{ name: 'forgot-password' }">Forgot password?</RouterLink></p>
  </NCard>
</template>
