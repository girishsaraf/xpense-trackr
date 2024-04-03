DROP DATABASE xpense_trackr_db;
-- Create database
CREATE DATABASE xpense_trackr_db;

-- Use the database
USE xpense_trackr_db;

-- Create users table
CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(100) NOT NULL,
                       email VARCHAR(100) NOT NULL,
                       password_hash VARCHAR(255) NOT NULL
);

-- Create investment table
CREATE TABLE investments (
                            id INT AUTO_INCREMENT PRIMARY KEY,
                            amount DECIMAL(10, 2) NOT NULL,
                            investment_type ENUM('stocks', '401k', 'mutual_funds', 'real_estate', 'savings', 'bonds') NOT NULL,
                            date DATE NOT NULL
);

-- Create categories table
CREATE TABLE categories (
                            id INT AUTO_INCREMENT PRIMARY KEY,
                            name VARCHAR(50) UNIQUE NOT NULL
);

-- Create expenses table
CREATE TABLE expenses (id INT AUTO_INCREMENT PRIMARY KEY,
                          category_id INT NOT NULL,
                          amount DECIMAL(10, 2) NOT NULL,
                          description VARCHAR(255) NOT NULL,
                          date DATE NOT NULL,
                          FOREIGN KEY (category_id) REFERENCES categories(id));

-- Create budgets table
CREATE TABLE budgets (id INT AUTO_INCREMENT PRIMARY KEY,
                         category_id INT NOT NULL,
                         amount DECIMAL(10, 2) NOT NULL,
                         FOREIGN KEY (category_id) REFERENCES categories(id));

-- Investments for the past 6 months
-- Insert query for 'stocks'
INSERT INTO investments (amount, investment_type, date) VALUES
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'stocks', '2024-03-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'stocks', '2024-02-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'stocks', '2024-01-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'stocks', '2023-12-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'stocks', '2023-11-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'stocks', '2023-10-01');

-- Insert query for '401k'
INSERT INTO investments (amount, investment_type, date) VALUES
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, '401k', '2024-03-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, '401k', '2024-02-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, '401k', '2024-01-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, '401k', '2023-12-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, '401k', '2023-11-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, '401k', '2023-10-01');

-- Insert query for 'mutual_funds'
INSERT INTO investments (amount, investment_type, date) VALUES
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'mutual_funds', '2024-03-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'mutual_funds', '2024-02-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'mutual_funds', '2024-01-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'mutual_funds', '2023-12-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'mutual_funds', '2023-11-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'mutual_funds', '2023-10-01');

-- Insert query for 'real_estate'
INSERT INTO investments (amount, investment_type, date) VALUES
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'real_estate', '2024-03-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'real_estate', '2024-02-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'real_estate', '2024-01-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'real_estate', '2023-12-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'real_estate', '2023-11-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'real_estate', '2023-10-01');

-- Insert query for 'savings'
INSERT INTO investments (amount, investment_type, date) VALUES
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'savings', '2024-03-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'savings', '2024-02-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'savings', '2024-01-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'savings', '2023-12-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'savings', '2023-11-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'savings', '2023-10-01');

-- Insert query for 'bonds'
INSERT INTO investments (amount, investment_type, date) VALUES
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'bonds', '2024-03-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'bonds', '2024-02-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'bonds', '2024-01-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'bonds', '2023-12-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'bonds', '2023-11-01'),
    (FLOOR(RAND() * (500 - 10 + 1)) + 10, 'bonds', '2023-10-01');


INSERT INTO categories (name) VALUES
('Food'),
('Transportation'),
('Utilities'),
('Entertainment'),
('Healthcare');

-- Insert query for 'Food'
INSERT INTO expenses (category_id, amount, description, date) VALUES
      (1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Groceries', '2024-03-01'),
      (1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Dining out', '2024-02-01'),
      (1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Food delivery', '2024-01-01'),
      (1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Home cooking', '2023-12-01'),
      (1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Snacks', '2023-11-01'),
      (1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Restaurant', '2023-10-01');

-- Insert query for 'Transportation'
INSERT INTO expenses (category_id, amount, description, date) VALUES
      (2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Gasoline', '2024-03-01'),
      (2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Public transportation', '2024-02-01'),
      (2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Taxi', '2024-01-01'),
      (2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Ride-sharing', '2023-12-01'),
      (2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Car maintenance', '2023-11-01'),
      (2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Parking fees', '2023-10-01');

-- Insert query for 'Utilities'
INSERT INTO expenses (category_id, amount, description, date) VALUES
      (3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Electricity', '2024-03-01'),
      (3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Water', '2024-02-01'),
      (3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Internet', '2024-01-01'),
      (3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Phone bill', '2023-12-01'),
      (3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Gas', '2023-11-01'),
      (3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Trash collection', '2023-10-01');

-- Insert query for 'Entertainment'
INSERT INTO expenses (category_id, amount, description, date) VALUES
      (4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Movie tickets', '2024-03-01'),
      (4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Concert', '2024-02-01'),
      (4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Video games', '2024-01-01'),
      (4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Books', '2023-12-01'),
      (4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Streaming service', '2023-11-01'),
      (4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Board games', '2023-10-01');

-- Insert query for 'Healthcare'
INSERT INTO expenses (category_id, amount, description, date) VALUES
      (5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Medications', '2024-03-01'),
      (5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Doctor visit', '2024-02-01'),
      (5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Health insurance', '2024-01-01'),
      (5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Dental care', '2023-12-01'),
      (5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Gym membership', '2023-11-01'),
      (5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Prescription glasses', '2023-10-01');


INSERT INTO budgets (category_id, amount) VALUES
(1, 1500.00), -- Food
(2, 500.00),  -- Transportation
(3, 200.00),  -- Utilities
(4, 300.00),  -- Entertainment
(5, 1000.00); -- Healthcare
