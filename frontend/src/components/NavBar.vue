<script setup lang="ts">
import { NDropdown, NButton, NSpace } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { computed } from 'vue'

const authStore = useAuthStore()
const router = useRouter()

const themePreference = useLocalStorage('theme-preference', 'system')

const themeLabel = computed(() => {
  return themePreference.value.charAt(0).toUpperCase() + themePreference.value.slice(1)
})

const themeOptions = [
  { label: 'Light', key: 'light' },
  { label: 'Dark', key: 'dark' },
  { label: 'System', key: 'system' }
]

const handleThemeSelect = (key: string) => {
  themePreference.value = key
}

const handleLogout = async () => {
  await authStore.logout()
  router.push({ name: 'login' })
}
</script>


<template>
  <div style="display: flex; justify-content: space-between; align-items: center; padding: 16px;">
    <NButton text style="font-size: 18px; font-weight: bold;" @click="router.push({ name: 'main' })">
      Basketball Matchup Analyzer
    </NButton>

    <NSpace align="center">
      <template v-if="authStore.accessToken">
        <NButton text @click="router.push({ name: 'profile' })">Profile</NButton>
        <NButton size="small" @click="handleLogout">Logout</NButton>
      </template>
      <template v-else>
        <NButton text @click="router.push({ name: 'login' })">Log in</NButton>
        <NButton type="primary" size="small" @click="router.push({ name: 'register' })">Register</NButton>
      </template>

      <NDropdown :options="themeOptions" @select="handleThemeSelect">
        <NButton quaternary size="small">{{ themeLabel }}</NButton>
      </NDropdown>
    </NSpace>
  </div>
</template>
