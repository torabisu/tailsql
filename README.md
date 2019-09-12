# tailsql

The beginnings of a command line script that runs SQL queries,
outputs them in an ASCII table, and refreshes every X seconds.

Meant to be used in the scenario where you tailing a log table
in a database, so something like:

```sql
SELECT * FROM logs ORDER BY ID desc LIMIT 20;
```


