<script setup lang="ts">

import { ref } from 'vue'
import { apiClient } from '@/api/client';

const currentPassword = ref('')
const newPassword = ref('')
const newPasswordConfirm = ref('')
const errorMessage = ref('')
const message = ref('')

const handleChangePassword = async () => {
    errorMessage.value = ''
    message.value = ''

    if (newPassword.value !== newPasswordConfirm.value){
        errorMessage.value = 'New passwords do not match'
        return
    }

    try{
        await apiClient.post('/auth/change-password', {
            current_password: currentPassword.value,
            new_password: newPassword.value,
        })
        message.value = 'Password change successfully.'
        currentPassword.value = ''
        newPassword.value = ''
        newPasswordConfirm.value = ''
    } catch (error: any) {
        errorMessage.value = error.response?.data?.error ?? 'Password change failed.'
    }
}
</script>


<template>
  <h1>My Profile</h1>
  <h2>Change Password</h2>
  <form @submit.prevent="handleChangePassword">
    <div>
      <label for="currentPassword">Current Password</label>
      <input id="currentPassword" type="password" v-model="currentPassword">
    </div>
    <div>
      <label for="newPassword">New Password</label>
      <input id="newPassword" type="password" v-model="newPassword">
    </div>
    <div>
      <label for="newPasswordConfirm">Confirm New Password</label>
      <input id="newPasswordConfirm" type="password" v-model="newPasswordConfirm">
    </div>
    <p v-if="message">{{ message }}</p>
    <p v-if="errorMessage">{{ errorMessage }}</p>
    <button type="submit">Change Password</button>
  </form>
</template>