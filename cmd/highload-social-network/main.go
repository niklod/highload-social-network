package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/streadway/amqp"

	"github.com/niklod/highload-social-network/config"
	"github.com/niklod/highload-social-network/internal/cache"
	"github.com/niklod/highload-social-network/internal/queue/feed"
	"github.com/niklod/highload-social-network/internal/queue/feed/producer"
	"github.com/niklod/highload-social-network/internal/queue/feed/receiver"
	"github.com/niklod/highload-social-network/internal/server"
	"github.com/niklod/highload-social-network/internal/user"
	"github.com/niklod/highload-social-network/internal/user/city"
	"github.com/niklod/highload-social-network/internal/user/interest"
	"github.com/niklod/highload-social-network/internal/user/post"
	"github.com/niklod/highload-social-network/internal/websocket"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Connections
	db, err := dbConnect("mysql", cfg.DB.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	conn, err := amqp.Dial(cfg.RabbitMQ.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	userRepo := user.NewRepository(db)
	cityRepo := city.NewRepository(db)
	interestRepo := interest.NewRepository(db)
	postRepo := post.NewRepository(db)

	feedCache := cache.NewFeedCache()
	ch, err := feed.NewQueueChannel(conn, cfg.RabbitMQ)
	if err != nil {
		log.Fatal(err)
	}

	// WebSockets pool
	wsPool := websocket.NewPool()
	go wsPool.Start()

	// Services
	cityService := city.NewService(cityRepo)
	interestService := interest.NewService(interestRepo)
	userService := user.NewService(userRepo, cityService, interestService)
	feedProducer := producer.NewFeedProducer(ch, cfg.RabbitMQ, feedCache)
	postService := post.NewService(postRepo, feedCache, feedProducer)
	feedReceiver := receiver.NewFeedReceiver(ch, cfg.RabbitMQ, feedCache, postService, userService, wsPool)

	// Starting feed update receivers
	for i := 0; i < cfg.RabbitMQ.ReceiversCount; i++ {
		go feedReceiver.Run()
	}

	cookieStore := sessions.NewCookieStore([]byte(cfg.SecretKey))
	gob.Register(user.User{})

	// Handlers
	userHandler := user.NewHandler(
		userService,
		cityService,
		postService,
		cookieStore,
		interestService,
	)
	wsHandler := websocket.NewWebsocketHandler(wsPool, userService)

	srv := server.NewHTTPServer(cfg.Server)
	srv.BaseRouterGroup.Use(userHandler.AuthMiddleware)

	// Главная
	srv.BaseRouterGroup.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/login")
	})

	// Регистрациия
	srv.BaseRouterGroup.GET("/registrate", userHandler.HandleUserRegistrate)
	srv.BaseRouterGroup.POST("/registrate", userHandler.HandleUserRegistrateSubmit)

	// Вход Выход
	srv.BaseRouterGroup.GET("/login", userHandler.HandleUserLogin)
	srv.BaseRouterGroup.POST("/login", userHandler.HandleUserLoginSubmit)
	srv.BaseRouterGroup.GET("/logout", userHandler.HandleUserLogout)

	// User detail page
	srv.BaseRouterGroup.GET("/user/:login", userHandler.HandleUserDetail)

	// Добавление Удаление из друзей
	srv.BaseRouterGroup.POST("/user/:login/add_friend", userHandler.HandleAddFriend)
	srv.BaseRouterGroup.POST("/user/:login/delete_friend", userHandler.HandleDeleteFriend)

	srv.BaseRouterGroup.POST("/user/:login/add_post", userHandler.HandleAddPost)

	// Список пользователей
	srv.BaseRouterGroup.GET("/users", userHandler.HandleUsersList)

	srv.BaseRouterGroup.GET("/feed", userHandler.HandleFeed)

	// Static
	srv.BaseRouterGroup.Static("/public/", "./static")

	srv.BaseRouterGroup.GET("/ws/feed/:login", wsHandler.HandleWS)

	srv.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	log.Printf("received signal %s, stopping program...", sig)

	srv.Shutdown()
	ch.Close()
	conn.Close()
	signal.Stop(sigCh)

	log.Println("program stopped")
}

func dbConnect(driver, connectionString string) (*sql.DB, error) {
	var connErr error

	for i := 1; i <= 5; i++ {
		fmt.Printf("trying to connect to DB, try %d\n", i)

		if i != 1 {
			time.Sleep(5 * time.Second)
		}

		db, err := sql.Open(driver, connectionString)
		if err != nil {
			connErr = err
			continue
		}

		if err := db.Ping(); err != nil {
			connErr = err
			continue
		}

		return db, connErr
	}

	return nil, connErr
}
