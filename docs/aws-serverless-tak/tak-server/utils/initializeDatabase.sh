#!/bin/bash

echo "Checking database configuration"

export PGPASSWORD=$DB_PASSWORD
java -jar /opt/tak/db-utils/SchemaManager.jar validate >/dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "Database has already been initialized"
else
    echo "Database not yet configured. Initializing database"
    psql --host "$DB_HOST" --dbname postgres --username postgres --command="CREATE DATABASE cot WITH OWNER postgres"
    java -jar /opt/tak/db-utils/SchemaManager.jar SetupRds
    java -jar /opt/tak/db-utils/SchemaManager.jar upgrade
    echo "Database initialized"
fi