package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sk-todos/auth"
	"sk-todos/router"
	"sk-todos/store"
	"sk-todos/todo"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(100, 1)

func limitedHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func main() {
	// Liveness Probe -> Create file then kube I check that there's that file or not
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer os.Remove("tmp/live")

	// ENV
	err = godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variables: %s \n", err.Error())
	}

	// Database
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	// Mongo
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongoadmin:secret@localhost:27017"))
	if err != nil {
		panic("fail to connect mongodb")
	}
	collection := client.Database("myapp-mongo").Collection("todos")

	r := router.NewMyRouter()

	// Readiness Probe
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Load test
	r.GET("/load-test", limitedHandler)

	// lgflags
	r.GET("/x", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))
	//protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	//gormStore := store.NewGormStore(db)
	//todoHandler := todo.NewTodoHandler(gormStore)

	mongoStore := store.NewMongoDBStore(collection)
	todoHandler := todo.NewTodoHandler(mongoStore)

	r.POST("/todos", todoHandler.NewTask)
	//protected.POST("/todos", router.NewGinHandler(todoHandler.NewTask))
	//protected.GET("/todos", todoHandler.List)
	//protected.DELETE("/todos/:id", todoHandler.Delete)

	// Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err.Error())
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}
