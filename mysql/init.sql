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
(1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Groceries', DATE_SUB(CURRENT_DATE(), INTERVAL 0 MONTH)),
(1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Dining out', DATE_SUB(CURRENT_DATE(), INTERVAL 1 MONTH)),
(1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Food delivery', DATE_SUB(CURRENT_DATE(), INTERVAL 2 MONTH)),
(1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Home cooking', DATE_SUB(CURRENT_DATE(), INTERVAL 3 MONTH)),
(1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Snacks', DATE_SUB(CURRENT_DATE(), INTERVAL 4 MONTH)),
(1, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Restaurant', DATE_SUB(CURRENT_DATE(), INTERVAL 5 MONTH));

-- Insert query for 'Transportation'
INSERT INTO expenses (category_id, amount, description, date) VALUES
(2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Gasoline', DATE_SUB(CURRENT_DATE(), INTERVAL 0 MONTH)),
(2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Public transportation', DATE_SUB(CURRENT_DATE(), INTERVAL 1 MONTH)),
(2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Taxi', DATE_SUB(CURRENT_DATE(), INTERVAL 2 MONTH)),
(2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Ride-sharing', DATE_SUB(CURRENT_DATE(), INTERVAL 3 MONTH)),
(2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Car maintenance', DATE_SUB(CURRENT_DATE(), INTERVAL 4 MONTH)),
(2, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Parking fees', DATE_SUB(CURRENT_DATE(), INTERVAL 5 MONTH));

-- Insert query for 'Utilities'
INSERT INTO expenses (category_id, amount, description, date) VALUES
(3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Electricity', DATE_SUB(CURRENT_DATE(), INTERVAL 0 MONTH)),
(3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Water', DATE_SUB(CURRENT_DATE(), INTERVAL 1 MONTH)),
(3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Internet', DATE_SUB(CURRENT_DATE(), INTERVAL 2 MONTH)),
(3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Phone bill', DATE_SUB(CURRENT_DATE(), INTERVAL 3 MONTH)),
(3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Gas', DATE_SUB(CURRENT_DATE(), INTERVAL 4 MONTH)),
(3, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Trash collection', DATE_SUB(CURRENT_DATE(), INTERVAL 5 MONTH));

-- Insert query for 'Entertainment'
INSERT INTO expenses (category_id, amount, description, date) VALUES
(4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Movie tickets', DATE_SUB(CURRENT_DATE(), INTERVAL 0 MONTH)),
(4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Concert', DATE_SUB(CURRENT_DATE(), INTERVAL 1 MONTH)),
(4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Video games', DATE_SUB(CURRENT_DATE(), INTERVAL 2 MONTH)),
(4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Books', DATE_SUB(CURRENT_DATE(), INTERVAL 3 MONTH)),
(4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Streaming service', DATE_SUB(CURRENT_DATE(), INTERVAL 4 MONTH)),
(4, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Board games', DATE_SUB(CURRENT_DATE(), INTERVAL 5 MONTH));

-- Insert query for 'Healthcare'
INSERT INTO expenses (category_id, amount, description, date) VALUES
(5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Medications', DATE_SUB(CURRENT_DATE(), INTERVAL 0 MONTH)),
(5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Doctor visit', DATE_SUB(CURRENT_DATE(), INTERVAL 1 MONTH)),
(5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Health insurance', DATE_SUB(CURRENT_DATE(), INTERVAL 2 MONTH)),
(5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Dental care', DATE_SUB(CURRENT_DATE(), INTERVAL 3 MONTH)),
(5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Gym membership', DATE_SUB(CURRENT_DATE(), INTERVAL 4 MONTH)),
(5, FLOOR(RAND() * (250 - 10 + 1)) + 10, 'Prescription glasses', DATE_SUB(CURRENT_DATE(), INTERVAL 5 MONTH));



INSERT INTO budgets (category_id, amount) VALUES
(1, 1500.00), -- Food
(2, 500.00),  -- Transportation
(3, 200.00),  -- Utilities
(4, 300.00),  -- Entertainment
(5, 1000.00); -- Healthcare


-- Recurring expenses and investments

CREATE TABLE recurring_expense (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    frequency VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    category_id INT NOT NULL CHECK (category_id BETWEEN 1 AND 5) -- Add category_id column
);

CREATE TABLE recurring_investment (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    frequency VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    type ENUM('stocks', '401k', 'mutual_funds', 'real_estate', 'savings', 'bonds') NOT NULL -- Add type column
);

INSERT INTO recurring_expense (description, amount, frequency, start_date, end_date, category_id)
VALUES
    ('Internet Subscription', 49.99, 'monthly', CURRENT_DATE(), NULL, 1),
    ('Gym Membership', 29.99, 'monthly', CURRENT_DATE(), NULL, 2),
    ('Mobile Phone Plan', 59.99, 'monthly', CURRENT_DATE(), NULL, 3),
    ('Scoop to Work', 6.99, 'daily', CURRENT_DATE(), NULL, 3);

INSERT INTO recurring_investment (description, amount, frequency, start_date, end_date, type)
VALUES
    ('401(k) Contribution', 500.00, 'monthly', CURRENT_DATE(), NULL, '401k'),
    ('Mutual Funds Investment', 1000.00, 'monthly', CURRENT_DATE(), NULL, 'mutual_funds'),
    ('Stocks Investment', 1500.00, 'monthly', CURRENT_DATE(), NULL, 'stocks');



