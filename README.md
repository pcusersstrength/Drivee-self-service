# Drivee-self-service

<img width="1536" height="1024" alt="image" src="https://github.com/user-attachments/assets/8d50fe73-03ed-4a4a-8329-f2d22666d063" />

https://drivee-ai.ru/

MVP интеллектуальной платформы, которая позволяет пользователям получать данные из базы данных **без знания SQL**, **используя естественный язык**.

---

## Возможности
- Чат с локальной ИИ-моделью в реальном времени
- Простая и приятная веб-интерфейс (тёмная/светлая тема, иконки, отправка сообщений)
- Генерация SQL запросов из текста на естественном языке
- Получение данных из базы данных 
- Графики и статистика

## Технологический стек
- Backend: Go + Gorilla WebSocket + SQLite/PostgreSQL
- AI Service: Python + FastAPI (или аналог) + LLM deepseek-coder
- Frontend: Чистый HTML + CSS + JavaScript + WebSocket
- Контейнеризация: Docker + Docker Compose
- Логирование: slog (с красивым выводом)

## Запуск

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/pcusersstrength/Drivee-self-service.git
cd Drivee-self-service
```

### 2. Настройте токен

```bash
cat > .env <<'ENV'
TOKEN=your_secret_token_here
ENV

```


### 3. Запустите докер

```bash
docker compose up -d
```

### 4. Импортируйте модель

```bash
docker exec ollama ollama pull deepseek-coder:6.7b-instruct-q4_K_M
```

### 5. Импортируйте датасет
добавьте файл train.csv в файлы проекта и запустите скрипт
```bash
python script.py
```
