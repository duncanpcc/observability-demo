USE master

if not exists (select * from sys.databases where name = 'demo_workload')
BEGIN
   CREATE DATABASE demo_workload;
END

ALTER DATABASE demo_workload SET RECOVERY SIMPLE;

