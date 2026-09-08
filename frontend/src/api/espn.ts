const SITE_BASE = 'https://site.api.espn.com/apis/site/v2/sports/basketball/nba'
const CORE_BASE = 'https://site.web.api.espn.com/apis/common/v3/sports/basketball/nba'


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