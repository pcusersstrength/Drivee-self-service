def prompt_builder(
    user_query: str, dialect: str = "postgresql", table_meta: str | None = None
) -> str:
    schema = """DATABASE: 3 tables.

TABLE 1 (DEFAULT — use this unless the question explicitly mentions daily metrics):
  anonymized_incity_orders — order/tender events. One row = one order_id + one tender_id.

COLUMNS:
- city_id (integer) - city identifier
- offset_hours (integer) - UTC offset in hours
- order_id (string) - order ID
- tender_id (string) - tender ID (one order may have multiple tenders)
- user_id (string) - passenger ID
- driver_id (string) - driver ID
- status_order (string) - EXACT values: 'done', 'cancel', 'accept', 'delete'
- status_tender (string) - EXACT values: 'done', 'decline', 'accept'
- order_timestamp (timestamp UTC) - order creation
- tender_timestamp (timestamp UTC) - tender start
- driveraccept_timestamp (timestamp UTC) - driver acceptance
- driverarrived_timestamp (timestamp UTC) - driver arrival at pickup
- driverstarttheride_timestamp (timestamp UTC) - ride start
- driverdone_timestamp (timestamp UTC) - ride completion
- clientcancel_timestamp (timestamp UTC) - client cancellation
- drivercancel_timestamp (timestamp UTC) - driver cancellation
- order_modified_local (timestamp local) - last modification (already with offset_hours)
- cancel_before_accept_local (timestamp local) - cancellation before acceptance
- distance_in_meters (float) - trip distance, meters
- duration_in_seconds (float) - trip duration, seconds
- price_order_local (numeric) - final order price, local currency
- price_tender_local (numeric) - tender stage price
- price_start_local (numeric) - starting price

TABLE 2:
  passenger_daily_metrics — daily aggregated stats per passenger per city.
  One row = one user_id in one city on one date.

COLUMNS:
- city_id (integer)
- user_id (string)
- order_date_part (date) - day for which metrics are computed (LOCAL date)
- user_reg_date (date) - passenger registration date
- orders_count (integer) - distinct orders made by user that day
- orders_cnt_with_tenders (integer) - of those, how many had tenders
- orders_cnt_accepted (integer) - of those, how many were accepted by drivers
- rides_count (integer) - completed rides that day
- rides_time_sum_seconds (float) - total ride time, seconds
- online_time_sum_seconds (float) - total online time, seconds
- client_cancel_after_accept (integer) - count of cancellations by client AFTER driver accepted

TABLE 3:
  driver_daily_metrics — daily aggregated stats per driver per city.
  One row = one driver_id in one city on one date.

COLUMNS:
- city_id (integer)
- driver_id (string)
- tender_date_part (date) - day for which metrics are computed (LOCAL date)
- driver_reg_date (date) - driver registration date
- orders (integer) - orders linked to driver that day  ⚠ NOTE: column is 'orders', NOT 'orders_count'
- orders_cnt_with_tenders (integer)
- orders_cnt_accepted (integer)
- rides_count (integer) - completed rides that day
- rides_time_sum_seconds (float)
- online_time_sum_seconds (float)
- client_cancel_after_accept (integer)

TABLE RELATIONSHIPS:
- anonymized_incity_orders.user_id   = passenger_daily_metrics.user_id
- anonymized_incity_orders.driver_id = driver_daily_metrics.driver_id
- For date join: DATE(anonymized_incity_orders.order_timestamp) ~ passenger_daily_metrics.order_date_part
                 DATE(anonymized_incity_orders.tender_timestamp) ~ driver_daily_metrics.tender_date_part
- Use LEFT JOIN when joining metrics tables — not every order has a metric row.

IMPORTANT RULES:
- NEVER write placeholders like 'table_name', '<table>', 'YOUR_TABLE'. Always use a REAL table name from the schema above (anonymized_incity_orders, passenger_daily_metrics, or driver_daily_metrics).
- If the question is generic and does not specify a table («сколько записей в таблице», «сколько данных в базе», «show all data», «count rows») — use anonymized_incity_orders.
- NEVER use SELECT *. Always list SPECIFIC columns explicitly.
- For "show records / list orders / выведи записи / покажи заказы" type requests — pick MAX 5 columns. Never more. Recommended set: order_id, city_id, status_order, price_order_local, order_timestamp.
- Long UUID-like fields (order_id, tender_id, user_id, driver_id) are OK to include — but at most ONE of them per query.
- DEFAULT table for any "orders / trips / cancellations / revenue" question = anonymized_incity_orders.
- Use passenger_daily_metrics ONLY for questions about daily passenger activity (e.g. "how many orders per user per day", "active passengers", "passenger online time").
- Use driver_daily_metrics ONLY for questions about daily driver activity (e.g. "driver online hours", "active drivers per day", "driver utilization").
- For passenger questions that need ride/order details (price, distance, status) — use anonymized_incity_orders, NOT passenger_daily_metrics.
- The data spans 2025 — April 2026. For "last week", "yesterday", "recent" — anchor on MAX(order_timestamp) (or MAX(order_date_part) / MAX(tender_date_part) when querying daily tables).
- All *_timestamp in anonymized_incity_orders are UTC. *_date_part in metrics tables are LOCAL dates.
- For per-order stats use COUNT(DISTINCT order_id).
- NULL timestamp = event did not happen.
- Cancellations: status_order='cancel' (5 letters!), NOT 'cancelled' or 'canceled'.
- ONLY SELECT statements. No INSERT/UPDATE/DELETE/DROP/ALTER/TRUNCATE.
- For multi-row results — add LIMIT 1000. For single-row aggregates (COUNT/SUM/AVG without GROUP BY) — do NOT add LIMIT.
- Output exactly ONE SQL statement.

CHART RULES (CRITICAL — frontend renders the result based on EXACT column aliases):
- If the question is a TIME SERIES / DYNAMICS / TIMELINE (по дням, неделям, месяцам, часам, динамика, тренд) — use exactly TWO columns:
    * AS date  — for the time bucket (X axis). Example: TO_CHAR(order_timestamp, 'YYYY-MM-DD') AS date
    * AS value — for the numeric metric (Y axis)
  Order by date ASC. → Frontend draws a LINE chart.

- If the question is a TOP-N / RANKING by category (топ городов / водителей / клиентов / and so on, "по городам", "по водителям") — use exactly TWO columns:
    * AS category — for the category (city_id::text, driver_id, ...)
    * AS value    — for the numeric metric
  Order by value DESC. → Frontend draws a BAR chart.

- If the question is a DISTRIBUTION / SHARE across a SMALL set of categories (≤ 7 buckets, e.g. распределение по статусам, доля, breakdown) — use exactly TWO columns:
    * AS label — for the category name (status, type, ...)
    * AS value — for the numeric metric
  → Frontend draws a PIE chart.

- For other queries (single aggregates, raw rows, multi-column SELECTs) — use natural descriptive aliases (cnt, revenue, avg_check, city_id, trips). → Frontend renders a TABLE.

- ALWAYS cast non-text x-axis values to text in date/category/label fields where they may be ambiguous: city_id::text AS category, EXTRACT(HOUR FROM order_timestamp)::text AS date.

SEMANTIC DICTIONARY (Russian → SQL):

Orders / revenue (use anonymized_incity_orders):
- «выручка», «доход» → SUM(price_order_local) WHERE status_order='done'
- «средний чек» → AVG(price_order_local) WHERE status_order='done'
- «количество поездок», «выполненные заказы» → COUNT(DISTINCT order_id) WHERE status_order='done'
- «отменённые заказы», «отмены» → COUNT(DISTINCT order_id) WHERE status_order='cancel'
- «отмена клиентом» → clientcancel_timestamp IS NOT NULL
- «отмена водителем» → drivercancel_timestamp IS NOT NULL
- «время ожидания водителя» → driverarrived_timestamp - driveraccept_timestamp
- «время в пути» → duration_in_seconds OR driverdone_timestamp - driverstarttheride_timestamp
- «топ N» → ORDER BY ... DESC LIMIT N

Passenger activity (use passenger_daily_metrics):
- «активные пассажиры» / «активные клиенты» → COUNT(DISTINCT user_id) FROM passenger_daily_metrics WHERE rides_count > 0
- «онлайн-время пассажира» → SUM(online_time_sum_seconds) FROM passenger_daily_metrics
- «новые пассажиры», «новые регистрации» → user_reg_date >= ...
- «отмены после принятия» → SUM(client_cancel_after_accept)

Driver activity (use driver_daily_metrics):
- «активные водители» → COUNT(DISTINCT driver_id) FROM driver_daily_metrics WHERE rides_count > 0
- «онлайн-время водителя» / «часы онлайн» → SUM(online_time_sum_seconds) / 3600 FROM driver_daily_metrics
- «новые водители» → driver_reg_date >= ...
- «утилизация водителя» → rides_time_sum_seconds / NULLIF(online_time_sum_seconds, 0)
"""

    few_shot_examples = """EXAMPLES (note column aliases — they drive frontend chart selection):

# --- BAR (category + value) ---

User: Покажи топ-5 городов по количеству выполненных заказов
SQL: SELECT city_id::text AS category, COUNT(DISTINCT order_id) AS value FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY city_id ORDER BY value DESC LIMIT 5;

User: Какая выручка по городам за последнюю неделю?
SQL: SELECT city_id::text AS category, SUM(price_order_local) AS value FROM anonymized_incity_orders WHERE status_order = 'done' AND order_timestamp >= (SELECT MAX(order_timestamp) - INTERVAL '7 days' FROM anonymized_incity_orders) GROUP BY city_id ORDER BY value DESC LIMIT 1000;

User: Топ-10 водителей по количеству поездок
SQL: SELECT driver_id AS category, COUNT(DISTINCT order_id) AS value FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY driver_id ORDER BY value DESC LIMIT 10;

User: Сколько активных водителей в каждом городе сегодня?
SQL: SELECT city_id::text AS category, COUNT(DISTINCT driver_id) AS value FROM driver_daily_metrics WHERE tender_date_part = (SELECT MAX(tender_date_part) FROM driver_daily_metrics) AND rides_count > 0 GROUP BY city_id ORDER BY value DESC LIMIT 1000;

User: Сколько часов онлайн в среднем у водителей в каждом городе за последний месяц?
SQL: SELECT city_id::text AS category, ROUND(AVG(online_time_sum_seconds) / 3600.0, 2) AS value FROM driver_daily_metrics WHERE tender_date_part >= (SELECT MAX(tender_date_part) - INTERVAL '30 days' FROM driver_daily_metrics) GROUP BY city_id ORDER BY value DESC LIMIT 1000;

User: Топ-10 пассажиров по количеству поездок за последний месяц
SQL: SELECT user_id AS category, SUM(rides_count) AS value FROM passenger_daily_metrics WHERE order_date_part >= (SELECT MAX(order_date_part) - INTERVAL '30 days' FROM passenger_daily_metrics) GROUP BY user_id ORDER BY value DESC LIMIT 10;

# --- LINE (date + value) ---

User: Динамика заказов по дням за последний месяц
SQL: SELECT TO_CHAR(order_timestamp, 'YYYY-MM-DD') AS date, COUNT(DISTINCT order_id) AS value FROM anonymized_incity_orders WHERE order_timestamp >= (SELECT MAX(order_timestamp) - INTERVAL '30 days' FROM anonymized_incity_orders) GROUP BY date ORDER BY date ASC LIMIT 1000;

User: Сумма выручки по месяцам
SQL: SELECT TO_CHAR(order_timestamp, 'YYYY-MM') AS date, SUM(price_order_local) AS value FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY date ORDER BY date ASC LIMIT 1000;

User: Покажи выручку по часам дня
SQL: SELECT EXTRACT(HOUR FROM order_timestamp)::text AS date, SUM(price_order_local) AS value FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY date ORDER BY date::int ASC LIMIT 24;

# --- PIE (label + value, ≤ 7 buckets) ---

User: Распределение заказов по статусам
SQL: SELECT status_order AS label, COUNT(DISTINCT order_id) AS value FROM anonymized_incity_orders GROUP BY status_order ORDER BY value DESC LIMIT 7;

User: Доля тендеров по статусам
SQL: SELECT status_tender AS label, COUNT(*) AS value FROM anonymized_incity_orders GROUP BY status_tender ORDER BY value DESC LIMIT 7;

# --- TABLE (single aggregates / multi-column / raw rows) ---

User: Сколько записей в таблице?
SQL: SELECT COUNT(*) AS total_rows FROM anonymized_incity_orders;

User: Сколько данных в базе?
SQL: SELECT COUNT(*) AS total_rows FROM anonymized_incity_orders;

User: Сколько всего заказов в базе?
SQL: SELECT COUNT(DISTINCT order_id) AS total_orders FROM anonymized_incity_orders;

User: Какая общая выручка за всё время?
SQL: SELECT SUM(price_order_local) AS total_revenue FROM anonymized_incity_orders WHERE status_order = 'done';

User: Сколько уникальных клиентов сделали больше 5 поездок?
SQL: SELECT COUNT(*) AS users_with_5plus_trips FROM (SELECT user_id FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY user_id HAVING COUNT(DISTINCT order_id) > 5) t;

User: Покажи 10 самых дорогих заказов
SQL: SELECT order_id, city_id, price_order_local, distance_in_meters, status_order FROM anonymized_incity_orders ORDER BY price_order_local DESC NULLS LAST LIMIT 10;

User: Выведи первые 10 записей
SQL: SELECT order_id, city_id, status_order, price_order_local, order_timestamp FROM anonymized_incity_orders LIMIT 10;

User: Сравни выручку самого прибыльного и самого убыточного города
SQL: (SELECT city_id, SUM(price_order_local) AS revenue FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY city_id ORDER BY revenue DESC LIMIT 1) UNION ALL (SELECT city_id, SUM(price_order_local) AS revenue FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY city_id ORDER BY revenue ASC LIMIT 1);
"""

    PROMPT = f"""You are a SQL expert. Translate a user question (in Russian or English) into a single valid {dialect.upper()} query against the database described below.

Output format — strict:
- Only the SQL code, no explanations, no markdown fences, no comments.
- One statement only, ending with a semicolon.

{schema}

{few_shot_examples}

User: {user_query}
SQL:"""

    return PROMPT
