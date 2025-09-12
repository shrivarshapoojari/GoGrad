-- Create teachers table for GoGrad database
-- Run this SQL script in your MySQL database

CREATE DATABASE IF NOT EXISTS gograd;
USE gograd;

CREATE TABLE IF NOT EXISTS teachers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    class VARCHAR(20) NOT NULL,
    subject VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Insert some sample data
INSERT INTO teachers (first_name, last_name, class, subject) VALUES
('John', 'Doe', '9', 'Math'),
('Jane', 'Smith', '10', 'Science'),
('Alice', 'Johnson', '11', 'History'),
('Bob', 'Wilson', '12', 'English');