// @title LinkShorty API
// @version 1.0
// @description Сервис для управления короткими ссылками
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"LinkShorty/configs"
	"LinkShorty/internal/auth"
	"LinkShorty/internal/link"
	"LinkShorty/internal/stat"
	"LinkShorty/internal/user"
	"LinkShorty/pkg/db"
	"LinkShorty/pkg/event"
	"LinkShorty/pkg/middleware"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	database := db.NewDB(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Swagger
	router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Убедитесь, что этот путь правильный
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.PersistAuthorization(false),
	))

	// Обработчик для swagger/doc.json
	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	//Repository
	linkRepository := link.NewLinkRepository(database)
	userRepository := user.NewUserRepository(database)
	statRepository := stat.NewStatRepository(database)

	//Service
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	//Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		EventBus:       eventBus,
		Config:         conf,
	})

	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	go statService.AddClick()

	//Middleware
	stack := middleware.Chain(
		middleware.Cors,
		middleware.Logging,
	)
	return stack(router)
}

func main() {
	app := App()

	server := &http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("server is listening on :8081")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
