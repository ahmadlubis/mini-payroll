-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create database for testing if not exists
SELECT 'CREATE DATABASE payslip_test_db' 
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'payslip_test_db');