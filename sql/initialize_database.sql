-- Create the "gorm" user with the password "gorm"
CREATE USER gorm WITH PASSWORD 'gorm';

-- Create the "gorm" database and grant privileges to the "gorm" user
CREATE DATABASE gorm;
GRANT ALL PRIVILEGES ON DATABASE gorm TO gorm;