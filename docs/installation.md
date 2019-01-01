
```
$ mysql -u root -p

mysql> create database taxi_scout;
Query OK, 1 row affected (0.01 sec)

mysql> CREATE USER 'taxi_scout_user'@'localhost' IDENTIFIED BY 'taxi_scout_pwd';
Query OK, 0 rows affected (0.01 sec)

mysql> GRANT ALL PRIVILEGES ON taxi_scout.* TO 'taxi_scout_user'@'localhost';
Query OK, 0 rows affected (0.00 sec)

```
