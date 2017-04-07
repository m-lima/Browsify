package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"io/ioutil"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/m-lima/browsify/auther"
)

const browse = "/browse"

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

func indexHandler(response http.ResponseWriter, request *http.Request) {
	user, err := auther.GetUser(response, request)
	if err == nil {
		fmt.Fprintf(response, "Logged in as "+user.Name)
	} else {
		fmt.Fprintf(response, "Nope")
	}

	// http.Redirect(response, request, browse, 302)

	// if _, err := gothic.CompleteUserAuth(response, request); err == nil {
	// 	http.Redirect(response, request, browse, 302)
	// 	return
	// }

	// http.ServeFile(response, request, "main.html")
}

func browseHandler(response http.ResponseWriter, request *http.Request) {
	// if _, err := gothic.CompleteUserAuth(response, request); err != nil {
	// 	gothic.BeginAuthHandler(response, request)
	// 	return
	// }

	// if true {
	// 	BrowseHandler(response, request)
	// 	return
	// }

	path := request.URL.Path
	var systemPath string

	if path == browse || path == browse+"/" {
		systemPath = home
		path = browse
	} else {
		if path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		systemPath = home + "/" + path[len(browse)+1:]
	}

	fmt.Println(systemPath)

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

func apiHandler(response http.ResponseWriter, request *http.Request) {
	_, err := auther.GetUser(response, request)
	if err == nil {
		ApiHandler(response, request)
	} else {
		response.WriteHeader(http.StatusForbidden)
	}
}

func uiHandler(response http.ResponseWriter, request *http.Request) {
	path := "frontend/build" + strings.Replace(request.URL.Path, "/ui", "", 1)
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
	auther.InitProvider("localhost", "")
	auther.RedirectSuccess = "ui"

	patter := pat.New()
	patter.Get(browse, browseHandler)
	patter.Get("/api", apiHandler)
	patter.Get("/ui", uiHandler)
	patter.Get("/authcallback", auther.AuthCallback)
	patter.Get("/login", auther.LoginHandler)
	patter.Get("/logout", auther.LogoutHandler)
	patter.Get("/", indexHandler)

	launchServer(patter)
}
