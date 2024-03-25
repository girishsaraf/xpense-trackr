-- Create database
CREATE DATABASE xpense_trackr_db;

-- Use the database
USE xpense_trackr_db;

-- Create users table
CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL
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
