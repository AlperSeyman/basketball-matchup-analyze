<script setup lang="ts">
import { ref } from 'vue'
import { NCard, NButton } from 'naive-ui'
import { useRouter } from 'vue-router'
import { divisions } from '@/data/teams'
const router = useRouter()

const failedLogos = ref(new Set<string>())

const handleLogoError = (abbr: string) => {
  failedLogos.value.add(abbr)
}


</script>

<template>
  <div style="max-width: 1000px; margin: 40px auto; padding: 0 16px;">
    <NCard title="Basketball Matchup Analyzer" style="margin-bottom: 24px;">
      <p>Welcome! Analyze basketball matchups using AI.</p>
    </NCard>

    <h2 style="margin-bottom: 16px;">All Teams</h2>
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 24px;">
      <div v-for="division in divisions" :key="division.name">
        <h3 style="text-transform: uppercase; margin-bottom: 8px;">{{ division.name }}</h3>
        <div v-for="team in division.teams" :key="team.abbr" style="display: flex; align-items: center; gap: 10px; padding: 6px 0;">
          <img
            v-if="!failedLogos.has(team.abbr)"
            :src="`https://a.espncdn.com/i/teamlogos/nba/500/${team.espnSlug || team.abbr.toLowerCase()}.png`"
            @error="handleLogoError(team.abbr)"
            style="width: 28px; height: 28px; object-fit: contain;"
          />
          <div
            v-else
            :style="{ backgroundColor: team.color, width: '28px', height: '28px', borderRadius: '50%' }"
          ></div>
          <span style="cursor: pointer;" @click="router.push({ name: 'team-detail', params: { id: team.espnId } })">{{ team.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
