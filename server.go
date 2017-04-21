package exploregithub

import (
	"fmt"
	"github.com/denkyl08/negroniJWT"
	"github.com/go-pg/pg"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"github.com/ykhrustalev/exploregithub/githubapi"
	"github.com/ykhrustalev/exploregithub/handlers"
	"github.com/ykhrustalev/exploregithub/jwtutils"
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
	// JWT auto settings
	negroniJWT.Init(false, "private.key", "public.key")

	githubAuth := githubapi.NewAuth(config.OAuth.ClientId, config.OAuth.ClientSecret)

	db := OpenDb(&config.Database)

	jwtAuthHandler := negroni.HandlerFunc(jwtutils.CreateAuthMiddleware(db))

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handlers.CreateLoginHandler(db, githubAuth))
	mux.HandleFunc("/auth_callback", handlers.CreateAuthCallbackHandler(db, githubAuth))

	mux.Handle("/github/user", negroni.New(
		jwtAuthHandler,
		negroni.Wrap(handlers.CreateGithubUserHandler(db, githubAuth)),
	))
	mux.Handle("/logout", negroni.New(
		jwtAuthHandler,
		negroni.Wrap(handlers.CreateLogoutHandler(db)),
	))

	n := negroni.New(
		cors.Default(),
		negroni.HandlerFunc(negroniJWT.Middleware),
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)

	n.UseHandler(mux)

	return n
}
