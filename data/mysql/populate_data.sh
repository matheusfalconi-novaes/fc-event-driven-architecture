#!/bin/bash

mysql --host mysql-db -u root --port 3306 -p"${MYSQL_ROOT_PASSWORD}" < wallet-db.sql