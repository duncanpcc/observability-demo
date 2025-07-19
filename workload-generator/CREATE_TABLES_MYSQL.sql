CREATE TABLE demo_workload.dbo.demo_table
(
    id INT IDENTITY(1,1) PRIMARY KEY CLUSTERED,
    name NVARCHAR(100),
    value INT,
    created_at DATETIME DEFAULT GETDATE()
)
