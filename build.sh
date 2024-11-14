#!/bin/bash

# Variables
DB_USER="gorm"      # replace with the desired PostgreSQL username
DB_PASSWORD="gorm"  # replace with the desired password
DB_NAME="gorm"  # replace with the desired database name

# Check if PostgreSQL is installed; if not, install it
if ! command -v psql &> /dev/null; then
    echo "PostgreSQL is not installed. Installing..."
    sudo apt update
    sudo apt install -y postgresql postgresql-contrib
else
    echo "PostgreSQL is already installed."
fi

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Switch to the postgres user and run SQL commands to create user and database
sudo -u postgres psql <<EOF
-- Create a new PostgreSQL user with password
DO \$$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_roles
        WHERE rolname = '${DB_USER}'
    ) THEN
        CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';
    END IF;
END
\$$;

-- Create a new PostgreSQL database owned by the new user
DO \$$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_database
        WHERE datname = '${DB_NAME}'
    ) THEN
        CREATE DATABASE ${DB_NAME} OWNER ${DB_USER};
    END IF;
END
\$$;
EOF

# Variables
PG_HBA_PATH="/var/lib/pgsql/data/pg_hba.conf"
BACKUP_PATH="/var/lib/pgsql/backups/pg_hba.conf.bak"

# Backup the original pg_hba.conf file
if [ ! -f "$BACKUP_PATH" ]; then
    echo "Backing up the original pg_hba.conf to $BACKUP_PATH"
    sudo cp "$PG_HBA_PATH" "$BACKUP_PATH"
else
    echo "Backup already exists at $BACKUP_PATH"
fi

# Update authentication method to md5
echo "Updating authentication method to md5 in $PG_HBA_PATH"
sudo sed -i 's/^\(local\s\+all\s\+all\s\+\)peer/\1md5/' "$PG_HBA_PATH"
sudo sed -i 's/^\(host\s\+all\s\+all\s\+127.0.0.1\/32\s\+\)ident/\1md5/' "$PG_HBA_PATH"
sudo sed -i 's/^\(host\s\+all\s\+all\s\+::1\/128\s\+\)ident/\1md5/' "$PG_HBA_PATH"

# Restart PostgreSQL to apply changes
echo "Restarting PostgreSQL service to apply changes..."
sudo systemctl restart postgresql

echo "Authentication method updated to md5. Please test your connection."

echo "User and database setup complete."
