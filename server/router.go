package server

import (
	"gornyakWarehouse/internal/database"

	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

var indexTmpl = `
<html>
  <h1>Welcome</h1>
  <ul>
    <li><a href="/api/users/123">/api/users/123</a></li>
    <li><a href="/api/users/current">/api/users/current</a></li>
    <li><a href="/api/users/foo/bar">/api/users/foo/bar</a></li>
    <li><a href="/404">/404</a></li>
    <li><a href="/405">/405</a></li>
  <li><a href="/api/users/all">/api/users/all</a></li>
  </ul>
</html>
`

func setRouter() *bunrouter.Router {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware(
			reqlog.FromEnv("BUNDEBUG"),
		)),
	)

	router.GET("/", indexHandler)
	router.POST("/405", indexHandler)

	router.WithGroup("/api", func (g *bunrouter.Group) {
		g.GET("/users/:id", debugHandler)
		g.GET("/users/current", debugHandler)
		g.GET("/users/*path", debugHandler)
		g.GET("/users/all", getAllUsers)
	})

	log.Println("Listening on http://localhost:9999")
	log.Println(http.ListenAndServe(":9999", router))

	return router
}

func debugHandler(w http.ResponseWriter, req bunrouter.Request) error {
	return bunrouter.JSON(w, bunrouter.H{
		"route": req.Route(),
		"params": req.Params().Map(),
	})
}

func indexHandler(w http.ResponseWriter, req bunrouter.Request) error {
	return indexTemplate().Execute(w, nil)
}

func indexTemplate() *template.Template {
	return template.Must(template.New("index").Parse(indexTmpl))
}

func getAllUsers(w http.ResponseWriter, req bunrouter.Request) error {
	db := database.NewDBConnection()
	ctx := context.Background()

	users, err := database.GetAllUsers(ctx, db)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, bunrouter.H{
		"result": users,
	})
}