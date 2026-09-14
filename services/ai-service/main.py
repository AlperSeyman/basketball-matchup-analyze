from fastapi import FastAPI
from contextlib import asynccontextmanager
from db.session import engine
from routers import health, analyze
from redis_client import redis_client


@asynccontextmanager
async def lifespan(_app: FastAPI):
    yield
    await engine.dispose()
    await redis_client.aclose()


app = FastAPI(lifespan=lifespan)

app.include_router(health.router)
app.include_router(analyze.router)