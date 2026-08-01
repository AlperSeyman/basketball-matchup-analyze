from redis.asyncio import Redis
from config import load

_config = load()
redis_client = Redis.from_url(_config.redis_url)

