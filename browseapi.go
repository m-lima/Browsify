package main

import (
	"os"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/m-lima/browsify/auther"
)

const (
	Api  = "/api/"
	User = "/user/"
)

var (
	Home          = os.Getenv("HOME")
	ShowHidden    = false
	ShowProtected = false
)

type directoryEntry struct {
	Name      string
	Directory bool
	Size      int64
	Date      time.Time
}

func shouldDisplayFile(file os.FileInfo) bool {
	if ShowProtected || file.Mode()&0004 != 0 {
		if ShowHidden || file.Name()[0] != '.' {
			return true
		}
	}
	return false
}

func ApiHandler(response http.ResponseWriter, request *http.Request) {
	_, err := auther.GetUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	path := request.URL.Path
	systemPath := ""

	if path == Api || path == Api[:len(Api)-1] {
		systemPath = Home
		path = Api
	} else {
		if path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		systemPath = Home + "/" + path[len(Api):]
	}

	// Not found
	if _, err := os.Stat(systemPath); err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	// Should not display
	if file, err := os.Stat(systemPath); err != nil || !shouldDisplayFile(file) {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	// Try folder
	{
		files, err := ioutil.ReadDir(systemPath)
		if err == nil {
			response.Header().Set("Content-type", "application/json")
			var entries []directoryEntry
			for _, file := range files {
				if shouldDisplayFile(file) {
					entries = append(entries, directoryEntry{
						Name:      file.Name(),
						Directory: file.IsDir(),
						Size:      file.Size(),
						Date:      file.ModTime(),
					})
				}
			}
			json.NewEncoder(response).Encode(entries)
			return
		}
	}

	// Try file
	{
		http.ServeFile(response, request, systemPath)
		return
	}
}

func UserHandler(response http.ResponseWriter, request *http.Request) {
	user, _ := auther.GetUser(response, request)
	json.NewEncoder(response).Encode(user)
}
