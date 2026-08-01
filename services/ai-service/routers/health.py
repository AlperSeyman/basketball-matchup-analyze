from fastapi import APIRouter, Response
from redis.exceptions import RedisError
from redis_client import redis_client


router = APIRouter()

@router.get("/health")
async def health_check(response: Response):
    try:
        await redis_client.ping()
    except RedisError:
        response.status_code = 503
        return {"status":"down", "redis": "down"}
    return {"status": "ok", "redis": "ok"}


