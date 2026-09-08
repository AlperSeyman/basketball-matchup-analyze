<script setup lang="ts">
import { useRoute } from 'vue-router'
import { ref, onMounted } from 'vue'
import { NCard, NSpin, NDataTable } from 'naive-ui'
import { getPlayerStats, getPlayerGamelog, getTeamRoster } from '@/api/espn'

const route = useRoute()
const player = ref<any>(null)
const ppg = ref('')
const rpg = ref('')
const apg = ref('')
const lastFiveGames = ref<any[]>([])

const gameColumns = [
  { title: 'Date', key: 'date' },
  { title: 'Opponent', key: 'opponent' },
  { title: 'Result', key: 'result' },
  { title: 'MIN', key: 'min' },
  { title: 'PTS', key: 'pts' },
  { title: 'REB', key: 'reb' },
  { title: 'AST', key: 'ast' },
]

onMounted(async () => {
  const playerId = String(route.params.id)
  const statsData = await getPlayerStats(playerId)

  if (statsData.categories) {
    const category = statsData.categories.find((c: any) => c.name === 'averages')
    const names = category.names
    const seasons = category.statistics
    const latestSeason = seasons[seasons.length - 1]

    ppg.value = latestSeason.stats[names.indexOf('avgPoints')]
    rpg.value = latestSeason.stats[names.indexOf('avgRebounds')]
    apg.value = latestSeason.stats[names.indexOf('avgAssists')]
  }

  const teamId = String(route.query.team)
  const rosterData = await getTeamRoster(teamId)
  player.value = rosterData.athletes.find((a: any) => a.id === playerId)

  const gamelogData = await getPlayerGamelog(playerId)

  if (gamelogData.seasonTypes && gamelogData.seasonTypes.length > 0) {
    const gameNames = gamelogData.names
    const firstSeasonType = gamelogData.seasonTypes[0]
    const allEvents: any[] = []

    for (const category of firstSeasonType.categories) {
      if (category.events) {
        for (const event of category.events) {
          allEvents.push(event)
        }
      }
    }

    lastFiveGames.value = allEvents.slice(0, 5).map((event: any) => {
      const gameInfo = gamelogData.events[event.eventId]
      return {
        date: gameInfo.gameDate.slice(0, 10),
        opponent: gameInfo.opponent.displayName,
        result: gameInfo.gameResult,
        min: event.stats[gameNames.indexOf('minutes')],
        pts: event.stats[gameNames.indexOf('points')],
        reb: event.stats[gameNames.indexOf('totalRebounds')],
        ast: event.stats[gameNames.indexOf('assists')],
      }
    })
  }
})
</script>

<template>
  <div v-if="player" style="max-width: 700px; margin: 40px auto;">
    <NCard :title="player.displayName">
      <div style="display: flex; gap: 16px; align-items: center;">
        <img :src="player.headshot?.href" style="width: 100px; height: 100px; border-radius: 50%;" />
        <div>
          <p>{{ player.position?.displayName }} | #{{ player.jersey }}</p>
          <p>{{ player.displayHeight }} | {{ player.displayWeight }} | {{ player.age }} years old</p>
          <p>College: {{ player.college?.name || 'N/A' }}</p>
        </div>
      </div>
      <div style="display: flex; gap: 32px; margin-top: 16px;">
        <div>
          <div>PPG</div>
          <div style="font-size: 20px; font-weight: bold;">{{ ppg }}</div>
        </div>
        <div>
          <div>RPG</div>
          <div style="font-size: 20px; font-weight: bold;">{{ rpg }}</div>
        </div>
        <div>
          <div>APG</div>
          <div style="font-size: 20px; font-weight: bold;">{{ apg }}</div>
        </div>
      </div>

      <h3 style="margin-top: 24px;">Last 5 Games</h3>
      <NDataTable :columns="gameColumns" :data="lastFiveGames" :bordered="false" />
    </NCard>
  </div>
  <div v-else style="text-align: center; margin-top: 80px;">
    <NSpin size="large" />
  </div>
</template>
