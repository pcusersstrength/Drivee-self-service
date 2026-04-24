def prompt_builder(
    user_query: str, dialect: str = "postgresql", table_meta: str | None = None
) -> str:
    schema = """TABLE NAME: anonymized_incity_orders

COLUMNS:
- city_id (integer) - city identifier
- offset_hours (integer) - UTC offset in hours for local time
- order_id (string) - anonymized order ID
- tender_id (string) - anonymized tender ID (one order may have multiple tenders)
- user_id (string) - anonymized user (client) ID
- driver_id (string) - anonymized driver ID
- status_order (string) - final order status. EXACT values: 'done' (completed), 'cancel' (cancelled by client or driver), 'accept' (accepted but not finished), 'delete' (deleted)
- status_tender (string) - tender/driver matching status. EXACT values: 'done', 'decline', 'accept'
- order_timestamp (timestamp, UTC) - order creation time
- tender_timestamp (timestamp, UTC) - tender (driver search) start time
- driveraccept_timestamp (timestamp, UTC) - driver acceptance time
- driverarrived_timestamp (timestamp, UTC) - driver arrival time at pickup
- driverstarttheride_timestamp (timestamp, UTC) - ride start time
- driverdone_timestamp (timestamp, UTC) - ride completion time
- clientcancel_timestamp (timestamp, UTC) - client cancellation time
- drivercancel_timestamp (timestamp, UTC) - driver cancellation time
- order_modified_local (timestamp, local) - last modification (already with offset_hours applied)
- cancel_before_accept_local (timestamp, local) - cancellation before acceptance
- distance_in_meters (float) - trip distance in meters
- duration_in_seconds (float) - trip duration in seconds
- price_order_local (numeric) - final order price in local currency
- price_tender_local (numeric) - tender stage price
- price_start_local (numeric) - starting price

IMPORTANT RULES:
- One row = one order_id + one tender_id. A single order may have multiple tenders. For per-order stats use COUNT(DISTINCT order_id).
- All *_timestamp columns are in UTC. For local time: timestamp + INTERVAL '1 hour' * offset_hours, or use the *_local columns.
- Prices are in each city's local currency (price_order_local).
- NULL in a timestamp = event did not happen (e.g. driverdone_timestamp IS NULL means the ride was never finished).
- The data spans roughly March-April 2026. For queries like "last week", "yesterday", "recent" — anchor on MAX(order_timestamp) in the table, NOT on NOW().
- EXACT status values (do not invent others): status_order IN ('done', 'cancel', 'accept', 'delete'); status_tender IN ('done', 'decline', 'accept'). Note: it's 'cancel' (5 letters), NOT 'cancelled' or 'canceled'.
- For percentage/ratio queries involving cancellations, DO NOT filter by status_order='done' — this removes all cancellations and gives 0. Either use no status filter, or filter by the status you are measuring.
- GENERATE ONLY SELECT statements. No INSERT/UPDATE/DELETE/DROP/ALTER/TRUNCATE.
- For queries that return many rows (plain SELECT without aggregation, or GROUP BY that yields multiple rows) — add LIMIT 1000 if no explicit limit is given.
- For aggregate queries that inherently return a single row (SELECT COUNT/SUM/AVG/MAX/MIN without GROUP BY) — do NOT add LIMIT.
- Output exactly ONE SQL statement. Do not output multiple statements separated by semicolons.

SEMANTIC DICTIONARY (Russian term → SQL mapping):
- «выручка», «доход», «сумма» → SUM(price_order_local) WHERE status_order = 'done'
- «средний чек» → AVG(price_order_local) WHERE status_order = 'done'
- «количество поездок», «выполненные заказы», «поездки» → COUNT(DISTINCT order_id) WHERE status_order = 'done'
- «отменённые заказы», «отмены» → COUNT(DISTINCT order_id) WHERE status_order = 'cancel'
- «отмена клиентом» → clientcancel_timestamp IS NOT NULL
- «отмена водителем» → drivercancel_timestamp IS NOT NULL
- «водители» → driver_id, «клиенты» / «пользователи» → user_id, «города» → city_id
- «время ожидания водителя» → driverarrived_timestamp - driveraccept_timestamp
- «время в пути» → duration_in_seconds (or driverdone_timestamp - driverstarttheride_timestamp)
- «топ N» → ORDER BY ... DESC LIMIT N"""

    few_shot_examples = """EXAMPLES:

User: Покажи топ-5 городов по количеству выполненных заказов
SQL: SELECT city_id, COUNT(DISTINCT order_id) AS completed_orders FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY city_id ORDER BY completed_orders DESC LIMIT 5;

User: Какая выручка по городам за последнюю неделю?
SQL: SELECT city_id, SUM(price_order_local) AS revenue FROM anonymized_incity_orders WHERE status_order = 'done' AND order_timestamp >= (SELECT MAX(order_timestamp) - INTERVAL '7 days' FROM anonymized_incity_orders) GROUP BY city_id ORDER BY revenue DESC LIMIT 1000;

User: Топ-10 водителей по количеству поездок
SQL: SELECT driver_id, COUNT(DISTINCT order_id) AS trips FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY driver_id ORDER BY trips DESC LIMIT 10;

User: Сколько заказов отменили вчера?
SQL: SELECT COUNT(DISTINCT order_id) AS cancelled FROM anonymized_incity_orders WHERE status_order = 'cancel' AND order_timestamp::date = (SELECT MAX(order_timestamp)::date - INTERVAL '1 day' FROM anonymized_incity_orders);

User: Средний чек по городам
SQL: SELECT city_id, AVG(price_order_local) AS avg_check FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY city_id ORDER BY avg_check DESC LIMIT 1000;

User: Динамика заказов по дням за последний месяц
SQL: SELECT DATE(order_timestamp) AS day, COUNT(DISTINCT order_id) AS orders FROM anonymized_incity_orders WHERE order_timestamp >= (SELECT MAX(order_timestamp) - INTERVAL '30 days' FROM anonymized_incity_orders) GROUP BY day ORDER BY day ASC LIMIT 1000;

User: Сравни выручку самого прибыльного и самого убыточного города
SQL: (SELECT city_id, SUM(price_order_local) AS revenue FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY city_id ORDER BY revenue DESC LIMIT 1) UNION ALL (SELECT city_id, SUM(price_order_local) AS revenue FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY city_id ORDER BY revenue ASC LIMIT 1);

User: Какой процент заказов отменяется клиентами до принятия водителем?
SQL: SELECT ROUND(100.0 * COUNT(DISTINCT CASE WHEN clientcancel_timestamp IS NOT NULL AND driveraccept_timestamp IS NULL THEN order_id END) / NULLIF(COUNT(DISTINCT order_id), 0), 2) AS cancel_before_accept_pct FROM anonymized_incity_orders;

User: Какой процент заказов отменяется клиентами?
SQL: SELECT ROUND(100.0 * COUNT(DISTINCT CASE WHEN status_order = 'cancel' AND clientcancel_timestamp IS NOT NULL THEN order_id END) / NULLIF(COUNT(DISTINCT order_id), 0), 2) AS client_cancel_pct FROM anonymized_incity_orders;

User: Сколько уникальных клиентов сделали больше 5 поездок?
SQL: SELECT COUNT(*) AS users_with_more_than_5_trips FROM (SELECT user_id FROM anonymized_incity_orders WHERE status_order = 'done' GROUP BY user_id HAVING COUNT(DISTINCT order_id) > 5) t;"""

    PROMPT = f"""You are a SQL expert. Translate a user question (written in Russian or English) into a single valid {dialect.upper()} query.

Output format — strict:
- Only the SQL code, no explanations, no markdown fences, no comments.
- One statement only, ending with a semicolon.

{schema}

{few_shot_examples}

User: {user_query}
SQL:"""

    return PROMPT
