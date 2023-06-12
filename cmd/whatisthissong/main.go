package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"whatisthissong/cmd/whatisthissong/config"
	"whatisthissong/cmd/whatisthissong/controller"
	"whatisthissong/cmd/whatisthissong/server"

	"github.com/gocolly/colly/v2"
	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI))
	ch   = make(chan *spotify.Client)
)

// crawl
func Crawl() {
	c := colly.NewCollector()

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		h.Request.Visit(h.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}

func LoginToSpotify() {

}

func main() {
	log.Println("starting server...")
	cnf, err := config.AppConfig()
	//로컬 저장소 없으면 생성

	c, err := controller.NewController()
	if err != nil {
		panic(err)
	}
	r := c.NewRouter()
	svr := server.NewServer(r, fmt.Sprintf(":%d", cnf.Port))
	go func() {
		if err := svr.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error while running server: %s\n", err.Error())
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	// stop server
	if err := svr.Stop(ctx); err != nil {
		log.Fatalf("failed to stop server: %v", err)
	}
}
