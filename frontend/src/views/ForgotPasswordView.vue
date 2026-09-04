<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiClient } from '@/api/client'
import { NCard, NForm, NFormItem, NInput, NButton, NAlert } from 'naive-ui'

const router = useRouter()
const email = ref('')
const message = ref('')
const errorMessage = ref('')

const handleForgotPassword = async () => {
    message.value = ''
    errorMessage.value = ''
    try {
        await apiClient.post('/auth/forgot-password', { email: email.value })
        message.value = 'If that email exists, a reset link has been sent.'
    } catch (error: any){
        if (error.response?.status === 400){
            errorMessage.value = 'Please check your information and try again.'
        } else {
            errorMessage.value = error.response?.data?.error ?? 'Something went wrong.'
        }
    }
}
</script>


<template>
  <NCard title="Forgot Password" style="max-width: 400px; margin: 40px auto;">
    <form @submit.prevent="handleForgotPassword">
        <NFormItem label="Email" required>
            <NInput v-model:value="email" />
        </NFormItem>
      <NAlert v-if="message" type="success" style="margin-bottom: 16px;">{{ message }}</NAlert>
      <NAlert v-if="errorMessage" type="error" style="margin-bottom: 16px;">{{ errorMessage }}</NAlert>
      <NButton type="primary" attr-type="submit" block>Send Reset Link</NButton>
    </form>
    <p style="margin-top: 16px;">
      Remembered your password?
      <NButton text @click="router.push({ name: 'login' })">Log in</NButton>
    </p>
  </NCard>
</template>