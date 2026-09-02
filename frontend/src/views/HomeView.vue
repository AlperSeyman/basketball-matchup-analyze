<script setup lang="ts">
import { NButton, NCard } from 'naive-ui'
</script>

<template>
  <NCard title="Basketball Matchup Analyzer" style="max-width: 500px; margin: 40px auto;">
    <p>Welcome! Analyze basketball matchups using AI.</p>
    <NButton tag="a" type="primary" @click="$router.push({ name: 'login' })">Log in</NButton>
    <NButton tag="a" @click="$router.push({ name: 'register' })">Register</NButton>
  </NCard>
</template>
