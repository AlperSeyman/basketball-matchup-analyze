<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiClient } from '@/api/client'

const router = useRouter()

const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const passwordConfirm = ref('')
const errorMessage = ref('')

const handleRegister = async () => {
    errorMessage.value = ''

    if (password.value !== passwordConfirm.value) {
        errorMessage.value = 'Passwords do not match.'
        return
    }

    try {
        await apiClient.post('/auth/register', {
            first_name: firstName.value,
            last_name: lastName.value,
            email: email.value,
            password: password.value,
        })
        router.push({ name: 'login'})
    } catch (error: any) {
        errorMessage.value = error.response?.data?.error ?? 'Registration failed.'
    }
}
</script>


<template>
    <h1>Register</h1>
    <form @submit.prevent="handleRegister">
        <div>
            <label for="firstName">First Name</label>
            <input type="text" id="firstName" v-model="firstName">
        </div>
        <div>
            <label for="lastName">Last Name</label>
            <input type="text" id="lastName" v-model="lastName">
        </div>
        <div>
            <label for="email">Email</label>
            <input type="email" id="email" v-model="email">
        </div>
        <div>
            <label for="password">Password</label>
            <input type="password" id="password" v-model="password">
        </div>
        <div>
            <label for="passwordConfirm">Confirm Password</label>
            <input type="password" id="passwordConfirm" v-model="passwordConfirm">
        </div>
        <p v-if="errorMessage">{{ errorMessage }}</p>
        <button type="submit">Register</button>
        <p>Already have an account? <RouterLink :to="{name: 'login'}">Log in</RouterLink></p>
    </form>
</template>