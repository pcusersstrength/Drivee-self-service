import hashlib
import os
import re
from collections import OrderedDict

import requests
from dotenv import load_dotenv
from fastapi import Body, FastAPI, Header, HTTPException, Query
from fastapi.middleware.cors import CORSMiddleware

from utils import prompt_builder

load_dotenv()

TOKEN = os.getenv("TOKEN")
API_URL = os.getenv("API_URL")
CACHE_MAX = 256

_sql_cache: "OrderedDict[str, dict]" = OrderedDict()

app = FastAPI(title="SQL Generator API")
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

FORBIDDEN_KEYWORDS = [
    "DROP",
    "DELETE",
    "TRUNCATE",
    "ALTER",
    "UPDATE",
    "INSERT",
    "GRANT",
    "REVOKE",
    "CREATE",
]


async def verify_token(token: str):
    if token != TOKEN:
        raise HTTPException(status_code=401, detail="Invalid token")
    return True


DEFAULT_COLUMNS = {
    "anonymized_incity_orders": "order_id, city_id, status_order, price_order_local, order_timestamp",
    "passenger_daily_metrics": "user_id, city_id, order_date_part, orders_count, rides_count",
    "driver_daily_metrics": "driver_id, city_id, tender_date_part, rides_count, online_time_sum_seconds",
}


def fix_select_star(sql: str) -> str:
    pattern = re.compile(
        r"SELECT\s+\*\s+FROM\s+(\w+)",
        flags=re.IGNORECASE,
    )

    def replacer(match: "re.Match[str]") -> str:
        table = match.group(1).lower()
        cols = DEFAULT_COLUMNS.get(table)
        if cols:
            return f"SELECT {cols} FROM {match.group(1)}"
        return match.group(0)

    return pattern.sub(replacer, sql)


def postprocess_sql(raw: str) -> str:
    sql = raw.strip()
    sql = sql.replace("```sql", "").replace("```", "").strip()
    sql = sql.split("User:")[0].strip()

    if ";" in sql:
        sql = sql.split(";")[0].strip() + ";"

    sql = re.sub(r"\s+", " ", sql).strip()
    sql = fix_select_star(sql)
    return sql


def validate_sql(sql: str) -> None:
    upper = sql.upper()
    for keyword in FORBIDDEN_KEYWORDS:
        if re.search(rf"\b{keyword}\b", upper):
            raise HTTPException(
                status_code=400,
                detail=f"Forbidden SQL operation detected: {keyword}. Only SELECT is allowed.",
            )
    if not re.match(r"^\s*(\(|WITH\s|SELECT\b)", upper):
        raise HTTPException(
            status_code=400,
            detail="Only SELECT / WITH / parenthesized SELECT statements are allowed.",
        )


def call_ollama(prompt: str) -> str:
    try:
        response = requests.post(
            f"{API_URL}/api/generate",
            json={
                "model": "deepseek-coder:6.7b-instruct-q4_K_M",
                "prompt": prompt,
                "stream": False,
                "options": {
                    "temperature": 0.1,
                    "num_predict": 500,
                    "num_ctx": 6144,
                    "top_k": 40,
                    "top_p": 0.9,
                    "repeat_penalty": 1.0,
                    "seed": 42,
                    "num_thread": 16,
                    "num_batch": 512,
                    "stop": ["\nUser:", "\n\n", "```"],
                },
            },
            timeout=180,
        )
    except requests.exceptions.Timeout:
        raise HTTPException(status_code=504, detail="Request timeout")
    except requests.exceptions.ConnectionError:
        raise HTTPException(status_code=503, detail="Cannot connect to Ollama")

    if response.status_code != 200:
        raise HTTPException(
            status_code=response.status_code,
            detail=f"Ollama API error: {response.status_code}",
        )

    return response.json().get("response", "")


def cache_key(q: str, dialect: str, table_meta: str | None) -> str:
    normalized = f"{q.strip().lower()}|{dialect.lower()}|{(table_meta or '').strip()}"
    return hashlib.md5(normalized.encode("utf-8")).hexdigest()


def build_response(q: str, dialect: str, table_meta: str | None) -> dict:
    key = cache_key(q, dialect, table_meta)
    if key in _sql_cache:
        _sql_cache.move_to_end(key)
        cached = dict(_sql_cache[key])
        cached["cached"] = True
        return cached

    prompt = prompt_builder(q, dialect, table_meta)
    raw_sql = call_ollama(prompt)
    sql = postprocess_sql(raw_sql)
    validate_sql(sql)

    result = {
        "success": True,
        "question": q,
        "dialect": dialect,
        "sql": sql,
        "cached": False,
    }
    _sql_cache[key] = result
    if len(_sql_cache) > CACHE_MAX:
        _sql_cache.popitem(last=False)
    return result


@app.get("/api/sql")
async def generate_sql(
    q: str,
    dialect: str = Query(
        default="postgresql", description="SQL dialect: postgresql, mysql, sqlite"
    ),
    table_meta: str | None = Body(None, description="Custom table schema (optional)"),
    authorization: str | None = Header(None),
):
    if not authorization:
        raise HTTPException(status_code=401, detail="Authorization header required")
    await verify_token(authorization.replace("Bearer ", ""))
    return build_response(q, dialect, table_meta)


@app.post("/api/sql")
async def generate_sql_post(
    request_data: dict,
    authorization: str = Header(None),
):
    if not authorization:
        raise HTTPException(status_code=401, detail="Authorization header required")
    await verify_token(authorization.replace("Bearer ", ""))

    q = request_data.get("q")
    if not q:
        raise HTTPException(status_code=400, detail="Field 'q' is required")

    dialect = request_data.get("dialect", "postgresql")
    table_meta = request_data.get("table_meta", None)
    return build_response(q, dialect, table_meta)


@app.post("/api/cache/clear")
async def cache_clear(authorization: str = Header(None)):
    if not authorization:
        raise HTTPException(status_code=401, detail="Authorization header required")
    await verify_token(authorization.replace("Bearer ", ""))
    _sql_cache.clear()
    return {"success": True, "cache_size": 0}


@app.get("/api/cache/stats")
async def cache_stats(authorization: str = Header(None)):
    if not authorization:
        raise HTTPException(status_code=401, detail="Authorization header required")
    await verify_token(authorization.replace("Bearer ", ""))
    return {"cache_size": len(_sql_cache), "cache_max": CACHE_MAX}


@app.get("/health")
async def health_check():
    try:
        response = requests.get(f"{API_URL}/api/tags", timeout=5)
        if response.status_code == 200:
            return {"status": "healthy", "ollama": "connected"}
    except Exception:
        pass
    return {"status": "unhealthy", "ollama": "disconnected"}
