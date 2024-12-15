CREATE OR ALTER PROCEDURE dbo.demo_workload
AS
BEGIN	

    --randomly select 1 of 4 operations to perform --using a value from 1 to 100 to determine the operation
    DECLARE @operation INT = (SELECT ABS(CHECKSUM(NEWID()) % 100) + 1)
    IF @operation >= 1 AND @operation <= 50
    BEGIN
        PRINT 'selecting a random row'
        --determine max id
        DECLARE @max_id INT = (SELECT MAX(id) FROM demo_workload.dbo.demo_table)
        --SELECT a random row between 1 and max id
        SELECT * FROM demo_workload.dbo.demo_table WHERE id = (SELECT ABS(CHECKSUM(NEWID()) % @max_id) + 1)
    END
    ELSE IF @operation >= 51 AND @operation <= 74
    BEGIN
        PRINT 'inserting a random row'
        --insert a random row
        DECLARE @name NVARCHAR(100) = (SELECT 'name' + CAST(ABS(CHECKSUM(NEWID()) % 100) + 1 AS NVARCHAR(3)))
        DECLARE @value INT = (SELECT ABS(CHECKSUM(NEWID()) % 100) + 1)
        INSERT INTO demo_workload.dbo.demo_table (name, value) VALUES (@name, @value)
    END
    ELSE IF @operation >= 75 AND @operation <= 89
    BEGIN
        PRINT 'updating a random row'
        --update a random row
        DECLARE @max_id2 INT = (SELECT MAX(id) FROM demo_workload.dbo.demo_table)
        DECLARE @id INT = (SELECT ABS(CHECKSUM(NEWID()) % @max_id2) + 1)
        DECLARE @newname NVARCHAR(100) = (SELECT 'name' + CAST(ABS(CHECKSUM(NEWID()) % 100) + 1 AS NVARCHAR(3)))
        DECLARE @newvalue INT = (SELECT ABS(CHECKSUM(NEWID()) % 100) + 1)
        UPDATE demo_workload.dbo.demo_table SET name = @newname, value = @newvalue WHERE id = @id
    END
    ELSE IF @operation >= 90 AND @operation <= 100
    BEGIN
        PRINT 'deleting a random row'
        --delete a random row
        DECLARE @max_id3 INT = (SELECT MAX(id) FROM demo_workload.dbo.demo_table)
        DECLARE @id2 INT = (SELECT ABS(CHECKSUM(NEWID()) % @max_id3) + 1)
        DELETE FROM demo_workload.dbo.demo_table WHERE id = @id2
    END
END