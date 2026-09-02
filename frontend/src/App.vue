<script setup lang="ts">
import { computed } from 'vue'
import { darkTheme, NConfigProvider, NGlobalStyle } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'
import NavBar from '@/components/NavBar.vue'

const themePreference = useLocalStorage('theme-preference', 'system')
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches

const activeTheme = computed(() => {
  if (themePreference.value === 'dark') return darkTheme
  if (themePreference.value === 'light') return null
  return prefersDark ? darkTheme : null
})

const themeOverrides = {
  common: {
    primaryColor: '#FF6B35',
    primaryColorHover: '#FF8659',
    primaryColorPressed: '#E85A2A',
    borderRadius: '10px',
  },
}
</script>

<template>
  <NConfigProvider :theme="activeTheme" :theme-overrides="themeOverrides">
    <NGlobalStyle />
    <NavBar />
    <RouterView />
  </NConfigProvider>
</template>
