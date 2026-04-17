-- Create the database if it doesn't exist
-- This script ensures the simplebank database is created
SELECT 'CREATE DATABASE simplebank'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'simplebank')\gexec