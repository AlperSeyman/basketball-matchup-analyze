<script setup lang="ts">
import { ref } from 'vue'
import { apiClient } from '@/api/client'

const email = ref('')
const message = ref('')
const errorMessage = ref('')

const handleForgotPassword = async () => {
    message.value = ''
    errorMessage.value = ''
    try {
        await apiClient.post('/auth/forgot-password', {email: email.value})
        message.value = 'If that email exists, a reset link has been sent.'
    } catch (error: any){
        errorMessage.value = error.response?.data?.error ?? 'Something went wrong.'
    }
}

</script>


<template>
    <h1>Forgot Password</h1>
    <form @submit.prevent="handleForgotPassword">
        <div>
            <label for="email">Email</label>
            <input type="email" id="email" v-model="email">
        </div>
        <p v-if="message">{{  message  }}</p>
        <p v-if="errorMessage">{{  errorMessage  }}</p>
        <button type="submit">Send Reset Link</button>
        <p>Remembered your password? <RouterLink :to="{name: 'login'}">Log in</RouterLink></p>
    </form>
</template>