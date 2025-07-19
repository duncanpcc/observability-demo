CREATE OR REPLACE FUNCTION demo_workload.demo_workload()
    RETURNS void
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
DECLARE
    operation INT;
    max_id INT;
    max_id2 INT;
    max_id3 INT;
    random_id INT;    
    id2 INT;
    name VARCHAR(100);
    value INT;
    newname VARCHAR(100);
    newvalue INT;
BEGIN
    -- randomly select 1 of 4 operations to perform using a value from 1 to 100 to determine the operation
    operation := (SELECT FLOOR(RANDOM() * 100) + 1);

    IF operation >= 1 AND operation <= 50 THEN
        RAISE NOTICE 'selecting a random row';
        -- determine max id
        SELECT MAX(id) INTO max_id FROM demo_workload.demo_table;
        -- SELECT a random row between 1 and max id
        random_id := (FLOOR(RANDOM() * max_id) + 1);
        SELECT id into max_id FROM demo_workload.demo_table WHERE id = random_id;
    ELSIF operation >= 51 AND operation <= 74 THEN
        RAISE NOTICE 'inserting a random row';
        -- insert a random row
        name := 'name' || (FLOOR(RANDOM() * 100) + 1);
        value := FLOOR(RANDOM() * 100) + 1;
        INSERT INTO demo_workload.demo_table (name, value) VALUES (name, value);
    ELSIF operation >= 75 AND operation <= 89 THEN
        RAISE NOTICE 'updating a random row';
        -- update a random row
        SELECT MAX(id) INTO max_id2 FROM demo_workload.demo_table;
        random_id := FLOOR(RANDOM() * max_id2) + 1;
        newname := 'name' || (FLOOR(RANDOM() * 100) + 1);
        newvalue := FLOOR(RANDOM() * 100) + 1;
        UPDATE demo_workload.demo_table SET name = newname, value = newvalue WHERE id = random_id;
    ELSIF operation >= 90 AND operation <= 100 THEN
        RAISE NOTICE 'deleting a random row';
        -- delete a random row
        SELECT MAX(id) INTO max_id3 FROM demo_workload.demo_table;
        id2 := FLOOR(RANDOM() * max_id3) + 1;
        DELETE FROM demo_workload.demo_table WHERE id = id2;
    END IF;
END;
$BODY$;

ALTER FUNCTION demo_workload.demo_workload()
    OWNER TO postgres;
