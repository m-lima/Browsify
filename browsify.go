package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"

	"io/ioutil"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

const browse = "/browse"
const providerName = "openid-connect"

const pageHeader = `
<html>
	<head>
		<title>Browsify</title>
	</head>
	<body>
`

const linkString = `
		<a href="%s">%s</a>
		<br>
`

func init() {
	store := sessions.NewFilesystemStore(os.TempDir(), []byte("browsify"))
	store.MaxLength(math.MaxInt64)
	gothic.Store = store

	gothic.GetProviderName = func(request *http.Request) (string, error) {
		return providerName, nil
	}
}

func authCallback(response http.ResponseWriter, request *http.Request) {
	_, err := gothic.CompleteUserAuth(response, request)
	if err != nil {
		fmt.Fprintln(response, err)
		return
	}
	http.Redirect(response, request, browse, 302)
}

func indexHandler(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, browse, 302)
	// if _, err := gothic.CompleteUserAuth(response, request); err == nil {
	// 	http.Redirect(response, request, browse, 302)
	// 	return
	// }

	// http.ServeFile(response, request, "main.html")
}

func browseHandler(response http.ResponseWriter, request *http.Request) {
	if _, err := gothic.CompleteUserAuth(response, request); err != nil {
		gothic.BeginAuthHandler(response, request)
		return
	}

	path := request.URL.Path
	systemPath := ""

	if path == browse || path == browse+"/" {
		systemPath = home
		path = browse
	} else {
		if path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		systemPath = home + path[len(browse)+1:]
	}

	// Try folder
	{
		files, err := ioutil.ReadDir(systemPath)
		if err == nil {
			fmt.Fprint(response, pageHeader)
			fmt.Fprint(response, `
	<form method="post" action="/logout">
		<button type="submit">Logout</button>
	</form>
	`)

			if systemPath != home {
				upDir := path[:strings.LastIndex(path, "/")]
				fmt.Fprintf(response, linkString, upDir, "..")
			}

			for _, file := range files {
				fmt.Fprintf(response, linkString, path+"/"+file.Name(), file.Name())
			}

			fmt.Fprintln(response, "</html>")
			return
		}
	}

	// Try file
	{
		_, err := os.Stat(systemPath)

		if err == nil {
			http.ServeFile(response, request, systemPath)
			return
		}
	}

	// Not found
	response.WriteHeader(http.StatusNotFound)
	fmt.Fprint(response, "404 Not found")
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := gothic.Store.Get(request, providerName+gothic.SessionName)
	if session != nil {
		session.Options.MaxAge = -1
		if err := gothic.Store.Save(request, response, session); err != nil {
			fmt.Fprintln(os.Stderr, "Failed saving at store:", err)
			return
		}
		http.Redirect(response, request, "/", 302)
	} else {
		fmt.Fprintln(os.Stderr, "Failed: session is null")
	}
}

func loadOauthConfig(idFile string, secretFile string) (id string, secret string, err error) {
	file, err := ioutil.ReadFile(idFile)
	if err != nil {
		return "", "", err
	}
	id = string(file[:len(file)-1])

	file, err = ioutil.ReadFile(secretFile)
	if err != nil {
		return "", "", err
	}
	secret = string(file[:len(file)-1])

	return
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
	clientId, clientSecret, err := loadOauthConfig("oauth.client.id.hide", "oauth.client.secret.hide")
	if err != nil {
		log.Fatal("Failed loading oauth config files:\n", err)
		return
	}

	provider, err := openidConnect.New(clientId, clientSecret, "https://localhost/authcallback", "https://accounts.google.com/.well-known/openid-configuration")
	if provider != nil {
		goth.UseProviders(provider)
	} else {
		log.Fatal("Failed creating provider:\n", err)
		return
	}

	patter := pat.New()
	patter.Get(browse, browseHandler)
	patter.Get("/authcallback", authCallback)
	patter.Post("/logout", logoutHandler)
	patter.Get("/", indexHandler)

	launchServer(patter)
}
