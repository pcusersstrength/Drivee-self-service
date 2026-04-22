import os

import requests
from dotenv import load_dotenv
from fastapi import FastAPI, Header, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from utils import prompt_builder

load_dotenv()

TOKEN = os.getenv("TOKEN")
API_URL = os.getenv("API_URL")

app = FastAPI(title="SQL Generator API")
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

async def verify_token(token: str):
    if token != TOKEN:
        raise HTTPException(status_code=401, detail="Invalid token")
    return True


@app.get("/api/sql")
async def generate_sql(
    q: str, dialect: str = "postgresql", authorization: str = Header(None)
):
    if not authorization:
        raise HTTPException(status_code=401, detail="Authorization header required")

    token = authorization.replace("Bearer ", "")
    await verify_token(token)

    prompt = prompt_builder(q, dialect)

    try:
        response = requests.post(
            f"{API_URL}/api/generate",
            json={
                "model": "deepseek-coder:6.7b-instruct-q4_K_M",
                "prompt": prompt,
                "temperature": 0.7,
                "max_tokens": 500,
                "stream": False,
                "options": {
                    "num_thread": 10,
                    "num_batch": 512,
                    "repeat_penalty": 1.2, # Штраф за повторения (помогает, если модель зацикливается)
                    # "top_k": 40,           # Ограничивает выбор слов (уменьшает вариативность)
                    # "top_p": 0.9           # Ядерное сэмплирование
                }
            },
            timeout=80,
        )

        if response.status_code == 200:
            result = response.json()
            sql = result.get("response", "").strip()
            sql = sql.replace("```sql", "").replace("```", "").strip()
            print(sql)
            if sql == "атата":
                return {
                    "success": False,
                    "question": q,
                    "dialect": dialect,
                    "sql": sql,
                }
            return {
                "success": True,
                "question": q,
                "dialect": dialect,
                "sql": sql,
            }
        else:
            raise HTTPException(
                status_code=response.status_code,
                detail=f"Ollama API error: {response.status_code}",
            )

    except requests.exceptions.Timeout:
        raise HTTPException(status_code=504, detail="Request timeout")
    except requests.exceptions.ConnectionError:
        raise HTTPException(status_code=503, detail="Cannot connect to Ollama")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/health")
async def health_check():
    try:
        response = requests.get(f"{API_URL}/api/tags", timeout=5)
        if response.status_code == 200:
            return {"status": "healthy", "ollama": "connected"}
    except Exception:
        pass
    return {"status": "unhealthy", "ollama": "disconnected"}
