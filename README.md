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
