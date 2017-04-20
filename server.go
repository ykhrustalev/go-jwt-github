package exploregithub

import (
	"fmt"
	"github.com/denkyl08/negroniJWT"
	"github.com/go-pg/pg"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"github.com/ykhrustalev/exploregithub/githubapi"
	"github.com/ykhrustalev/exploregithub/handlers"
	"net/http"
)

func OpenDb(conf *DatabaseConfig) *pg.DB {
	return pg.Connect(&pg.Options{
		User:     conf.User,
		Password: conf.Password,
		Database: conf.Database,
		Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		PoolSize: conf.PoolSize,
	})
}

func Server(config *Config) *negroni.Negroni {
	mux := http.NewServeMux()

	// TODO: remove
	negroniJWT.Init(false, "private.key", "public.key")

	githubAuth := githubapi.NewAuth(config.OAuth.ClientId, config.OAuth.ClientSecret)
	db := OpenDb(&config.Database)

	// Serve our event subscription web handler
	mux.HandleFunc("/login", handlers.CreateLoginHandler(db, githubAuth))
	mux.HandleFunc("/logout", handlers.CreateLogoutHandler(db))
	mux.HandleFunc("/auth_callback", handlers.CreateAuthCallbackHandler(db, githubAuth))
	mux.HandleFunc("/github/user", handlers.CreateGithubUserHandler(db, githubAuth))

	n := negroni.New(
		cors.Default(),
		negroni.HandlerFunc(negroniJWT.Middleware),
		negroni.NewRecovery(),
		negroni.NewLogger(),
		// TODO: need?
		negroni.NewStatic(http.Dir("public")),
	)

	n.UseHandler(mux)

	return n
}
