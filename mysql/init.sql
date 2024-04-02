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
INSERT INTO investments (amount, investment_type, date) VALUES
(1000.00, 'stocks', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(500.50, '401k', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1500.25, 'mutual_funds', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2000.75, 'real_estate', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3000.00, 'savings', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(750.20, 'bonds', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1200.00, 'stocks', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(800.75, '401k', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2500.50, 'mutual_funds', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1800.25, 'real_estate', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5000.00, 'savings', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1000.20, 'bonds', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1500.00, 'stocks', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(600.75, '401k', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2000.50, 'mutual_funds', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3500.25, 'real_estate', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4000.00, 'savings', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(950.20, 'bonds', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1800.00, 'stocks', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(700.75, '401k', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3000.50, 'mutual_funds', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2200.25, 'real_estate', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(6000.00, 'savings', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(850.20, 'bonds', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY));

INSERT INTO categories (name) VALUES
('Food'),
('Transportation'),
('Utilities'),
('Entertainment'),
('Healthcare');

-- Expenses for the past 6 months
INSERT INTO expenses (category_id, amount, description, date) VALUES
(1, 50.00, 'Farmers Market Goodies', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 30.00, 'City Bus Commute', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 20.00, 'Energy Saver Gadgets', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 10.00, 'Blockbuster Movie Night', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 25.00, 'Holistic Healing Session', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1, 40.00, 'Epicurean Adventure', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 35.00, 'Ride-Hailing Journey', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 15.00, 'Hydroelectricity Bill', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 22.00, 'Concert Under the Stars', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 18.00, 'Pharmacy Essentials', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1, 60.00, 'Fresh Produce Haul', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 45.00, 'Train Adventure', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 12.00, 'Gas Utility Bill', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 28.00, 'Theme Park Escapade', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 30.00, 'Oral Hygiene Visit', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1, 55.00, 'Health Food Splurge', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 20.00, 'Public Transit Fare', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 28.00, 'Electricity Renewal', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 14.00, 'Big Screen Thrills', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 37.00, 'Medical Check-Up', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1, 70.00, 'Fine Dining Experience', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 32.00, 'Taxi Joyride', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 25.00, 'Water Utility Bill', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 18.00, 'Live Concert Extravaganza', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 40.00, 'Pharmaceutical Supplies', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1, 45.00, 'Organic Food Selection', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 27.00, 'City Train Journey', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 35.00, 'Gasoline Refill', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 22.00, 'Theme Park Fun', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 19.00, 'Dental Appointment', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1, 52.00, 'Fresh Market Finds', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 40.00, 'Public Bus Ride', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 18.00, 'Power Utility Bill', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 30.00, 'Blockbuster Cinema Experience', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 24.00, 'Health Check-Up', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1, 48.00, 'Farm-to-Table Feast', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 22.00, 'Urban Taxi Ride', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 30.00, 'Water Supply Bill', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 35.00, 'Live Music Concert Tickets', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 21.00, 'Medicinal Prescriptions', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1, 65.00, 'Grocery Shopping Extravaganza', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 18.00, 'Metro Train Ride', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 22.00, 'Utility Gas Bill', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 40.00, 'Theme Park Passes', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 29.00, 'Teeth Cleaning Session', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(1, 58.00, 'Healthy Grocery Haul', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(2, 24.00, 'City Bus Pass', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(3, 38.00, 'Electricity Service Fee', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(4, 16.00, 'Cinematic Experience', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY)),
(5, 33.00, 'Physician Appointment', DATE_SUB(CURRENT_DATE(), INTERVAL FLOOR(RAND() * 180) DAY));

INSERT INTO budgets (category_id, amount) VALUES
(1, 1500.00), -- Food
(2, 500.00),  -- Transportation
(3, 200.00),  -- Utilities
(4, 300.00),  -- Entertainment
(5, 1000.00); -- Healthcare
