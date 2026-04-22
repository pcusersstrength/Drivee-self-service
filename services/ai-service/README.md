# SQL Generator Service

Сервис для генерации SQL запросов из естественного языка с использованием LLM (DeepSeek-Coder).

## Возможности

- Генерация SQL запросов из текста на естественном языке
- Поддержка разных диалектов SQL (PostgreSQL, MySQL, SQLite)
- Простое развертывание через Docker Compose
- Аутентификация через Bearer token

## Быстрый старт

### 1. Клонируйте репозиторий

```bash
git clone <your-repo-url>
cd sql-generator-service
```

### 2. Настройте токен

```bash
TOKEN=your_secret_token_here
```

### 3. Запустите сервис

```bash
docker compose up -d
```

### 4. Импортируйте модель

```bash
docker exec ollama ollama pull deepseek-coder:6.7b-instruct-q4_K_M
```

### 5. Проверьте работу
```bash
# Health check
curl http://localhost:8000/health

# Генерация SQL
curl -X GET "http://localhost:8000/api/sql?q=find%20users%20over%2018&dialect=postgresql" \
  -H "Authorization: Bearer your_secret_token_here"
```

# API Endpoints

## GET `/api/sql`
Генерация SQL запроса из естественного языка.

### Параметры запроса

| Параметр | Тип | Обязательный | Описание | Пример |
|----------|-----|--------------|----------|--------|
| `q` | string | ✅ Да | Вопрос на естественном языке | `find users over 18` |
| `dialect` | string | ❌ Нет | Диалект SQL (по умолчанию: `postgresql`) | `mysql`, `sqlite`, `postgresql` |

### Заголовки

| Заголовок | Значение | Описание |
|-----------|----------|----------|
| `Authorization` | `Bearer <token>` | Токен авторизации |

### Пример запроса

```bash
curl -X GET "http://localhost:8000/api/sql?q=find%20users%20over%2018&dialect=postgresql" \
  -H "Authorization: Bearer your_secret_token_here"
```

### Пример ответа 
```json
{
  "success": true,
  "question": "find users over 18",
  "dialect": "postgresql",
  "sql": "SELECT * FROM users WHERE age > 18;"
}
```
