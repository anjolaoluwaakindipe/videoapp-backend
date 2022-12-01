package server

import (
	"context"
	"net/http"

	"github.com/anjolaoluwaakindipe/videoapp/routes"
	"github.com/anjolaoluwaakindipe/videoapp/utils/config"
	"github.com/anjolaoluwaakindipe/videoapp/utils/logger"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// injected param struct from fx
type InitServerParams struct {
	fx.In
	AppConfig config.AppConfig
	Logger    logger.Logger
	Routes    []routes.Route `group:"routes"`
}

func InitServer(lc fx.Lifecycle, isp InitServerParams) *http.Server {
	// create mutiplexer from gorilla mux
	gorrillaMux := mux.NewRouter()

	// get all routes from injected params
	routes := isp.Routes

	isp.Logger.Infoln("Registering routes...")

	// Hanling all routes with there respective hanlers
	for i :=0 ; i < len(routes); i++ {
		subRoutes := routes[i].Routes()

		for j := 0 ; j < len(subRoutes) ; j++ {
			gorrillaMux.HandleFunc(subRoutes[j].Path, subRoutes[j].HandlerFunc).Methods(subRoutes[j].Methods...)
		}
	}

	// global middleware
	gorrillaMux.Use(mux.CORSMethodMiddleware(gorrillaMux))

	// create http server
	httpServer := &http.Server{
		Addr:    "localhost:" + isp.AppConfig.Port,
		Handler: gorrillaMux,
	}

	// fx lifecycle methods
	lc.Append(fx.Hook{

		OnStart: func(ctx context.Context) error {
			
			// start server  
			go func() {

				isp.Logger.Infof("Server Starting on port %s... \n", isp.AppConfig.Port)
				if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
					isp.Logger.Fatal("Could not start server: " + err.Error())
				}

			}()
			return nil
		},

		OnStop: func(ctx context.Context) error {
			// gracefully shutdown server with context
			if err := httpServer.Shutdown(ctx); err != nil {
				isp.Logger.Fatalln("Error occured while gracefully shutdown server!!! " + err.Error())
				return err
			}
			isp.Logger.Infoln("Serving Gracefully terminated...")
			return nil
		},
	})

	return httpServer
}

var Module = fx.Module("server", fx.Provide(InitServer), fx.Invoke(func(*http.Server) {}))
