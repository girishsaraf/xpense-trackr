DROP DATABASE xpense_trackr_db;
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
       ('JohnDoe', 'john.doe@example.com', 'hashed_password_1'),
       ('JaneSmith', 'jane.smith@example.com', 'hashed_password_2'),
       ('AliceJohnson', 'alice.johnson@example.com', 'hashed_password_3'),
       ('BobBrown', 'bob.brown@example.com', 'hashed_password_4'),
       ('EmilyJones', 'emily.jones@example.com', 'hashed_password_5'),
       ('MichaelDavis', 'michael.davis@example.com', 'hashed_password_6'),
       ('SophiaWilson', 'sophia.wilson@example.com', 'hashed_password_7'),
       ('DavidMartinez', 'david.martinez@example.com', 'hashed_password_8'),
       ('OliviaTaylor', 'olivia.taylor@example.com', 'hashed_password_9'),
       ('WilliamAnderson', 'william.anderson@example.com', 'hashed_password_10');



INSERT INTO categories (name) VALUES
      ('Food'),
      ('Transportation'),
      ('Utilities'),
      ('Entertainment'),
      ('Healthcare');

INSERT INTO expenses (user_id, category_id, amount, description, date) VALUES
-- Expenses for user1
(1, 1, 50.00, 'Gourmet Coffee and Croissants', '2024-03-01'),
(1, 2, 30.00, 'Tech Gadgets Purchase', '2024-03-02'),
(1, 3, 20.00, 'Art Supplies for Creative Project', '2024-03-03'),
(1, 4, 10.00, 'Outdoor Adventure Gear', '2024-03-04'),
(1, 5, 25.00, 'Bookstore Splurge', '2024-03-05'),

-- Expenses for user2
(2, 1, 40.00, 'Vintage Vinyl Records', '2024-03-01'),
(2, 2, 35.00, 'Fitness Class Membership', '2024-03-02'),
(2, 3, 15.00, 'Sustainable Fashion Finds', '2024-03-03'),
(2, 4, 22.00, 'Home Decor Refresh', '2024-03-04'),
(2, 5, 18.00, 'Culinary Experiment Ingredients', '2024-03-05'),

-- Expenses for user3
(3, 1, 60.00, 'Farmers Market Bounty', '2024-03-01'),
(3, 2, 45.00, 'DIY Home Improvement Supplies', '2024-03-02'),
(3, 3, 12.00, 'Music Festival Tickets', '2024-03-03'),
(3, 4, 28.00, 'Pet Accessories Splurge', '2024-03-04'),
(3, 5, 30.00, 'Local Artisan Crafts', '2024-03-05'),

-- Expenses for user4
(4, 1, 55.00, 'Weekend Getaway Expenses', '2024-03-01'),
(4, 2, 20.00, 'Gardening Supplies', '2024-03-02'),
(4, 3, 28.00, 'Concert Tickets', '2024-03-03'),
(4, 4, 14.00, 'Yoga Retreat Booking', '2024-03-04'),
(4, 5, 37.00, 'Plant-Based Cooking Class', '2024-03-05'),

-- Expenses for user5
(5, 1, 70.00, 'New Running Shoes', '2024-03-01'),
(5, 2, 32.00, 'Tech Conference Pass', '2024-03-02'),
(5, 3, 25.00, 'Picnic Supplies', '2024-03-03'),
(5, 4, 18.00, 'Personal Development Books', '2024-03-04'),
(5, 5, 40.00, 'Beach Day Essentials', '2024-03-05'),

-- Expenses for user6
(6, 1, 45.00, 'Artisanal Cheese Collection', '2024-03-01'),
(6, 2, 27.00, 'Skiing Equipment Rental', '2024-03-02'),
(6, 3, 35.00, 'Pottery Class Registration', '2024-03-03'),
(6, 4, 22.00, 'Healthy Meal Delivery Subscription', '2024-03-04'),
(6, 5, 19.00, 'Vintage Bookstore Finds', '2024-03-05'),

-- Expenses for user7
(7, 1, 52.00, 'Wine Tasting Tour', '2024-03-01'),
(7, 2, 40.00, 'Rock Climbing Gear', '2024-03-02'),
(7, 3, 18.00, 'Coffee Brewing Workshop', '2024-03-03'),
(7, 4, 30.00, 'Home Office Upgrades', '2024-03-04'),
(7, 5, 24.00, 'Hiking Essentials', '2024-03-05'),

-- Expenses for user8
(8, 1, 48.00, 'Cooking Masterclass', '2024-03-01'),
(8, 2, 22.00, 'Home Brewing Kit', '2024-03-02'),
(8, 3, 30.00, 'Gym Membership Renewal', '2024-03-03'),
(8, 4, 35.00, 'DIY Home Renovation Materials', '2024-03-04'),
(8, 5, 21.00, 'Local Farmers Market Splurge', '2024-03-05'),

-- Expenses for user9
(9, 1, 65.00, 'Travel Photography Workshop', '2024-03-01'),
(9, 2, 18.00, 'Plant-Based Cooking Class', '2024-03-02'),
(9, 3, 22.00, 'Mindfulness Meditation Retreat', '2024-03-03'),
(9, 4, 40.00, 'Art Supplies for Creative Project', '2024-03-04'),
(9, 5, 29.00, 'Local Brewery Tour', '2024-03-05'),

-- Expenses for user10
(10, 1, 58.00, 'Vintage Clothing Haul', '2024-03-01'),
(10, 2, 24.00, 'Outdoor Adventure Gear', '2024-03-02'),
(10, 3, 38.00, 'Concert Tickets', '2024-03-03'),
(10, 4, 16.00, 'Cooking Masterclass', '2024-03-04'),
(10, 5, 33.00, 'Gourmet Food Delivery', '2024-03-05');


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
