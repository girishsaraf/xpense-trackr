-- Create database
CREATE DATABASE xpense_trackr_db;

-- Use the database
USE xpense_trackr_db;

-- Create users table
CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL,
                       password_hash VARCHAR(255) NOT NULL
);


-- Create categories table
CREATE TABLE categories (
                            id INT AUTO_INCREMENT PRIMARY KEY,
                            name VARCHAR(50) UNIQUE NOT NULL
);

-- Create expenses table
CREATE TABLE expenses (
                          id INT AUTO_INCREMENT PRIMARY KEY,
                          user_id INT NOT NULL,
                          category_id INT NOT NULL,
                          amount DECIMAL(10, 2) NOT NULL,
                          description VARCHAR(255) NOT NULL,
                          date DATE NOT NULL,
                          FOREIGN KEY (user_id) REFERENCES users(id),
                          FOREIGN KEY (category_id) REFERENCES categories(id)
);


-- Create budgets table
CREATE TABLE budgets (
                         id INT AUTO_INCREMENT PRIMARY KEY,
                         user_id INT NOT NULL,
                         category_id INT NOT NULL,
                         amount DECIMAL(10, 2) NOT NULL,
                         FOREIGN KEY (user_id) REFERENCES users(id),
                         FOREIGN KEY (category_id) REFERENCES categories(id)
);


INSERT INTO users (username, email, password_hash) VALUES
       ('user1', 'user1@example.com', 'hashed_password_1'),
       ('user2', 'user2@example.com', 'hashed_password_2'),
       ('user3', 'user3@example.com', 'hashed_password_3'),
       ('user4', 'user4@example.com', 'hashed_password_4'),
       ('user5', 'user5@example.com', 'hashed_password_5'),
       ('user6', 'user6@example.com', 'hashed_password_6'),
       ('user7', 'user7@example.com', 'hashed_password_7'),
       ('user8', 'user8@example.com', 'hashed_password_8'),
       ('user9', 'user9@example.com', 'hashed_password_9'),
       ('user10', 'user10@example.com', 'hashed_password_10');


INSERT INTO categories (name) VALUES
      ('Food'),
      ('Transportation'),
      ('Utilities'),
      ('Entertainment'),
      ('Healthcare');

INSERT INTO expenses (user_id, category_id, amount, description, date) VALUES
-- Expenses for user1
(1, 1, 50.00, 'Expense description 1 for user1', '2024-03-01'),
(1, 2, 30.00, 'Expense description 2 for user1', '2024-03-02'),
(1, 3, 20.00, 'Expense description 3 for user1', '2024-03-03'),
(1, 4, 10.00, 'Expense description 4 for user1', '2024-03-04'),
(1, 5, 25.00, 'Expense description 5 for user1', '2024-03-05'),

-- Expenses for user2
(2, 1, 40.00, 'Expense description 1 for user2', '2024-03-01'),
(2, 2, 35.00, 'Expense description 2 for user2', '2024-03-02'),
(2, 3, 15.00, 'Expense description 3 for user2', '2024-03-03'),
(2, 4, 22.00, 'Expense description 4 for user2', '2024-03-04'),
(2, 5, 18.00, 'Expense description 5 for user2', '2024-03-05'),

-- Expenses for user3
(3, 1, 60.00, 'Expense description 1 for user3', '2024-03-01'),
(3, 2, 45.00, 'Expense description 2 for user3', '2024-03-02'),
(3, 3, 12.00, 'Expense description 3 for user3', '2024-03-03'),
(3, 4, 28.00, 'Expense description 4 for user3', '2024-03-04'),
(3, 5, 30.00, 'Expense description 5 for user3', '2024-03-05'),

-- Expenses for user4
(4, 1, 55.00, 'Expense description 1 for user4', '2024-03-01'),
(4, 2, 20.00, 'Expense description 2 for user4', '2024-03-02'),
(4, 3, 28.00, 'Expense description 3 for user4', '2024-03-03'),
(4, 4, 14.00, 'Expense description 4 for user4', '2024-03-04'),
(4, 5, 37.00, 'Expense description 5 for user4', '2024-03-05'),

-- Expenses for user5
(5, 1, 70.00, 'Expense description 1 for user5', '2024-03-01'),
(5, 2, 32.00, 'Expense description 2 for user5', '2024-03-02'),
(5, 3, 25.00, 'Expense description 3 for user5', '2024-03-03'),
(5, 4, 18.00, 'Expense description 4 for user5', '2024-03-04'),
(5, 5, 40.00, 'Expense description 5 for user5', '2024-03-05'),

-- Expenses for user6
(6, 1, 45.00, 'Expense description 1 for user6', '2024-03-01'),
(6, 2, 27.00, 'Expense description 2 for user6', '2024-03-02'),
(6, 3, 35.00, 'Expense description 3 for user6', '2024-03-03'),
(6, 4, 22.00, 'Expense description 4 for user6', '2024-03-04'),
(6, 5, 19.00, 'Expense description 5 for user6', '2024-03-05'),

-- Expenses for user7
(7, 1, 52.00, 'Expense description 1 for user7', '2024-03-01'),
(7, 2, 40.00, 'Expense description 2 for user7', '2024-03-02'),
(7, 3, 18.00, 'Expense description 3 for user7', '2024-03-03'),
(7, 4, 30.00, 'Expense description 4 for user7', '2024-03-04'),
(7, 5, 24.00, 'Expense description 5 for user7', '2024-03-05'),

-- Expenses for user8
(8, 1, 48.00, 'Expense description 1 for user8', '2024-03-01'),
(8, 2, 22.00, 'Expense description 2 for user8', '2024-03-02'),
(8, 3, 30.00, 'Expense description 3 for user8', '2024-03-03'),
(8, 4, 35.00, 'Expense description 4 for user8', '2024-03-04'),
(8, 5, 21.00, 'Expense description 5 for user8', '2024-03-05'),

-- Expenses for user9
(9, 1, 65.00, 'Expense description 1 for user9', '2024-03-01'),
(9, 2, 18.00, 'Expense description 2 for user9', '2024-03-02'),
(9, 3, 22.00, 'Expense description 3 for user9', '2024-03-03'),
(9, 4, 40.00, 'Expense description 4 for user9', '2024-03-04'),
(9, 5, 29.00, 'Expense description 5 for user9', '2024-03-05'),

-- Expenses for user10
(10, 1, 58.00, 'Expense description 1 for user10', '2024-03-01'),
(10, 2, 24.00, 'Expense description 2 for user10', '2024-03-02'),
(10, 3, 38.00, 'Expense description 3 for user10', '2024-03-03'),
(10, 4, 16.00, 'Expense description 4 for user10', '2024-03-04'),
(10, 5, 33.00, 'Expense description 5 for user10', '2024-03-05');

INSERT INTO budgets (user_id, category_id, amount) VALUES
       (1, 1, 500), -- Budget for Food
       (1, 2, 300), -- Budget for Transportation
       (1, 3, 200), -- Budget for Utilities
       (1, 4, 150), -- Budget for Entertainment
       (1, 5, 100); -- Budget for Healthcare

-- Budgets for user2
INSERT INTO budgets (user_id, category_id, amount) VALUES
       (2, 1, 550), -- Budget for Food
       (2, 2, 320), -- Budget for Transportation
       (2, 3, 220), -- Budget for Utilities
       (2, 4, 160), -- Budget for Entertainment
       (2, 5, 110); -- Budget for Healthcare

-- Budgets for user3
INSERT INTO budgets (user_id, category_id, amount) VALUES
       (3, 1, 600), -- Budget for Food
       (3, 2, 350), -- Budget for Transportation
       (3, 3, 240), -- Budget for Utilities
       (3, 4, 180), -- Budget for Entertainment
       (3, 5, 120); -- Budget for Healthcare

-- Budgets for user4
INSERT INTO budgets (user_id, category_id, amount) VALUES
    (4, 1, 530), -- Budget for Food
    (4, 2, 290), -- Budget for Transportation
    (4, 3, 210), -- Budget for Utilities
    (4, 4, 170), -- Budget for Entertainment
    (4, 5, 130); -- Budget for Healthcare

-- Budgets for user5
INSERT INTO budgets (user_id, category_id, amount) VALUES
   (5, 1, 520), -- Budget for Food
   (5, 2, 280), -- Budget for Transportation
   (5, 3, 190), -- Budget for Utilities
   (5, 4, 140), -- Budget for Entertainment
   (5, 5, 150); -- Budget for Healthcare

-- Budgets for user6
INSERT INTO budgets (user_id, category_id, amount) VALUES
   (6, 1, 570), -- Budget for Food
   (6, 2, 310), -- Budget for Transportation
   (6, 3, 230), -- Budget for Utilities
   (6, 4, 155), -- Budget for Entertainment
   (6, 5, 105); -- Budget for Healthcare

-- Budgets for user7
INSERT INTO budgets (user_id, category_id, amount) VALUES
   (7, 1, 580), -- Budget for Food
   (7, 2, 330), -- Budget for Transportation
   (7, 3, 200), -- Budget for Utilities
   (7, 4, 160), -- Budget for Entertainment
   (7, 5, 100); -- Budget for Healthcare

-- Budgets for user8
INSERT INTO budgets (user_id, category_id, amount) VALUES
       (8, 1, 510), -- Budget for Food
       (8, 2, 270), -- Budget for Transportation
       (8, 3, 220), -- Budget for Utilities
       (8, 4, 175), -- Budget for Entertainment
       (8, 5, 120); -- Budget for Healthcare

-- Budgets for user9
INSERT INTO budgets (user_id, category_id, amount) VALUES
       (9, 1, 540), -- Budget for Food
       (9, 2, 320), -- Budget for Transportation
       (9, 3, 210), -- Budget for Utilities
       (9, 4, 150), -- Budget for Entertainment
       (9, 5, 130); -- Budget for Healthcare

-- Budgets for user10
INSERT INTO budgets (user_id, category_id, amount) VALUES
   (10, 1, 590), -- Budget for Food
   (10, 2, 350), -- Budget for Transportation
   (10, 3, 240), -- Budget for Utilities
   (10, 4, 170), -- Budget for Entertainment
   (10, 5, 110); -- Budget for Healthcare
