package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID          int    `json:"id"`          // added ID field
	IsDone      bool   `json:"is_done"`     // json tag for snake_case
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage (simple slice)
var todos = []Todo{}
var nextID = 1 // simple auto-increment

func main() {
	router := gin.Default()

	// Routes
	router.GET("/todos", getTodos)
	router.POST("/todos", createTodo)
	router.GET("/todos/:id", getTodoByID)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)

	router.Run("localhost:8080")
}

// GET /todos – return all todos
func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

// POST /todos – add a new todo
func createTodo(c *gin.Context) {
	var newTodo Todo

	// Bind JSON to newTodo
	if err := c.BindJSON(&newTodo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Assign ID and append
	newTodo.ID = nextID
	nextID++
	todos = append(todos, newTodo)

	c.IndentedJSON(http.StatusCreated, newTodo)
}

// GET /todos/:id – get one todo by ID
func getTodoByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for _, todo := range todos {
		if todo.ID == id {
			c.IndentedJSON(http.StatusOK, todo)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

// PUT /todos/:id – update an existing todo
func updateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedTodo Todo
	if err := c.BindJSON(&updatedTodo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			// Preserve ID, update other fields
			updatedTodo.ID = id
			todos[i] = updatedTodo
			c.IndentedJSON(http.StatusOK, updatedTodo)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

// DELETE /todos/:id – remove a todo
func deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			// Remove element from slice
			todos = append(todos[:i], todos[i+1:]...)
			c.IndentedJSON(http.StatusNoContent, nil) // 204 No Content
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}