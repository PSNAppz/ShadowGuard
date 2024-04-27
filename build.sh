#! /bin/bash

psql -U postgres -f sql/initialize_database.sql
psql -U postgres -d gorm -c "GRANT USAGE ON SCHEMA public TO gorm;"