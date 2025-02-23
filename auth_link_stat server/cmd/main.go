package main

import (
	"fmt"
	"net/http"
	"restapi/configs"
	"restapi/internal/auth"
	"restapi/internal/link"
	"restapi/internal/stat"
	"restapi/internal/user"
	"restapi/pkg/db"
	"restapi/pkg/event"
	"restapi/pkg/middleware"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDB(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus() // передаем как зависимость в handle

	// REPOSITORIES
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		StatRepository: statRepository,
		EventBus:       eventBus,
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

	//обработчик подписки ( бесконечно сидит отдельно и ждёт пока не придут сообщения)
	go statService.AddClick()

	//Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	return stack(router)
}

func main() {
	app := App()
	server := http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: app,
	}

	fmt.Println("dd")
	server.ListenAndServe()

}
