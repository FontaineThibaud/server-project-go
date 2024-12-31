package main

import (
	"api/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		log.Println("Ping endpoint called")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/health", func(c *gin.Context) {
		log.Println("Health check endpoint called")
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API is running",
		})
	})

	r.POST("/api/users", createUser)
	r.GET("/api/users", getUsers)
	r.GET("/api/users/:id", getUser)

	log.Printf("Server starting on port 8081")
	r.Run(":8081")
}

func createUser(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.QueryRow(
		"INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id",
		user.Username, user.Email,
	)

	var id int
	if err := result.Scan(&id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       id,
		"username": user.Username,
		"email":    user.Email,
	})
}

func getUsers(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, username, email, created_at FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []gin.H
	for rows.Next() {
		var id int
		var username, email, createdAt string
		if err := rows.Scan(&id, &username, &email, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, gin.H{
			"id":         id,
			"username":   username,
			"email":      email,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	var username, email, createdAt string
	err := database.DB.QueryRow(
		"SELECT username, email, created_at FROM users WHERE id = $1",
		id,
	).Scan(&username, &email, &createdAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         id,
		"username":   username,
		"email":      email,
		"created_at": createdAt,
	})
}
