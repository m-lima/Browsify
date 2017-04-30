package main

import (
	"log"
	"os"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/m-lima/browsify/auther"
)

const (
	ApiURL        = "/api/"
	UserURL       = "/user"
	UserUpdateURL = "/user/update"
	logStd = log.New(os.Stdout, "api: ")
	logErr = log.New(os.Stderr, "api: ")
)

type directoryEntry struct {
	Name      string
	Directory bool
	Size      int64
	Date      time.Time
}

func corsProtection(response http.ResponseWriter, request *http.Request) {
	if Configuration.Server.DisableCors {
		response.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
		response.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func shouldDisplayFile(user *User, file os.FileInfo) bool {
	if !user.ShouldShowProtected && file.Mode()&0004 == 0 {
		return false
	}

	if !user.ShouldShowHidden && file.Name()[0] == '.' {
		return false
	}

	return true
}

func ApiHandler(response http.ResponseWriter, request *http.Request) {
	
	corsProtection(response, request)

	sessionUser, err := auther.GetUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := ReadUser(&sessionUser)
	if err != nil {
		response.WriteHeader(http.StatusForbidden)
		return
	}

	path := request.URL.Path
	systemPath := ""

	if path == ApiURL || path == ApiURL[:len(ApiURL)-1] {
		systemPath = Configuration.Server.Home
		path = ApiURL
	} else {
		if path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		systemPath = Configuration.Server.Home + "/" + path[len(ApiURL):]
	}
	logStd.Println("API:", user.Email, "is requesting", path)

	// Not found
	if _, err := os.Stat(systemPath); err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	// Should not display
	if file, err := os.Stat(systemPath); err != nil || !shouldDisplayFile(&user, file) {
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
				if shouldDisplayFile(&user, file) {
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
	corsProtection(response, request)

	sessionUser, err := auther.GetUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := ReadUser(&sessionUser)
	if err == nil {
		json.NewEncoder(response).Encode(user)
	} else {
		response.WriteHeader(http.StatusForbidden)
	}
}

func UserUpdateHandler(response http.ResponseWriter, request *http.Request) {
	corsProtection(response, request)

	if request.Method != "POST" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionUser, err := auther.GetUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := ReadUser(&sessionUser)
	if err != nil {
		response.WriteHeader(http.StatusForbidden)
		return
	}

	err = request.ParseForm()
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	user.ShouldShowHidden = userCanShowHidden ? request.PostFormValue("ShouldShowHidden") == "true" : false
	user.ShouldShowProtected = user.CanShowProtected ? request.PostFormValue("ShouldShowProtected") == "true" : false

	err = UpdateUser(&user)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err = ReadUser(&sessionUser)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	logStd.Println("API:", user.Email, "updated. Hidden =", user.ShouldShowHidden, "Protected =", user.ShouldShowProtected)

	json.NewEncoder(response).Encode(user)
}
