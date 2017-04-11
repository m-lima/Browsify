package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"net/http"

	"github.com/gorilla/pat"
	"github.com/m-lima/browsify/auther"
)

const (
	staticPath = "/static"
	uiPath     = "ui/build"
)

var (
	authCallback = "/authcallback"
	host         = "localhost"
	clientID     = "oauth.client.id.hide"
	clientSecret = "oauth.client.secret.hide"
	serverCert   = "server.crt.hide"
	serverKey    = "server.key.hide"
	hostedDomain = ""
	home         = ""
	ui           = "/ui"
)

func uiHandler(response http.ResponseWriter, request *http.Request) {
	path := strings.Replace(request.URL.Path, ui, "", 1)

	if strings.HasPrefix(path, staticPath) {
		path = uiPath + path
	} else {
		path = uiPath
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
				http.Redirect(response, request, "https://"+host+request.URL.Path, http.StatusPermanentRedirect)
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
			err := http.ListenAndServeTLS("", serverCert, serverKey, patter)
			if err != nil {
				log.Fatal("Could not start HTTPS server:\n", err)
			}
		}()
	}

	waiter.Wait()
}

func handleFlags() {
	flag.StringVar(&authCallback, "auth", authCallback, "callback path for authentication")
	flag.StringVar(&host, "host", host, "the host for this server")
	flag.StringVar(&clientID, "cid", clientID, "file path for the client ID file")
	flag.StringVar(&clientSecret, "cs", clientSecret, "file path for the client secret file")
	flag.StringVar(&serverCert, "sc", serverCert, "file path for the server certificate file")
	flag.StringVar(&serverKey, "sk", serverKey, "file path for the server key file")
	flag.StringVar(&hostedDomain, "hd", hostedDomain, "authorized domains for OpenID authentication")
	flag.StringVar(&home, "home", home, "base path for browsing")
	flag.StringVar(&ui, "ui", ui, "URL path for main UI")

	flag.Parse()
}

func main() {
	handleFlags()

	err := auther.InitProvider(host, authCallback, clientID, clientSecret)
	if err != nil {
		log.Fatal("Could not start OpenID provider", err)
	}

	auther.PathConfig.DefaultRedirectSuccess = ui
	auther.PathConfig.HostedDomain = hostedDomain
	if home != "" {
		Home = home
	}

	patter := pat.New()
	patter.Get(Api, ApiHandler)
	patter.Get(User, UserHandler)
	patter.Get(ui, uiHandler)
	patter.Get(authCallback, auther.AuthCallback)
	patter.Get("/login", auther.LoginHandler)
	patter.Post("/logout", auther.LogoutHandler)
	patter.Get("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, ui, http.StatusPermanentRedirect)
	})

	launchServer(patter)
}
