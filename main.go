package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"mymachine707/config"
	docs "mymachine707/docs" // docs is generated by Swag CLI, you have to import it.
	"mymachine707/handlars"
	blogpost "mymachine707/protogen/blogpost"
	"mymachine707/services/article"
	"mymachine707/services/author"
	"mymachine707/storage"
	"mymachine707/storage/postgres"

	_ "github.com/lib/pq"
)

func initGRPC(stg storage.Interfaces) {

	println("gRPC server tutorial in Go")

	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()

	AuthorService := author.NewAuthorService(stg)
	blogpost.RegisterAuthorServiceServer(s, AuthorService)

	ArticleService := article.NewArticleService(stg)
	blogpost.RegisterArticleServiceServer(s, ArticleService)

	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}

// @license.name Apache 2.0
func main() {

	fmt.Println("---------------------------------->>>")

	cfg := config.Load()

	psqlConfigString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	docs.SwaggerInfo.Title = cfg.App
	docs.SwaggerInfo.Version = cfg.AppVersion

	var err error
	var stg storage.Interfaces

	stg, err = postgres.InitDB(psqlConfigString)

	if err != nil {
		panic(err)
	}

	go initGRPC(stg)

	if cfg.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	fmt.Println("----->>")
	fmt.Printf("%+v\n", cfg)
	fmt.Println("---->>")

	r := gin.New()

	if cfg.Environment != "production" {
		r.Use(gin.Logger(), gin.Recovery()) // Later they will be replaced by custom Logger and Recovery
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	h := handlars.Handler{
		Stg: stg,
	}
	// Gruppirovka qilindi
	v1 := r.Group("v2")
	{
		v1.POST("/article", h.CreatArticle)
		v1.GET("/article/:id", h.GetArticleByID)
		v1.GET("/article", h.GetArticleList)
		v1.PUT("/article", h.ArticleUpdate)
		v1.DELETE("/article/:id", h.DeleteArticle)

		v1.POST("/author", h.CreatAuthor)
		v1.GET("/author/:id", h.GetAuthorByID)
		v1.GET("/author", h.GetAuthorList)
		v1.PUT("/author", h.AuthorUpdate)
		v1.DELETE("/author/:id", h.DeleteAuthor)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(cfg.HTTPPort) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
