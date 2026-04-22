def prompt_builder(user_query: str, dialect: str = "postgresql") -> str:
    PROMPT = f"""You are a SQL generator. Follow this format exactly:

User: [question in natural language]
SQL: [only SQL code, no explanations, no markdown]

Database schema for anonymized incity orders:
- city_id (integer) - city identifier
- offset_hours (integer) - UTC offset in hours for local time
- order_id (string) - anonymized order ID
- tender_id (string) - anonymized tender ID
- user_id (string) - anonymized user ID
- driver_id (string) - anonymized driver ID
- status_order (string) - final order status
- status_tender (string) - tender/driver matching status
- order_timestamp (timestamp) - order creation time
- tender_timestamp (timestamp) - tender start time
- driveraccept_timestamp (timestamp) - driver acceptance time
- driverarrived_timestamp (timestamp) - driver arrival time
- driverstarttheride_timestamp (timestamp) - ride start time
- driverdone_timestamp (timestamp) - ride completion time
- clientcancel_timestamp (timestamp) - client cancellation time
- drivercancel_timestamp (timestamp) - driver cancellation time
- order_modified_local (timestamp) - last modification time (local)
- cancel_before_accept_local (timestamp) - cancellation before acceptance
- distance_in_meters (float) - trip distance in meters
- duration_in_seconds (float) - trip duration in seconds
- price_order_local (decimal) - final order price in local currency
- price_tender_local (decimal) - tender stage price in local currency
- price_start_local (decimal) - starting price in local currency

Important notes:
- One row = one order_id + tender_id combination
- For time-based queries, consider offset_hours for local time conversion
- Use {dialect.upper()} syntax

User: {user_query}
SQL:"""

    return PROMPT
