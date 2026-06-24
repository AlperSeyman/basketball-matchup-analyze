import base64
import os
from dataclasses import dataclass

import yaml
from sqlalchemy import URL
from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PublicKey
from cryptography.hazmat.primitives.serialization import load_pem_public_key
from dotenv import load_dotenv

load_dotenv()

@dataclass
class Models:
    gemini: str

@dataclass
class Config:
    db_url: URL
    redis_url: str
    jwt_public_key: Ed25519PublicKey
    gemini_api_key: str
    models: Models
    environment: str
    port: int
    rate_limit_per_hour: int
    db_sslmode: str


def load() -> Config:
    db_url = _build_db_url()
    jwt_public_key = _load_jwt_public_key()
    models = _load_models()
    db_sslmode = os.getenv("POSTGRES_SSLMODE", "disable")

    return Config(
        db_url=db_url,
        redis_url=os.getenv("REDIS_URL", "redis://localhost:6379/0"),
        jwt_public_key=jwt_public_key,
        gemini_api_key=_require_env("GEMINI_API_KEY"),
        models=models,
        environment=os.getenv("ENVIRONMENT", "development"),
        port=int(os.getenv("AI_SERVICE_PORT", "8002")),
        rate_limit_per_hour=int(os.getenv("RATE_LIMIT_REQUESTS_PER_HOUR", "10")),
        db_sslmode=db_sslmode
    )


def _build_db_url() -> URL:
    user = _require_env("POSTGRES_USER")
    password = _require_env("POSTGRES_PASSWORD")
    host = os.getenv("POSTGRES_HOST", "localhost")
    port = os.getenv("POSTGRES_PORT", "5432")
    db = _require_env("POSTGRES_DB")

    url = URL.create(
        "postgresql+asyncpg",
        username=user,
        password=password,
        host=host,
        port=int(port),
        database=db
    )
    return url

def _load_jwt_public_key() -> Ed25519PublicKey:
    b64 = _require_env("JWT_PUBLIC_KEY_BASE64")
    pem_bytes = base64.b64decode(b64)
    key = load_pem_public_key(pem_bytes)
    if not isinstance(key, Ed25519PublicKey):
        raise ValueError("JWT_PUBLIC_KEY_BASE64 is not an Ed25519 public key")
    return key


def _load_models() -> Models:
    config_path = os.path.join(os.path.dirname(__file__), "config.yml")
    with open(config_path, "r") as f:
        data = yaml.safe_load(f)
    return Models(gemini=data["models"]["gemini"])


def _require_env(key: str) -> str:
    value = os.getenv(key)
    if not value:
        raise RuntimeError(f"Missing required environment variable: {key}")
    return value