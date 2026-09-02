<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { apiClient } from '@/api/client'
import { NCard, NForm, NFormItem, NInput, NButton, NAlert } from 'naive-ui'

const router = useRouter()

const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const passwordConfirm = ref('')
const errorMessage = ref('')

const passwordMismatch = computed(() => {
    return passwordConfirm.value !== '' && password.value !== passwordConfirm.value
})

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
        router.push({ name: 'login' })
    } catch (error: any) {
        errorMessage.value = error.response?.data?.error ?? 'Registration failed.'
    }
}
</script>

<template>
  <NCard title="Register" style="max-width: 400px; margin: 40px auto;">
    <form @submit.prevent="handleRegister">
      <NFormItem label="First Name">
        <NInput v-model:value="firstName" />
      </NFormItem>
      <NFormItem label="Last Name">
        <NInput v-model:value="lastName" />
      </NFormItem>
      <NFormItem label="Email">
        <NInput v-model:value="email" />
      </NFormItem>
      <NFormItem label="Password">
        <NInput v-model:value="password" type="password" show-password-on="click" />
      </NFormItem>
      <NFormItem label="Confirm Password">
        <NInput v-model:value="passwordConfirm" type="password" show-password-on="click" />
      </NFormItem>
      <NAlert v-if="passwordMismatch" type="warning" style="margin-bottom: 16px;">Passwords do not match.</NAlert>
      <NAlert v-if="errorMessage" type="error" style="margin-bottom: 16px;">{{ errorMessage }}</NAlert>
      <NButton type="primary" attr-type="submit" block>Register</NButton>
    </form>
    <p style="margin-top: 16px;">
      Already have an account?
      <NButton text @click="router.push({ name: 'login' })">Log in</NButton>
    </p>
  </NCard>
</template>
