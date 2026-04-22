package postgres


// WITH tables AS (
//     SELECT 
//         table_schema,
//         table_name
//     FROM information_schema.tables 
//     WHERE table_schema NOT IN ('information_schema', 'pg_catalog')
//       AND table_type = 'BASE TABLE'
// )
// SELECT 
//     t.table_schema,
//     t.table_name,
//     c.ordinal_position,
//     c.column_name,
//     c.data_type,
//     c.is_nullable,
//     c.column_default,
//     col_description((t.table_schema || '.' || t.table_name)::regclass, c.ordinal_position) AS column_comment,
//     obj_description((t.table_schema || '.' || t.table_name)::regclass) AS table_comment
// FROM tables t
// JOIN information_schema.columns c 
//     ON c.table_schema = t.table_schema 
//    AND c.table_name = t.table_name
// ORDER BY t.table_schema, t.table_name, c.ordinal_position;