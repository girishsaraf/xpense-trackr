# Use official MySQL Docker image
FROM mysql:latest

# Set environment variables
ENV MYSQL_ROOT_PASSWORD=root_password
ENV MYSQL_DATABASE=xpense-trackr
ENV MYSQL_USER=username
ENV MYSQL_PASSWORD=password

# Copy custom MySQL configuration file
#COPY ./my.cnf /etc/mysql/my.cnf

# Expose MySQL port
EXPOSE 3306
