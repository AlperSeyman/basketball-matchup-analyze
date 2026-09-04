<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiClient } from '@/api/client'
import { NCard, NForm, NFormItem, NInput, NButton, NAlert } from 'naive-ui'

const route = useRoute()
const router = useRouter()

const newPassword = ref('')
const passwordConfirm = ref('')
const errorMessage = ref('')

const passwordMismatch = computed(() => {
    return passwordConfirm.value !== '' && newPassword.value !== passwordConfirm.value
})

const handleResetPassword = async () => {
    errorMessage.value = ''

    if (newPassword.value !== passwordConfirm.value) {
        errorMessage.value = 'Passwords do not match'
        return
    }

    try {
        await apiClient.post('/auth/reset-password', {
            token: route.query.token,
            new_password: newPassword.value,
        })
        router.push({ name: 'login' })
    } catch (error: any) {
        if (error.response?.status === 400) {
            errorMessage.value = 'Please check your information and try again.'
        } else {
            errorMessage.value = error.response?.data?.error ?? 'Reset failed'
        }
    }
}
</script>

<template>
  <NCard title="Reset Password" style="max-width: 400px; margin: 40px auto;">
    <form @submit.prevent="handleResetPassword">
      <NFormItem label="New Password" required>
        <NInput v-model:value="newPassword" type="password" show-password-on="click" />
      </NFormItem>
      <NFormItem label="Confirm New Password" required>
        <NInput v-model:value="passwordConfirm" type="password" show-password-on="click" />
      </NFormItem>
      <NAlert v-if="passwordMismatch" type="warning" style="margin-bottom: 16px;">Passwords do not match.</NAlert>
      <NAlert v-if="errorMessage" type="error" style="margin-bottom: 16px;">{{ errorMessage }}</NAlert>
      <NButton type="primary" attr-type="submit" block>Reset Password</NButton>
    </form>
  </NCard>
</template>
