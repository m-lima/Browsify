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
	uiPath     = "web"

	authCallback = "/authcallback"
	login        = "/login"
	logout       = "/logout"
)

var (
	configFile     = flag.String("c", "securidash.conf", "Configuration file")
	generateConfig = flag.String("g", "", "File to be generated as default configuration")
	newUserEmail   = flag.String("a", "", "User email to be added as admin")
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
				http.Redirect(response, request, "https://"+Configuration.Server.Host+request.URL.Path, http.StatusPermanentRedirect)
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
			err := http.ListenAndServeTLS("", Configuration.Ssl.Certificate, Configuration.Ssl.Key, mux)
			if err != nil {
				log.Fatal("Could not start HTTPS server:\n", err)
			}
		}()
	}

	waiter.Wait()
}

func main() {
	flag.Parse()

	if *generateConfig != "" {
		err := GenerateDefaultConfig(*generateConfig)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := LoadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer FinalizeDB()

	if *newUserEmail != "" {
		err = CreateUser(&User{
			Email:            *newUserEmail,
			Admin:            true,
			CanShowHidden:    true,
			CanShowProtected: true,
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = auther.InitProvider(Configuration.Server.Host, authCallback, Configuration.Oauth.Id, Configuration.Oauth.Secret)
	if err != nil {
		log.Fatal("Could not start OpenID provider", err)
	}

	auther.PathConfig.DefaultRedirectSuccess = Configuration.Server.Ui
	auther.PathConfig.HostedDomain = Configuration.Server.HostedDomain
	auther.UserValidator = ValidateUser

	mux := http.NewServeMux()

	// Main redirect
	mux.Handle("/", http.RedirectHandler(Configuration.Server.Ui, http.StatusPermanentRedirect))

	// Web UI routes
	mux.HandleFunc("/favicon.ico", func(response http.ResponseWriter, request *http.Request) {
		http.ServeFile(response, request, uiPath+"/favicon.ico")
	})
	mux.HandleFunc(Configuration.Server.Ui, func(response http.ResponseWriter, request *http.Request) {
		http.ServeFile(response, request, uiPath+"/index.html")
	})
	mux.HandleFunc(staticPath, staticHandler)

	// Api routes
	mux.HandleFunc(ApiURL, ApiHandler)
	mux.HandleFunc(UserURL, UserHandler)
	mux.HandleFunc(UserUpdateURL, UserUpdateHandler)

	// Auth routes
	mux.HandleFunc(authCallback, auther.AuthCallback)
	mux.HandleFunc(login, auther.LoginHandler)
	mux.HandleFunc(logout, auther.LogoutHandler)

	launchServer(mux)
}
