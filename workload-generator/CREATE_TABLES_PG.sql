-- Create a table within the schema if it doesn't exist
CREATE TABLE IF NOT EXISTS demo_workload.demo_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    value INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);