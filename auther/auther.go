package auther

import (
	"errors"
	"log"
	"math"
	"os"
	"time"

	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

// Redirect paths for authentication
type Paths struct {
	HostedDomain           string
	RedirectFailure        string
	DefaultRedirectSuccess string
}

// Establish the authentication as OpenID-Connect
const (
	providerName = "openid-connect"
	autherStore  = "auther-store"
)

var (
	provider goth.Provider

	PathConfig = Paths{
		RedirectFailure:        "/",
		DefaultRedirectSuccess: "/",
	}

	LogStd = log.New(os.Stdout, "[auther] ", log.Ldate|log.Ltime)
	LogErr = log.New(os.Stderr, "ERROR [auther] ", log.Ldate|log.Ltime)

	UserValidator = func(*goth.User) bool {
		return true
	}
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Establish a filesystem cookie store
	store := sessions.NewFilesystemStore(os.TempDir(), []byte(autherStore))
	// Set the MaxLength to avoid `securecookie: the value is too long` issue
	store.MaxLength(math.MaxInt64)
	gothic.Store = store

	gothic.GetProviderName = func(request *http.Request) (string, error) {
		return providerName, nil
	}
}

// Load the OAuth configuration files
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

// Helper function to retrieve the session based on the request
func getSession(request *http.Request) (goth.Session, error) {
	session, err := gothic.Store.Get(request, providerName+gothic.SessionName)
	if err != nil {
		return nil, errors.New("could not find a session for the request")
	}

	value := session.Values[providerName]
	if value == nil {
		return nil, errors.New("could not find a session for the request")
	}

	return provider.UnmarshalSession(value.(string))
}

// Initialize the OpenID authentication
// This is needed in order to load the configuration files associated with the ConnectID key
func InitProvider(domain string, authCallback string, clientID string, clientSecret string) error {
	if domain == "" || authCallback == "" || clientID == "" || clientSecret == "" {
		log.Fatal("Arguments cannot be empty")
		return errors.New("could not initialize authentication provider")
	}

	clientID, clientSecret, err := loadOauthConfig(clientID, clientSecret)
	if err != nil {
		log.Fatal("Failed loading oauth config files:\n", err)
		return err
	}

	provider, err = openidConnect.New(clientID, clientSecret, "https://"+domain+authCallback, "https://accounts.google.com/.well-known/openid-configuration", "email")
	if provider != nil {
		goth.UseProviders(provider)
		return nil
	} else {
		log.Fatal("Failed creating provider:\n", err)
		return err
	}
}

// Login entry-point to authenticate a user
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	rawURL, err := gothic.GetAuthURL(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := url.Parse(rawURL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// Inject the hosted domain
	query := url.Query()
	query.Add("hd", PathConfig.HostedDomain)

	session, err := gothic.Store.Get(request, providerName+gothic.SessionName)
	if err == nil && session != nil {

		// Inject the redirect URL for after after logging in
		// This will match the current URL where the user is coming from
		session.Values["redirect"] = request.Header.Get("Referer")
		session.Save(request, response)
	} else {
		if session == nil {
			LogErr.Println("Session is null")
		}
		if err != nil {
			LogErr.Println("Error found: ", err)
		}
		http.Redirect(response, request, "/login", http.StatusPermanentRedirect)
		return
	}

	url.RawQuery = query.Encode()

	http.Redirect(response, request, url.String(), http.StatusTemporaryRedirect)
}

// Finalizes the login process
// This callcabk is called by Google when returning from their login screen
func AuthCallback(response http.ResponseWriter, request *http.Request) {
	url := request.URL

	// Check hosted domain
	{
		if hd := url.Query().Get("hd"); PathConfig.HostedDomain != "" && hd != PathConfig.HostedDomain {
			LogStd.Println("Hosted domain did not match. Got", hd)
			gothic.Logout(response, request)
			http.Redirect(response, request, PathConfig.RedirectFailure, http.StatusPermanentRedirect)
			return
		}
	}

	redirectSuccess := PathConfig.DefaultRedirectSuccess

	// Prepare redirect
	session, err := gothic.Store.Get(request, providerName+gothic.SessionName)
	if session != nil && err == nil {
		redirectSuccess = session.Values["redirect"].(string)
	}

	user, err := gothic.CompleteUserAuth(response, request)
	if err != nil {
		LogErr.Println("failed to complete login:", err)
		gothic.Logout(response, request)
		http.Redirect(response, request, PathConfig.RedirectFailure, http.StatusPermanentRedirect)
		return
	}

	// Hand-over finalization to calling process callback
	if !UserValidator(&user) {
		LogStd.Println("user invalid:", user)
		gothic.Logout(response, request)
		http.Redirect(response, request, PathConfig.RedirectFailure, http.StatusPermanentRedirect)
		return
	}

	http.Redirect(response, request, redirectSuccess, http.StatusPermanentRedirect)
}

// Retrieve the session user
func GetUser(response http.ResponseWriter, request *http.Request) (goth.User, error) {
	session, err := getSession(request)
	if err != nil || session == nil {
		return goth.User{}, err
	}

	return provider.FetchUser(session)
}

// Logs out and redirects to the URL from where the user is coming from
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	if err := gothic.Logout(response, request); err != nil {
		LogErr.Println("Failed: session is null")
	}

	redirectSuccess := request.Header.Get("Referer")
	if redirectSuccess == "" {
		redirectSuccess = PathConfig.DefaultRedirectSuccess
	}

	http.Redirect(response, request, redirectSuccess, http.StatusSeeOther)
}
