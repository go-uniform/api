package service

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-diary/diary"
	"net/http"
	"os"
	"service/service/_base"
	"service/service/info"
	"sync"
)

func RunAfter(shutdown chan bool, group *sync.WaitGroup, p diary.IPage) {
	port := fmt.Sprint(info.Args["port"])
	httpsCert := fmt.Sprint(info.Args["httpsCert"])
	httpsKey := fmt.Sprint(info.Args["httpsKey"])
	disableHttps := info.Args["disableHttps"].(bool)
	origin := fmt.Sprint(info.Args["origin"])
	jwtPub := fmt.Sprint(info.Args["jwtPub"])

	jwtPubData, err := os.ReadFile(jwtPub)
	if err != nil {
		panic(err)
	}
	info.JwtPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(jwtPubData)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	info.Engine = gin.Default()
	if err := info.Engine.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	// serve api html documentation on the root path
	p.Info("http.bind.main", diary.M{
		"path": "/",
	})
	info.Engine.Handle(http.MethodGet, "/", func(ctx *gin.Context) {
		writer := ctx.Writer
		writer.Header().Set("Content-Type", "text/html")
		writer.WriteHeader(200)
		// todo: generate and cache api documentation if not cached
		_, _ = writer.Write([]byte("todo: generate and cache api documentation if not cached"))
	})

	// serve openapi.json specification file
	p.Info("http.bind.openapi", diary.M{
		"path": "/openapi.json",
	})
	info.Engine.Handle(http.MethodGet, "openapi.json", func(ctx *gin.Context) {
		writer := ctx.Writer
		writer.Header().Set("Content-Type", "text/html")
		writer.WriteHeader(200)
		// todo: generate and cache api specification if not cached
		_, _ = writer.Write([]byte("todo: generate and cache api specification if not cached"))
	})

	if !disableHttps {
		if err := info.Engine.RunTLS(":8000", httpsCert, httpsKey); err != nil {
			panic(err)
		}
	} else {
		if err := info.Engine.Run(":8000"); err != nil {
			panic(err)
		}
	}

	srv := http.Server{
		Addr: ":" + port,
		// the always annoying CORS middleware, for added security of course ;)
		Handler: &_base.CorsMiddleware{Engine: info.Engine, Origin: origin},
	}
	p.Info("http.server", diary.M{
		"addr": ":" + port,
	})

	// wait for shutdown signal in separate thread
	go func() {
		group.Add(1)
		defer group.Done()

		// closing the shutdown chan will broadcast a close signal
		<-shutdown

		p.Notice("http.server.shutdown", diary.M{
			"addr": ":" + port,
		})

		if err := srv.Shutdown(context.TODO()); err != nil {
			p.Warning("http.server.stop.error", "failed to stop web server", diary.M{
				"addr":     ":" + port,
				"error":    err,
				"errorMsg": err.Error(),
			})
		} else {
			p.Notice("http.server.stop", diary.M{
				"addr": ":" + port,
			})
		}
	}()

	// run web server in separate thread
	go func() {
		group.Add(1)
		defer group.Done()

		p.Notice("http.server.start", diary.M{
			"addr": ":" + port,
		})

		if !disableHttps {
			fmt.Printf("\n\nhttps://127.0.0.1:%s\n\n\n", port)
			if err := srv.ListenAndServeTLS(httpsCert, httpsKey); err != nil {
				if err != http.ErrServerClosed {
					panic(err)
				}
			}
		} else {
			fmt.Printf("\n\nhttp://127.0.0.1:%s\n\n\n", port)
			if err := srv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					panic(err)
				}
			}
		}
	}()
}
