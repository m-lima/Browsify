package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"net/http"

	"github.com/m-lima/browsify/auther"
)

const (
	staticPath = "/static/"
	uiPath     = "web/build"

	authCallback = "/authcallback"
	login        = "/login"
	logout       = "/logout"
)

var (
	host         = "localhost"
	clientID     = "oauth.client.id.hide"
	clientSecret = "oauth.client.secret.hide"
	serverCert   = "server.crt.hide"
	serverKey    = "server.key.hide"
	hostedDomain = ""
	ui           = "/ui/"
)

func staticHandler(response http.ResponseWriter, request *http.Request) {
	path := uiPath + request.URL.Path
	file, _ := os.Stat(path)
	if file != nil && !file.IsDir() {
		http.ServeFile(response, request, path)
	} else {
		response.WriteHeader(http.StatusNotFound)
	}
}

func launchServer(mux *http.ServeMux) {
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
			err := http.ListenAndServeTLS("", serverCert, serverKey, mux)
			if err != nil {
				log.Fatal("Could not start HTTPS server:\n", err)
			}
		}()
	}

	waiter.Wait()
}

func handleFlags() {
	flag.StringVar(&host, "host", host, "the host for this server")
	flag.StringVar(&clientID, "cid", clientID, "file path for the client ID file")
	flag.StringVar(&clientSecret, "cs", clientSecret, "file path for the client secret file")
	flag.StringVar(&serverCert, "sc", serverCert, "file path for the server certificate file")
	flag.StringVar(&serverKey, "sk", serverKey, "file path for the server key file")
	flag.StringVar(&hostedDomain, "hd", hostedDomain, "authorized domains for OpenID authentication")
	flag.StringVar(&ui, "ui", ui, "URL path for main UI")
	flag.StringVar(&Home, "home", Home, "base path for browsing")
	flag.BoolVar(&ShowHidden, "sh", ShowHidden, "show hidden files")
	flag.BoolVar(&ShowProtected, "sp", ShowProtected, "show hidden files")

	flag.Parse()

	if ui[0] != '/' {
		ui = "/" + ui
	}

	if ui[len(ui)-1] != '/' {
		ui += "/"
	}

	if Home[len(ui)-1] == '/' {
		Home = Home[:len(Home)-1]
	}
}

func main() {
	handleFlags()

	err := auther.InitProvider(host, authCallback, clientID, clientSecret)
	if err != nil {
		log.Fatal("Could not start OpenID provider", err)
	}

	auther.PathConfig.DefaultRedirectSuccess = ui
	auther.PathConfig.HostedDomain = hostedDomain

	mux := http.NewServeMux()

	// Main redirect
	mux.Handle("/", http.RedirectHandler(ui, http.StatusPermanentRedirect))

	// Web UI routes
	mux.HandleFunc("/favicon.ico", func(response http.ResponseWriter, request *http.Request) {
		http.ServeFile(response, request, uiPath+"/favicon.ico")
	})
	mux.HandleFunc(ui, func(response http.ResponseWriter, request *http.Request) {
		http.ServeFile(response, request, uiPath+"/index.html")
	})
	mux.HandleFunc(staticPath, staticHandler)

	// Api routes
	mux.HandleFunc(Api, ApiHandler)
	mux.HandleFunc(User, UserHandler)

	// Auth routes
	mux.HandleFunc(authCallback, auther.AuthCallback)
	mux.HandleFunc(login, auther.LoginHandler)
	mux.HandleFunc(logout, auther.LogoutHandler)

	launchServer(mux)
}
