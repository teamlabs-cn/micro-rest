# micro-rest
resource-based, mini framework for web api

# Install

`go get github.com/teamlabs-cn/micro-rest/rest`

# Samples

### Definie resource

```go
func User() rest.HttpResource {
	res := rest.NewResource()

	res.Get("{id}", func(ctx *rest.HttpContext) {
		bytes, _ := json.Marshal(struct {
			Id   int
			Name string
		}{100000, "mixlatte"})

		resp.Json(ctx.Response, bytes)
	})

	res.Get("{id}/managers/{manager_id}", func(ctx *rest.HttpContext) {
		fmt.Println(ctx.RouteData)
		bytes, _ := json.Marshal(struct {
			Id   int
			Name string
		}{100111, "tomy lu"})

		resp.Json(ctx.Response, bytes)
	})

	res.Post("test/{id}", func(ctx *rest.HttpContext) {
		fmt.Println("POST /users/:id")
	})

	return res
}
```

## Build rest app

```go
func main() {
	PORT := ":8080"
	log.Print("Running server on " + PORT)
	
	// global middleware
	rest.Use(middlewares.Logging())
	// global resource
	rest.Resource("/files", resources.File)

	app := rest.NewApp("/userframework/v1")
	// app middleware
	app.Use(middlewares.Auth("123456"))
	// app resource
	app.Resource("/users", resources.User)

	log.Fatal(http.ListenAndServe(PORT, nil))
}
```
