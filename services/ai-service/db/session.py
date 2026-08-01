import config
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker
from sqlalchemy.orm import DeclarativeBase

_config = config.load()
if _config.db_sslmode == "disable":
    _ssl = False
else:
    _ssl = _config.db_sslmode
engine = create_async_engine(url=_config.db_url, connect_args={"ssl":_ssl})

AsyncSessionLocal = async_sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False
)

class Base(DeclarativeBase):
    pass

