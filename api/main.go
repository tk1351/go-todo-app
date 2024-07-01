package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Todo struct {
	*gorm.Model
	Content string `json:"content"`
}

func GetAll(c *gin.Context) {
	var todos []Todo
	result := db.Find(&todos)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "all todos",
		"data":    todos,
	})
}

func Create(c *gin.Context) {
	var data Todo

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := Todo{Content: data.Content}
	result := db.Create(&todo)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create new",
		"data":    todo,
	})
}

func Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var payload Todo

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var todo Todo
	result := db.First(&todo, id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	todo.Content = payload.Content
	db.Save(&todo)

	c.JSON(http.StatusOK, gin.H{
		"message": "update todo",
		"data":    todo,
	})
}

func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var todo Todo
	result := db.First(&todo, id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	result = db.Delete(&todo, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted todo",
		"data":    todo,
	})
}

var db *gorm.DB

func main() {
	loadEnv()
	r := gin.Default()

	var err error

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	failOnError(err, "Error DB Connection")

	db.AutoMigrate(&Todo{})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/todos", GetAll)

	r.POST("/todos", Create)

	r.PUT("/todos/:id", Update)

	r.DELETE("/todos/:id", Delete)

	fmt.Println("Database connection and setup successful")
	r.Run()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Can not load: %v", err)
	}

	message := os.Getenv("POSTGRES_HOST")

	fmt.Println(message)
}
