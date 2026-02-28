package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Todo struct {
    IsDone      bool   // exported field (capitalized) for JSON
    Name        string
    Description string
}

func main() {
    router := gin.Default() // note capital D

    router.GET("/todos", func(c *gin.Context) {
        c.IndentedJSON(http.StatusOK, []Todo{ // capital I, and import http
            {
                Name: "Do the laundry", // comma after field
            },
            {
                Name: "Clean the dishes", // comma after field
            },
        })
    })

    router.Run("localhost:8080")
}