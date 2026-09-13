const SITE_BASE = 'https://site.api.espn.com/apis/site/v2/sports/basketball/nba'
const CORE_BASE = 'https://site.web.api.espn.com/apis/common/v3/sports/basketball/nba'
const SEARCH_BASE = 'https://site.api.espn.com/apis/search/v2'

export async function getStandings(){
    const response = await fetch(`https://site.api.espn.com/apis/v2/sports/basketball/nba/standings`)
    return response.json()
}


export async function getTeamRoster(teamId: string){
    const response = await fetch(`${SITE_BASE}/teams/${teamId}/roster`)
    return response.json()
}

export async function getTeamStatistics(teamId: string){
    const response = await fetch(`${SITE_BASE}/teams/${teamId}/statistics`)
    return response.json()
}

export async function getPlayerStats(playerId: string){
    const response = await fetch(`${CORE_BASE}/athletes/${playerId}/stats`)
    return response.json()
}

export async function getPlayerGamelog(playerId: string){
    const response = await fetch(`${CORE_BASE}/athletes/${playerId}/gamelog`)
    return response.json()
}

export async function searchPlayers(query: string){
    const response = await fetch(`${SEARCH_BASE}?query=${query}&limit=10`)
    const data = await response.json()

    const playerGroup = data.results.find((group: any) => group.type === 'player')
    if (!playerGroup) return []

    return playerGroup.contents
    .filter((item: any) => item.sport === 'basketball' && item.defaultLeagueSlug === 'nba')
    .map((item: any) => ({
        id: item.uid.split('a:')[1],
        name: item.displayName,
        team: item.subtitle
    }))
}