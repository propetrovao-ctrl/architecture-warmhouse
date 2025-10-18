package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// TemperatureResponse represents the response structure
type TemperatureResponse struct {
	SensorID    string    `json:"sensorId"`
	SensorType  string    `json:"sensorType"`
	Location    string    `json:"location"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Temperature API is running!",
			"version": "1.0.0",
		})
	})

	r.GET("/temperature", getTemperature)
	r.GET("/temperature/:sensorId", getTemperatureBySensorID)
	r.GET("/temperature/health", healthCheck)

	// Start server on port 8081
	r.Run(":8081")
}

// getTemperature handles GET /temperature requests
func getTemperature(c *gin.Context) {
	location := c.Query("location")
	sensorID := c.Query("sensorId")

	// If no location is provided, use a default based on sensor ID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// If no sensor ID is provided, generate one based on location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	// Generate random temperature between 15 and 30 degrees
	temperature := 15 + rand.Float64()*15
	temperature = float64(int(temperature*10)) / 10 // Round to 1 decimal place

	// Determine status based on temperature
	var status string
	if temperature < 18 {
		status = "Cold"
	} else if temperature < 25 {
		status = "Comfortable"
	} else {
		status = "Hot"
	}

	response := TemperatureResponse{
		SensorID:    sensorID,
		SensorType:  "temperature",
		Location:    location,
		Value:       temperature,
		Unit:        "°C",
		Status:      status,
		Timestamp:   time.Now().UTC(),
		Description: "Temperature in " + location + ": " + strconv.FormatFloat(temperature, 'f', 1, 64) + "°C (" + status + ")",
	}

	c.JSON(http.StatusOK, response)
}

// getTemperatureBySensorID handles GET /temperature/{sensorId} requests
func getTemperatureBySensorID(c *gin.Context) {
	sensorID := c.Param("sensorId")

	// Determine location based on sensor ID
	var location string
	switch sensorID {
	case "1":
		location = "Living Room"
	case "2":
		location = "Bedroom"
	case "3":
		location = "Kitchen"
	default:
		location = "Unknown"
	}

	// Generate random temperature between 15 and 30 degrees
	temperature := 15 + rand.Float64()*15
	temperature = float64(int(temperature*10)) / 10 // Round to 1 decimal place

	// Determine status based on temperature
	var status string
	if temperature < 18 {
		status = "Cold"
	} else if temperature < 25 {
		status = "Comfortable"
	} else {
		status = "Hot"
	}

	response := TemperatureResponse{
		SensorID:    sensorID,
		SensorType:  "temperature",
		Location:    location,
		Value:       temperature,
		Unit:        "°C",
		Status:      status,
		Timestamp:   time.Now().UTC(),
		Description: "Temperature in " + location + " (sensor " + sensorID + "): " + strconv.FormatFloat(temperature, 'f', 1, 64) + "°C (" + status + ")",
	}

	c.JSON(http.StatusOK, response)
}

// healthCheck handles GET /temperature/health requests
func healthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Message:   "Temperature API is running!",
		Timestamp: time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

