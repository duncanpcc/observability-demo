#!/bin/bash

# Parameters
SQL_SERVER_SA_PASSWORD=${SA_PASSWORD}
TELEGRAF_SQL_SERVER_PASSWORD=${TELEGRAF_SQL_SERVER_PASSWORD}

/opt/mssql/bin/mssql-conf set network.tlscert /etc/ssl/certs/mssql.pem
/opt/mssql/bin/mssql-conf set network.tlskey /etc/ssl/private/mssql.key
/opt/mssql/bin/mssql-conf set network.tlsprotocols 1.2
/opt/mssql/bin/mssql-conf set network.forceencryption 0

# Start SQL Server in the background
/opt/mssql/bin/sqlservr &

# Wait for SQL Server to be ready
echo "Waiting for SQL Server to start... "
until /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P ${SQL_SERVER_SA_PASSWORD}  -N -C -Q "SELECT 1" &> /dev/null
do
  echo "Waiting for SQL Server to start... delaying 1 second"
  sleep 1
done

# Run the initialization script
echo "Running the initialization script..."
/home/init-sqlserver.sh ${SQL_SERVER_SA_PASSWORD} ${TELEGRAF_SQL_SERVER_PASSWORD}

echo "Initialization script executed. Waiting for SQL Server process to end..."

# Wait for SQL Server process to end
wait