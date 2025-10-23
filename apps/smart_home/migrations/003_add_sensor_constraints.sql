-- Migration: 003_add_sensor_constraints.sql
-- Description: Add constraints and validation to sensors table
-- Created: 2024-01-15

-- Add check constraints for status values
ALTER TABLE sensors ADD CONSTRAINT chk_sensors_status 
    CHECK (status IN ('active', 'inactive', 'maintenance', 'error'));

-- Add check constraints for type values
ALTER TABLE sensors ADD CONSTRAINT chk_sensors_type 
    CHECK (type IN ('temperature', 'humidity', 'pressure', 'motion', 'light'));

-- Add check constraints for value range (temperature sensors)
ALTER TABLE sensors ADD CONSTRAINT chk_sensors_value_range 
    CHECK (value >= -50 AND value <= 100);

-- Add unique constraint for name per location
ALTER TABLE sensors ADD CONSTRAINT uq_sensors_name_location 
    UNIQUE (name, location);
