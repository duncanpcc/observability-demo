IF EXISTS (SELECT * FROM sys.tables WHERE name = 'demo_table')
BEGIN
    DROP TABLE demo_workload.dbo.demo_table;
END

CREATE TABLE demo_workload.dbo.demo_table
(
    id INT IDENTITY(1,1) PRIMARY KEY CLUSTERED,
    name NVARCHAR(100),
    value INT,
    created_at DATETIME DEFAULT GETDATE()
)
