<script setup lang="ts">

import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiClient } from '@/api/client'


const route = useRoute()
const router = useRouter()

const newPassword = ref('')
const passwordConfirm = ref('')
const errorMessage = ref('')

const handleResetPassword = async () =>{
    errorMessage.value = ''

    if (newPassword.value !== passwordConfirm.value){
        errorMessage.value = 'Passwords do not match'
        return
    }

    try{
        await apiClient.post('/auth/reset-password', {
            token: route.query.token,
            new_password: newPassword.value
        })
        router.push({name: 'login'})
    } catch (error: any){
        errorMessage.value = error.response?.data?.error ?? 'Reset failed'
    }
}
</script>


<template>
  <h1>Reset Password</h1>
  <form @submit.prevent="handleResetPassword">
    <div>
      <label for="newPassword">New Password</label>
      <input id="newPassword" type="password" v-model="newPassword">
    </div>
    <div>
      <label for="passwordConfirm">Confirm New Password</label>
      <input id="passwordConfirm" type="password" v-model="passwordConfirm">
    </div>
    <p v-if="errorMessage">{{ errorMessage }}</p>
    <button type="submit">Reset Password</button>
  </form>
</template>