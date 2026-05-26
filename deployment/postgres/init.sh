#!/bin/bash
set -e

# Create additional databases for each domain service
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE wms_customer;
    CREATE DATABASE wms_inbound;
    CREATE DATABASE wms_inventory;
    CREATE DATABASE wms_order;
EOSQL

echo "Databases created: wms_customer, wms_inbound, wms_inventory, wms_order"
