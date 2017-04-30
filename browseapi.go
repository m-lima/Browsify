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
)

var (
	apiLogStd = log.New(os.Stdout, "[api] ", log.Ldate|log.Ltime)
	apiLogErr = log.New(os.Stderr, "ERROR [api] ", log.Ldate|log.Ltime)
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
	apiLogStd.Println(user.Email, "is requesting", path)

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

	apiLogStd.Println("Updating", user.Email)
	if user.CanShowHidden {
		user.ShouldShowHidden = request.PostFormValue("ShouldShowHidden") == "true"
		apiLogStd.Println("Show hidden:", user.ShouldShowHidden)
	} else if request.PostFormValue("ShouldShowHidden") == "true" {
		apiLogErr.Println(user.Email, "is trying to show hidden without permission")
	}
	if user.CanShowProtected {
		user.ShouldShowProtected = request.PostFormValue("ShouldShowProtected") == "true"
		apiLogStd.Println("Show protected:", user.ShouldShowProtected)
	} else if request.PostFormValue("ShouldShowProtected") == "true" {
		apiLogErr.Println(user.Email, "is trying to show protected without permission")
	}

	err = UpdateUser(&user)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		apiLogErr.Println("Could not update", user.Email, err)
		return
	}

	user, err = ReadUser(&sessionUser)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		apiLogErr.Println("Could not update", user.Email, err)
		return
	}

	apiLogStd.Println("Successfully updated", user.Email)
	json.NewEncoder(response).Encode(user)
}
