package main

type Todo struct  {
 isDone bool
 Name   string
 Description string
}

func main() {

  router := gin.default()

  router.GET("/todos", func(c *gin.Context) {
  c.indentedJSON(http.StatusOK, []Todo {
	{
		Name: "Do the laundry"
	},
	{
		Name: "Clean the dishes"
	},

  })
})

router.Run("localhost:8080")
}
