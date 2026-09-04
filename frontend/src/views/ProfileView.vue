<script setup lang="ts">

import { ref, computed } from 'vue'
import { apiClient } from '@/api/client'
import { NCard, NForm, NFormItem, NInput, NButton, NAlert } from 'naive-ui'

const currentPassword = ref('')
const newPassword = ref('')
const newPasswordConfirm = ref('')
const errorMessage = ref('')
const message = ref('')

const passwordMismatch = computed(() => {
  return newPasswordConfirm.value !== '' && newPassword.value !== newPasswordConfirm.value
})

const handleChangePassword = async () => {
  errorMessage.value = ''
  message.value = ''

  if (newPassword.value !== newPasswordConfirm.value) {
    errorMessage.value = 'New password do not match'
    return
  }

  try {
    await apiClient.post('/auth/change-password', {
      current_password: currentPassword.value,
      new_password: newPassword.value,
    })
    message.value = 'Password changed successfully.'
    currentPassword.value = ''
    newPassword.value = ''
    newPasswordConfirm.value = ''
  } catch (error: any){
    errorMessage.value = error.response?.data?.error ?? 'Password change failed.'
  }
}
</script>


<template>
  <NCard title="My Profile" style="max-width: 400px; margin: 40px auto;">
    <h3 style="margin-bottom: 16px;">Change Password</h3>
    <form @submit.prevent="handleChangePassword">
      <NFormItem label="Current Password">
        <NInput v-model:value="currentPassword" type="password" show-password-on="click" />
      </NFormItem>
      <NFormItem label="New Password">
        <NInput v-model:value="newPassword" type="password" show-password-on="click" />
      </NFormItem>
      <NFormItem label="Confirm New Password">
        <NInput v-model:value="newPasswordConfirm" type="password" show-password-on="click" />
      </NFormItem>
      <NAlert v-if="passwordMismatch" type="warning" style="margin-bottom: 16px;">Passwords do not match.</NAlert>
      <NAlert v-if="message" type="success" style="margin-bottom: 16px;">{{ message }}</NAlert>
      <NAlert v-if="errorMessage" type="error" style="margin-bottom: 16px;">{{ errorMessage }}</NAlert>
      <NButton type="primary" attr-type="submit" block>Change Password</NButton>
    </form>
  </NCard>
</template>