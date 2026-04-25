import pandas as pd
from sqlalchemy import create_engine, types

# ====================== НАСТРОЙКИ ======================
DB_URI = 'postgresql://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable'

engine = create_engine(DB_URI)

# ====================== 1. ЗАГРУЗКА incity.csv (train.csv) ======================
print("Загружаем incity.csv (детальные заказы и тендеры)...")

df_incity = pd.read_csv('train.csv')   # или 'incity.csv'

# Приведение типов
df_incity['city_id'] = pd.to_numeric(df_incity['city_id'], errors='coerce').astype('Int64')
df_incity['offset_hours'] = pd.to_numeric(df_incity['offset_hours'], errors='coerce').astype('Int64')

df_incity['distance_in_meters'] = pd.to_numeric(df_incity['distance_in_meters'], errors='coerce')
df_incity['duration_in_seconds'] = pd.to_numeric(df_incity['duration_in_seconds'], errors='coerce')

df_incity['price_order_local'] = pd.to_numeric(df_incity['price_order_local'], errors='coerce')
df_incity['price_tender_local'] = pd.to_numeric(df_incity['price_tender_local'], errors='coerce')
df_incity['price_start_local'] = pd.to_numeric(df_incity['price_start_local'], errors='coerce')

# Тимestamps
timestamp_cols = [
    'order_timestamp', 'tender_timestamp', 'driveraccept_timestamp',
    'driverarrived_timestamp', 'driverstarttheride_timestamp', 'driverdone_timestamp',
    'clientcancel_timestamp', 'drivercancel_timestamp',
    'order_modified_local', 'cancel_before_accept_local'
]

for col in timestamp_cols:
    df_incity[col] = pd.to_datetime(df_incity[col], errors='coerce', utc=True)

# Явное задание типов для PostgreSQL
dtype_incity = {
    'city_id': types.Integer(),
    'offset_hours': types.Integer(),

    'order_id': types.String(),
    'tender_id': types.String(),
    'user_id': types.String(),
    'driver_id': types.String(),
    'status_order': types.String(),
    'status_tender': types.String(),

    'order_timestamp': types.TIMESTAMP(timezone=True),
    'tender_timestamp': types.TIMESTAMP(timezone=True),
    'driveraccept_timestamp': types.TIMESTAMP(timezone=True),
    'driverarrived_timestamp': types.TIMESTAMP(timezone=True),
    'driverstarttheride_timestamp': types.TIMESTAMP(timezone=True),
    'driverdone_timestamp': types.TIMESTAMP(timezone=True),
    'clientcancel_timestamp': types.TIMESTAMP(timezone=True),
    'drivercancel_timestamp': types.TIMESTAMP(timezone=True),
    'order_modified_local': types.TIMESTAMP(timezone=True),
    'cancel_before_accept_local': types.TIMESTAMP(timezone=True),

    'distance_in_meters': types.Float(),
    'duration_in_seconds': types.Float(),

    'price_order_local': types.Numeric(12, 2),
    'price_tender_local': types.Numeric(12, 2),
    'price_start_local': types.Numeric(12, 2),
}

df_incity.to_sql(
    'anonymized_incity_orders',
    engine,
    if_exists='append',
    index=False,
    dtype=dtype_incity,
    chunksize=100_000,          # для больших файлов
    method='multi'
)

print(f"Таблица anonymized_incity_orders загружена. Строк: {len(df_incity):,}")

# # ====================== 2. ЗАГРУЗКА pass_detail.csv ======================
# print("\nЗагружаем pass_detail.csv (метрики пассажиров)...")

# df_pass = pd.read_csv('pass_detail.csv')

# # Приведение типов
# df_pass['city_id'] = pd.to_numeric(df_pass['city_id'], errors='coerce').astype('Int64')
# df_pass['orders_count'] = pd.to_numeric(df_pass['orders_count'], errors='coerce').astype('Int64')
# df_pass['orders_cnt_with_tenders'] = pd.to_numeric(df_pass['orders_cnt_with_tenders'], errors='coerce').astype('Int64')
# df_pass['orders_cnt_accepted'] = pd.to_numeric(df_pass['orders_cnt_accepted'], errors='coerce').astype('Int64')
# df_pass['rides_count'] = pd.to_numeric(df_pass['rides_count'], errors='coerce').astype('Int64')
# df_pass['client_cancel_after_accept'] = pd.to_numeric(df_pass['client_cancel_after_accept'], errors='coerce').astype('Int64')

# df_pass['rides_time_sum_seconds'] = pd.to_numeric(df_pass['rides_time_sum_seconds'], errors='coerce')
# df_pass['online_time_sum_seconds'] = pd.to_numeric(df_pass['online_time_sum_seconds'], errors='coerce')

# df_pass['order_date_part'] = pd.to_datetime(df_pass['order_date_part'], errors='coerce').dt.date
# df_pass['user_reg_date'] = pd.to_datetime(df_pass['user_reg_date'], errors='coerce').dt.date

# dtype_pass = {
#     'city_id': types.Integer(),
#     'user_id': types.String(),
#     'order_date_part': types.Date(),
#     'user_reg_date': types.Date(),

#     'orders_count': types.Integer(),
#     'orders_cnt_with_tenders': types.Integer(),
#     'orders_cnt_accepted': types.Integer(),
#     'rides_count': types.Integer(),
#     'client_cancel_after_accept': types.Integer(),

#     'rides_time_sum_seconds': types.BigInteger(),
#     'online_time_sum_seconds': types.BigInteger(),
# }

# df_pass.to_sql(
#     'passenger_daily_metrics',
#     engine,
#     if_exists='replace',
#     index=False,
#     dtype=dtype_pass,
#     chunksize=100_000,
#     method='multi'
# )

# print(f"Таблица passenger_daily_metrics загружена. Строк: {len(df_pass):,}")

# # ====================== 3. ЗАГРУЗКА driver_detail.csv ======================
# print("\nЗагружаем driver_detail.csv (метрики водителей)...")

# df_driver = pd.read_csv('driver_detail.csv')

# # Приведение типов
# df_driver['city_id'] = pd.to_numeric(df_driver['city_id'], errors='coerce').astype('Int64')

# for col in ['orders', 'orders_cnt_with_tenders', 'orders_cnt_accepted', 
#             'rides_count', 'client_cancel_after_accept']:
#     df_driver[col] = pd.to_numeric(df_driver[col], errors='coerce').astype('Int64')

# for col in ['rides_time_sum_seconds', 'online_time_sum_seconds']:
#     df_driver[col] = pd.to_numeric(df_driver[col], errors='coerce')

# df_driver['tender_date_part'] = pd.to_datetime(df_driver['tender_date_part'], errors='coerce').dt.date
# df_driver['driver_reg_date'] = pd.to_datetime(df_driver['driver_reg_date'], errors='coerce').dt.date

# dtype_driver = {
#     'city_id': types.Integer(),
#     'driver_id': types.String(),
#     'tender_date_part': types.Date(),
#     'driver_reg_date': types.Date(),

#     'orders': types.Integer(),
#     'orders_cnt_with_tenders': types.Integer(),
#     'orders_cnt_accepted': types.Integer(),
#     'rides_count': types.Integer(),
#     'client_cancel_after_accept': types.Integer(),

#     'rides_time_sum_seconds': types.BigInteger(),
#     'online_time_sum_seconds': types.BigInteger(),
# }

# df_driver.to_sql(
#     'driver_daily_metrics',
#     engine,
#     if_exists='replace',
#     index=False,
#     dtype=dtype_driver,
#     chunksize=100_000,
#     method='multi'
# )

# print(f"Таблица driver_daily_metrics загружена. Строк: {len(df_driver):,}")
# print("\nВсе таблицы успешно загружены в PostgreSQL!")