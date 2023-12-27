CREATE DATABASE IF NOT EXISTS analytics;

SET allow_experimental_object_type = 1;

CREATE TABLE IF NOT EXISTS analytics.hits
(
    eventTime DateTime,
    eventType String,
    sessionID String,
	platform String,
	meta JSON
)
ENGINE MergeTree
ORDER BY tuple()
TTL eventTime + INTERVAL 7 DAY;

CREATE TABLE IF NOT EXISTS analytics.visitors_1_day
(
    eventTime DateTime,
    sessionID String,
	platform String,
    totalHits UInt64
)
ENGINE = SummingMergeTree
ORDER BY (eventTime, sessionID)
TTL eventTime + INTERVAL 7 DAY;


CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.visitors_1_day_mv
TO analytics.visitors_1_day
AS
SELECT
    toStartOfDay(eventTime) AS eventTime,
    sessionID,
    platform,
    count(*) as totalHits
    FROM analytics.hits
    WHERE eventType = 'page_view'
    GROUP BY
        sessionID,
        eventTime,
        platform;
