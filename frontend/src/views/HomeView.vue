<script setup lang="ts">
import { ref } from 'vue'
import { NCard, NButton } from 'naive-ui'
import { useRouter } from 'vue-router'

const router = useRouter()

const failedLogos = ref(new Set<string>())

const handleLogoError = (abbr: string) => {
  failedLogos.value.add(abbr)
}

const divisions = [
  {
    name: 'Atlantic',
    teams: [
      { name: 'Boston Celtics', abbr: 'BOS', color: '#007A33' },
      { name: 'Brooklyn Nets', abbr: 'BKN', color: '#000000' },
      { name: 'New York Knicks', abbr: 'NYK', color: '#006BB6' },
      { name: 'Philadelphia 76ers', abbr: 'PHI', color: '#006BB6' },
      { name: 'Toronto Raptors', abbr: 'TOR', color: '#CE1141' },
    ],
  },
  {
    name: 'Central',
    teams: [
      { name: 'Chicago Bulls', abbr: 'CHI', color: '#CE1141' },
      { name: 'Cleveland Cavaliers', abbr: 'CLE', color: '#860038' },
      { name: 'Detroit Pistons', abbr: 'DET', color: '#C8102E' },
      { name: 'Indiana Pacers', abbr: 'IND', color: '#002D62' },
      { name: 'Milwaukee Bucks', abbr: 'MIL', color: '#00471B' },
    ],
  },
  {
    name: 'Southeast',
    teams: [
      { name: 'Atlanta Hawks', abbr: 'ATL', color: '#E03A3E' },
      { name: 'Charlotte Hornets', abbr: 'CHA', color: '#1D1160' },
      { name: 'Miami Heat', abbr: 'MIA', color: '#98002E' },
      { name: 'Orlando Magic', abbr: 'ORL', color: '#0077C0' },
      { name: 'Washington Wizards', abbr: 'WAS', color: '#002B5C' },
    ],
  },
  {
    name: 'Northwest',
    teams: [
      { name: 'Denver Nuggets', abbr: 'DEN', color: '#0E2240' },
      { name: 'Minnesota Timberwolves', abbr: 'MIN', color: '#0C2340' },
      { name: 'Oklahoma City Thunder', abbr: 'OKC', color: '#007AC1' },
      { name: 'Portland Trail Blazers', abbr: 'POR', color: '#E03A3E' },
      { name: 'Utah Jazz', abbr: 'UTA', espnSlug: 'utah', color: '#002B5C' },
    ],
  },
  {
    name: 'Pacific',
    teams: [
      { name: 'Golden State Warriors', abbr: 'GSW', color: '#1D428A' },
      { name: 'LA Clippers', abbr: 'LAC', color: '#C8102E' },
      { name: 'Los Angeles Lakers', abbr: 'LAL', color: '#552583' },
      { name: 'Phoenix Suns', abbr: 'PHX', color: '#1D1160' },
      { name: 'Sacramento Kings', abbr: 'SAC', color: '#5A2D81' },
    ],
  },
  {
    name: 'Southwest',
    teams: [
      { name: 'Dallas Mavericks', abbr: 'DAL', color: '#00538C' },
      { name: 'Houston Rockets', abbr: 'HOU', color: '#CE1141' },
      { name: 'Memphis Grizzlies', abbr: 'MEM', color: '#5D76A9' },
      { name: 'New Orleans Pelicans', abbr: 'NOP', espnSlug: 'no', color: '#0C2340' },
      { name: 'San Antonio Spurs', abbr: 'SAS', color: '#C4CED4' },
    ],
  },
]
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
          <span>{{ team.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
