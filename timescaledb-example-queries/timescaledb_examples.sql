select  "time"
,host
,measurement_db_type
,sql_instance
,wait_category
,wait_type
,max_wait_time_ms
,resource_wait_ms
,signal_wait_time_ms
,wait_time_ms
,waiting_tasks_count
from public.sqlserver_waitstats
LIMIT 1000;

SELECT time_bucket('5m', time) AS time, wait_type, (avg(wait_time_ms) + avg(signal_wait_time_ms)) as total_wait_time_ms
FROM sqlserver_waitstats
--WHERE $__timeFilter(time)
GROUP BY time , wait_type
ORDER BY time;

--increase in value over time, accounting for counter resets
select time , wait_type, 
(
    CASE
        --compare current value versus previous value (lag) in this window, if current value is greater than previous value then do the proper math to determine the change
        WHEN  (avg(wait_time_ms) + avg(signal_wait_time_ms)) >= lag((avg(wait_time_ms) + avg(signal_wait_time_ms))) OVER w 
            THEN (avg(wait_time_ms) + avg(signal_wait_time_ms)) - lag((avg(wait_time_ms) + avg(signal_wait_time_ms))) OVER w
        --if the previous value (lag) in the window is null, then don't do any math and simply return null (we could potentially return 0 here depending on what we are trying to accomplish)
        WHEN  lag((avg(wait_time_ms) + avg(signal_wait_time_ms))) OVER w IS NULL 
            THEN NULL
        --if the new value is lesser than the old value, return the new value. This will occur in case of restarts where the counters get zeroed
        ELSE (avg(wait_time_ms) + avg(signal_wait_time_ms))
    END
) as total_wait_time_ms_chg 
,(avg(wait_time_ms) + avg(signal_wait_time_ms)) as actual_ms
from sqlserver_waitstats ws 
--WHERE $__timeFilter(time)
--where wait_Type in ('SOS_SCHEDULER_YIELD')
--where wait_Type in ('SOS_SCHEDULER_YIELD','LCK_M_U')
group by time, wait_type
WINDOW w AS (PARTITION BY wait_type ORDER BY time)
order by time;


--increase in value over time, accounting for counter resets
--calculating rate by dividing by the time difference
select time , wait_type, 
(
    CASE
        --compare current value versus previous value (lag) in this window, if current value is greater than previous value then do the proper math to determine the change
        WHEN  (avg(wait_time_ms) + avg(signal_wait_time_ms)) >= lag((avg(wait_time_ms) + avg(signal_wait_time_ms))) OVER w 
            THEN (avg(wait_time_ms) + avg(signal_wait_time_ms)) - lag((avg(wait_time_ms) + avg(signal_wait_time_ms))) OVER w
        --if the previous value (lag) in the window is null, then don't do any math and simply return null (we could potentially return 0 here depending on what we are trying to accomplish)
        WHEN  lag((avg(wait_time_ms) + avg(signal_wait_time_ms))) OVER w IS NULL 
            THEN NULL
        --if the new value is lesser than the old value, return the new value. This will occur in case of restarts where the counters get zeroed
        ELSE (avg(wait_time_ms) + avg(signal_wait_time_ms))
    END
) / EXTRACT(EPOCH FROM time - lag(time) OVER w) as wait_time_ms_per_s 
, (
    CASE
        --compare current value versus previous value (lag) in this window, if current value is greater than previous value then do the proper math to determine the change
        WHEN  (avg(wait_time_ms) + avg(signal_wait_time_ms)) >= lag((avg(wait_time_ms) + avg(signal_wait_time_ms))) OVER w 
            THEN (avg(wait_time_ms) + avg(signal_wait_time_ms)) - lag((avg(wait_time_ms) + avg(signal_wait_time_ms))) OVER w
        --if the previous value (lag) in the window is null, then don't do any math and simply return null (we could potentially return 0 here depending on what we are trying to accomplish)
        WHEN  lag((avg(wait_time_ms) + avg(signal_wait_time_ms))) OVER w IS NULL 
            THEN NULL
        --if the new value is lesser than the old value, return the new value. This will occur in case of restarts where the counters get zeroed
        ELSE (avg(wait_time_ms) + avg(signal_wait_time_ms))
    END
)  as difference_wait_time_ms 
,(avg(wait_time_ms) + avg(signal_wait_time_ms)) as actual_ms
from sqlserver_waitstats ws 
--WHERE $__timeFilter(time)
--where wait_Type in ('SOS_SCHEDULER_YIELD')
--where wait_Type in ('SOS_SCHEDULER_YIELD','LCK_M_U')
group by time, wait_type
WINDOW w AS (PARTITION BY wait_type ORDER BY time)
order by time;

--display the last value
--select * from sqlserver_Server_properties limit 10
select LAST(cpu_count, time),LAST(sql_version, time) as version, LAST(sku,time) as sku
from sqlserver_Server_properties


--take a look at the different tables and the data
select * from sqlserver_cpu limit 5;
select * from sqlserver_database_io limit 5;
select * from sqlserver_memory_clerks limit 5;
select * from sqlserver_performance limit 5;
select * from sqlserver_recentbackup limit 5;
select * from sqlserver_requests limit 5;
select * from sqlserver_schedulers limit 5;
select * from sqlserver_server_properties limit 5;
select * from sqlserver_volume_space limit 5;
select * from sqlserver_waitstats limit 5;