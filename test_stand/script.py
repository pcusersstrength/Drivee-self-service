
import pandas as pd
from sqlalchemy import create_engine

# 1. Загружаем CSV в DataFrame
df = pd.read_csv('train.csv')

# 2. Создаем подключение к БД (пример для PostgreSQL)
# Формат: postgresql://username:password@host:port/database
engine = create_engine('postgresql://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable')

# 3. Заливаем данные
# if_exists='append' добавит данные в существующую таблицу
# if_exists='replace' создаст таблицу заново
df.to_sql('table_name', engine, if_exists='append', index=False)
print("hello world")
