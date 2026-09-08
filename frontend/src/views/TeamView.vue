<script setup lang="ts">
import { ref, onMounted, watch, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NCard, NSpin, NDataTable, NSelect } from 'naive-ui'
import { divisions } from '@/data/teams'
import { getStandings, getTeamRoster, getTeamStatistics } from '@/api/espn'

const route = useRoute()
const router = useRouter()
const team = ref<any>(null)
const stats = ref<any[]>([])
const roster = ref<any[]>([])
const detailedStats = ref<any[]>([])
const teamColor = ref('')
const teamLogo = ref('')
const conferenceAbbr = ref('')

const columns = [
  { title: '#', key: 'jersey' },
  {
    title: 'Name',
    key: 'displayName',
    render: (row: any) => h('a', {
      style: 'cursor: pointer; color: inherit; text-decoration: underline;',
      onClick: () => router.push({ name: 'player-detail', params: { id: row.id }, query: { team: route.params.id } }),
    }, row.displayName),
  },
  { title: 'Pos', key: 'position', render: (row: any) => row.position?.abbreviation },
  { title: 'Height', key: 'displayHeight' },
  { title: 'Weight', key: 'displayWeight' },
  { title: 'Age', key: 'age' },
]

const teamOptions: any[] = []
for (const division of divisions) {
  for (const t of division.teams) {
    teamOptions.push({
      label: t.name,
      value: t.espnId,
      espnSlug: t.espnSlug,
      abbr: t.abbr,
    })
  }
}

const renderTeamLabel = (option: any) => {
  return h('div', { style: 'display: flex; align-items: center; gap: 8px;' }, [
    h('img', {
      src: `https://a.espncdn.com/i/teamlogos/nba/500/${option.espnSlug || option.abbr.toLowerCase()}.png`,
      style: 'width: 20px; height: 20px; object-fit: contain;',
    }),
    option.label,
  ])
}

const handleTeamSelect = (value: string) => {
  router.push({ name: 'team-detail', params: { id: value } })
}

const loadTeamData = async () => {
  team.value = null
  stats.value = []
  detailedStats.value = []
  roster.value = []

  const teamId = String(route.params.id)
  const data = await getStandings()

  const rosterData = await getTeamRoster(teamId)
  roster.value = rosterData.athletes

  const statsData = await getTeamStatistics(teamId)
  teamColor.value = statsData.team.color
  teamLogo.value = statsData.team.logo

  for (const category of statsData.results.stats.categories) {
    for (const stat of category.stats) {
      detailedStats.value.push(stat)
    }
  }

  for (const conference of data.children) {
    for (const entry of conference.standings.entries) {
      if (entry.team.id === teamId) {
        team.value = entry.team
        stats.value = entry.stats
        conferenceAbbr.value = conference.abbreviation
      }
    }
  }
}

onMounted(() => {
  loadTeamData()
})

watch(() => route.params.id, () => {
  loadTeamData()
})

const getStat = (statName: string) => {
  const stat = stats.value.find((s: any) => s.name === statName)
  return stat ? stat.displayValue : '...'
}

const getDetailedStat = (statName: string) => {
  const stat = detailedStats.value.find((s: any) => s.name === statName)
  return stat ? stat.displayValue : '...'
}

const ordinal = (n: number) => {
  if (n % 10 === 1 && n % 100 !== 11) return n + 'st'
  if (n % 10 === 2 && n % 100 !== 12) return n + 'nd'
  if (n % 10 === 3 && n % 100 !== 13) return n + 'rd'
  return n + 'th'
}
</script>

<template>
  <div v-if="team" style="max-width: 900px; margin: 40px auto;">
    <div
      :style="{
        backgroundColor: '#' + teamColor,
        padding: '24px',
        borderRadius: '8px',
        color: 'white',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        marginBottom: '24px',
        flexWrap: 'wrap',
        gap: '16px',
      }"
    >
      <div style="display: flex; align-items: center; gap: 16px;">
        <img :src="teamLogo" style="width: 64px; height: 64px;" />
        <div>
          <div style="display: flex; align-items: center; gap: 8px;">
            <h1 style="margin: 0;">{{ team?.displayName }}</h1>
            <NSelect
              :options="teamOptions"
              :value="String(route.params.id)"
              @update:value="handleTeamSelect"
              :render-label="renderTeamLabel"
              style="width: 220px;"
            />
          </div>
          <p style="margin: 0;">{{ getStat('overall') }} | {{ ordinal(Number(getStat('playoffSeed'))) }} in {{ conferenceAbbr }}</p>
        </div>
      </div>
      <div style="display: flex; gap: 24px; text-align: center;">
        <div>
          <div style="font-size: 12px;">PPG</div>
          <div style="font-size: 20px; font-weight: bold;">{{ getStat('avgPointsFor') }}</div>
        </div>
        <div>
          <div style="font-size: 12px;">RPG</div>
          <div style="font-size: 20px; font-weight: bold;">{{ getDetailedStat('avgRebounds') }}</div>
        </div>
        <div>
          <div style="font-size: 12px;">APG</div>
          <div style="font-size: 20px; font-weight: bold;">{{ getDetailedStat('avgAssists') }}</div>
        </div>
        <div>
          <div style="font-size: 12px;">OPPG</div>
          <div style="font-size: 20px; font-weight: bold;">{{ getStat('avgPointsAgainst') }}</div>
        </div>
      </div>
    </div>

    <NCard title="Roster">
      <NDataTable :columns="columns" :data="roster" :bordered="false" />
    </NCard>
  </div>
  <div v-else style="text-align: center; margin-top: 80px;">
    <NSpin size="large" />
  </div>
</template>
