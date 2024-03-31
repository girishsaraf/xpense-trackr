# Use official MySQL Docker image
FROM mysql:latest

# Copy the SQL script to the docker-entrypoint-initdb.d directory
COPY backend/database/init.sql /docker-entrypoint-initdb.d/

# Set environment variables
ENV MYSQL_ROOT_PASSWORD=root_password
ENV MYSQL_DATABASE=xpense_trackr_db
ENV MYSQL_USER=username
ENV MYSQL_PASSWORD=password

# Copy custom MySQL configuration file
#COPY ./my.cnf /etc/mysql/my.cnf

# Expose MySQL port and Go Backend Port
EXPOSE 3306
