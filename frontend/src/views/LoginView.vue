<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const router = useRouter()
const email = ref('')
const password = ref('')
const errorMessage = ref('')
const authStore = useAuthStore()


const handleLogin = async () => {
    errorMessage.value = ''
    try {
        await authStore.login(email.value, password.value)
        router.push({name: 'profile'})
    } catch (error: any){
        errorMessage.value = error.response?.data?.error ?? 'Login failed.'
    }
}

</script>



<template>
    <h1>Login</h1>
    <form @submit.prevent="handleLogin">
        <div>
            <label for="email">Email</label>
            <input type="email" id="email" v-model="email">
        </div>
        <div>
            <label for="password">Password</label>
            <input type="password" id="password" v-model="password">
        </div>
        <p v-if="errorMessage">{{  errorMessage  }}</p>
        <button type="submit">Log In</button>
        <p>Don't have an account? <RouterLink :to="{ name: 'register' }">Register</RouterLink></p>
        <p><RouterLink :to="{ name: 'forgot-password' }">Forgot password?</RouterLink></p>
    </form>
</template>