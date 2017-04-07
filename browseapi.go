package main

import (
	"os"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

const Api = "/api"

var home = os.Getenv("HOME")

type directoryEntry struct {
	Name      string
	Directory bool
	Size      int64
	Date      time.Time
}

func ApiHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")

	path := request.URL.Path
	systemPath := ""

	if path == Api || path == Api+"/" {
		systemPath = home
		path = Api
	} else {
		if path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		systemPath = home + "/" + path[len(Api)+1:]
	}

	// Not found
	if _, err := os.Stat(systemPath); err != nil {
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
				entries = append(entries, directoryEntry{
					Name:      file.Name(),
					Directory: file.IsDir(),
					Size:      file.Size(),
					Date:      file.ModTime(),
				})
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
