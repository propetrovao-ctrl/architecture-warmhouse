package models

import (
	"time"
)

// SensorType represents the type of sensor
type SensorType string

const (
	Temperature SensorType = "temperature"
	Humidity    SensorType = "humidity"
	Pressure    SensorType = "pressure"
	Motion      SensorType = "motion"
	Light       SensorType = "light"
)

// Sensor represents a smart home sensor
type Sensor struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Type        SensorType `json:"type"`
	Location    string     `json:"location"`
	Value       float64    `json:"value"`
	Unit        string     `json:"unit"`
	Status      string     `json:"status"`
	LastUpdated time.Time  `json:"last_updated"`
	CreatedAt   time.Time  `json:"created_at"`
}

// SensorCreate represents the data needed to create a new sensor
type SensorCreate struct {
	Name     string     `json:"name" binding:"required,min=1,max=100"`
	Type     SensorType `json:"type" binding:"required,oneof=temperature humidity pressure motion light"`
	Location string     `json:"location" binding:"required,min=1,max=100"`
	Unit     string     `json:"unit" binding:"omitempty,max=20"`
}

// SensorUpdate represents the data that can be updated for a sensor
type SensorUpdate struct {
	Name     string     `json:"name"`
	Type     SensorType `json:"type"`
	Location string     `json:"location"`
	Value    *float64   `json:"value"`
	Unit     string     `json:"unit"`
	Status   string     `json:"status"`
}
