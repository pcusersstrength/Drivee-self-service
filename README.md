# Drivee-self-service

MVP интеллектуальной платформы, которая позволяет пользователям получать данные из базы данных **без знания SQL**, **используя естественный язык**.

---

## Возможности

- Генерация SQL запросов из текста на естественном языке
- Получение данных из базы данных 
- Вывод статистики в виде интерактивных графиков

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

```bash
cat > ./services/backend/config/local.yaml <<'YAML'
env: "local" #local, dev, prod
storage_path: "storage.db"
http_server:
  address: "localhost:8080"
  timeout: 4s
  idle_timeout: 60s

jwt_secret: "your_secret_token_here"
YAML
```
