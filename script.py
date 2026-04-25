import pandas as pd
from sqlalchemy import create_engine, types

# 1. Загружаем CSV
df = pd.read_csv('train.csv')

# 2. Приводим типы

# integer
df['city_id'] = pd.to_numeric(df['city_id'], errors='coerce').astype('Int64')
df['offset_hours'] = pd.to_numeric(df['offset_hours'], errors='coerce').astype('Int64')

# float
df['distance_in_meters'] = pd.to_numeric(df['distance_in_meters'], errors='coerce')
df['duration_in_seconds'] = pd.to_numeric(df['duration_in_seconds'], errors='coerce')

# decimal (оставляем float, но в БД зададим Numeric)
df['price_order_local'] = pd.to_numeric(df['price_order_local'], errors='coerce')
df['price_tender_local'] = pd.to_numeric(df['price_tender_local'], errors='coerce')
df['price_start_local'] = pd.to_numeric(df['price_start_local'], errors='coerce')

# timestamps
timestamp_cols = [
    'order_timestamp',
    'tender_timestamp',
    'driveraccept_timestamp',
    'driverarrived_timestamp',
    'driverstarttheride_timestamp',
    'driverdone_timestamp',
    'clientcancel_timestamp',
    'drivercancel_timestamp',
    'order_modified_local',
    'cancel_before_accept_local'
]

for col in timestamp_cols:
    df[col] = pd.to_datetime(df[col], errors='coerce')

# 3. Подключение
engine = create_engine(
    'postgresql://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable'
)

# 4. Явно задаём типы колонок в PostgreSQL
dtype = {
    'city_id': types.Integer(),
    'offset_hours': types.Integer(),

    'order_id': types.String(),
    'tender_id': types.String(),
    'user_id': types.String(),
    'driver_id': types.String(),
    'status_order': types.String(),
    'status_tender': types.String(),

    'order_timestamp': types.TIMESTAMP(),
    'tender_timestamp': types.TIMESTAMP(),
    'driveraccept_timestamp': types.TIMESTAMP(),
    'driverarrived_timestamp': types.TIMESTAMP(),
    'driverstarttheride_timestamp': types.TIMESTAMP(),
    'driverdone_timestamp': types.TIMESTAMP(),
    'clientcancel_timestamp': types.TIMESTAMP(),
    'drivercancel_timestamp': types.TIMESTAMP(),
    'order_modified_local': types.TIMESTAMP(),
    'cancel_before_accept_local': types.TIMESTAMP(),

    'distance_in_meters': types.Float(),
    'duration_in_seconds': types.Float(),

    'price_order_local': types.Numeric(12, 2),
    'price_tender_local': types.Numeric(12, 2),
    'price_start_local': types.Numeric(12, 2),
}

# 5. Загрузка
df.to_sql(
    'anonymized_incity_orders',
    engine,
    if_exists='replace',   # важно: чтобы пересоздать с типами
    index=False,
    dtype=dtype
)

print("done")