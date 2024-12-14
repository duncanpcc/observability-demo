#!/bin/bash

#parameters
SQL_SERVER_SA_PASSWORD=$1
TELEGRAF_SQL_SERVER_PASSWORD=$2

/opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P ${SQL_SERVER_SA_PASSWORD} -d master -N -C -Q "
CREATE LOGIN [telegraf] WITH PASSWORD = N'${TELEGRAF_SQL_SERVER_PASSWORD}';
GO
GRANT VIEW SERVER STATE TO [telegraf];
GO
GRANT VIEW ANY DEFINITION TO [telegraf];
GO
"