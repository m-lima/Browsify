package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"net/http"

	"github.com/gorilla/pat"
	"github.com/m-lima/browsify/auther"
)

const authCallback = "/authcallback"

func uiHandler(response http.ResponseWriter, request *http.Request) {
	path := strings.Replace(request.URL.Path, "/ui", "", 1)

	if strings.HasPrefix(path, "/static") {
		path = "ui/build" + path
	} else {
		path = "ui/build"
	}

	_, err := os.Stat(path)

	if err == nil {
		http.ServeFile(response, request, path)
	} else {
		response.WriteHeader(http.StatusNotFound)
		fmt.Fprint(response, "404 Not found")
	}
}

func launchServer(patter http.Handler) {
	var waiter sync.WaitGroup
	waiter.Add(2)

	{
		go func() {
			defer waiter.Done()
			http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
				http.Redirect(response, request, "https://localhost"+request.URL.Path, 302)
			})
			err := http.ListenAndServe("", nil)
			if err != nil {
				log.Fatal("Could not start HTTP server:\n", err)
			}
		}()
	}

	{
		go func() {
			defer waiter.Done()
			err := http.ListenAndServeTLS("", "server.crt.hide", "server.key.hide", patter)
			if err != nil {
				log.Fatal("Could not start HTTPS server:\n", err)
			}
		}()
	}

	waiter.Wait()
}

func main() {
	auther.InitProvider("localhost", authCallback, "oauth.client.id.hide", "oauth.client.secret.hide")

	auther.PathConfig.DefaultRedirectSuccess = "ui"
	auther.PathConfig.HostedDomain = "telenordigital.com"
	Api = "/api"

	patter := pat.New()
	patter.Get(Api, ApiHandler)
	patter.Get("/user", UserHandler)
	patter.Get("/ui", uiHandler)
	patter.Get(authCallback, auther.AuthCallback)
	patter.Get("/login", auther.LoginHandler)
	patter.Post("/logout", auther.LogoutHandler)
	patter.Get("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/ui", http.StatusPermanentRedirect)
	})

	launchServer(patter)
}
