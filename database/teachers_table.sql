-- Create teachers table for GoGrad database
-- Run this SQL script in your MySQL database

CREATE DATABASE IF NOT EXISTS gograd;
USE gograd;

CREATE TABLE IF NOT EXISTS teachers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    class VARCHAR(20) NOT NULL,
    subject VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Insert some sample data
INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES
('John', 'Doe', 'john.doe@school.edu', '9', 'Math'),
('Jane', 'Smith', 'jane.smith@school.edu', '10', 'Science'),
('Alice', 'Johnson', 'alice.johnson@school.edu', '11', 'History'),
('Bob', 'Wilson', 'bob.wilson@school.edu', '12', 'English');