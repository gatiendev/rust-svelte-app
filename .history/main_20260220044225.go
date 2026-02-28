package main

type Todo struct  {
 isDone bool
 Name   string
 Description string
}

func main() {

  router := gin.default()

  router.GET("/todos", func(c *gin.Context) {
  c.indentedJson(
}
}
