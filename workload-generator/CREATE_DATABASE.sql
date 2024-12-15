USE master

if exists (select * from sys.databases where name = 'demo_workload')
BEGIN
    ALTER DATABASE demo_workload SET SINGLE_USER WITH ROLLBACK IMMEDIATE
    DROP DATABASE demo_workload
END

CREATE DATABASE demo_workload;

