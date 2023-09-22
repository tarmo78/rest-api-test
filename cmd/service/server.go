package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/cors"

	"github.com/tarmo78/rest-api-test/cmd/service/api"
)

type (
	serverContext struct {
		listenPort string
		router     *mux.Router

		controller *api.Controller // necessary?
	}

	apiStarter struct {
		serverCtx *serverContext
		config    *serverConfig
	}
)

func main() {

	//client := twilio.NewRestClient()

	//sendMsg(client)

	config := &serverConfig{
		Port: "8080",
	}
	apiServer := &apiStarter{
		config: config,
		serverCtx: &serverContext{
			listenPort: config.Port,
		},
	}

	httpServer, err := apiServer.serverCtx.createServer()
	if err != nil {
		fmt.Printf("Failed to create/init api server: %s", err)
		os.Exit(1)
	}

	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("HTTP server exited with error: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func (s *serverContext) createServer() (*http.Server, error) {
	if err := s.initController(); err != nil {
		return nil, errors.Wrap(err, "failed to init controller")
	}

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT"},
		//AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	server := &http.Server{
		Addr:    ":" + s.listenPort,
		Handler: corsWrapper.Handler(s.router),
	}

	return server, nil
}

func (s *serverContext) initController() error {
	s.router = mux.NewRouter()

	controller, err := api.NewController()
	if err != nil {
		return errors.Wrap(err, "failed to create controller")
	}
	s.controller = controller
	s.controller.SetupRouter(s.router)

	return nil
}
