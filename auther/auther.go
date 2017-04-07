package auther

import (
	"errors"
	"log"
	"math"
	"os"
	"time"

	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

const (
	providerName = "openid-connect"
	autherStore  = "auther-store"
)

var (
	domain       = "localhost"
	authCallback = "/authcallback"
	provider     goth.Provider
	states       map[string]string

	HostedDomain    = "telenordigital.com"
	RedirectSuccess = "/"
	RedirectFailure = "/"
	LogStd          = log.New(os.Stdout, "auther: ", 0)
	LogErr          = log.New(os.Stderr, "auther: ", 0)
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	store := sessions.NewFilesystemStore(os.TempDir(), []byte(autherStore))
	store.MaxLength(math.MaxInt64)
	gothic.Store = store

	gothic.GetProviderName = func(request *http.Request) (string, error) {
		return providerName, nil
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

func generateState() string {
	bytes := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		bytes[i] = byte(rand.Int())
	}

	sha := sha256.New()
	sha.Write(bytes)

	return hex.EncodeToString(sha.Sum(nil))
}

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

func InitProvider(newDomain string, newAuthCallback string) {
	if newDomain != "" {
		domain = newDomain
	}
	if newAuthCallback != "" {
		authCallback = newAuthCallback
	}

	states = make(map[string]string)

	clientID, clientSecret, err := loadOauthConfig("oauth.client.id.hide", "oauth.client.secret.hide")
	if err != nil {
		log.Fatal("Failed loading oauth config files:\n", err)
		return
	}

	provider, err = openidConnect.New(clientID, clientSecret, "https://"+domain+authCallback, "https://accounts.google.com/.well-known/openid-configuration", "email")
	if provider != nil {
		goth.UseProviders(provider)
	} else {
		log.Fatal("Failed creating provider:\n", err)
		return
	}
}

func LoginHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
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

	query := url.Query()
	query.Add("hd", HostedDomain)

	session, err := getSession(request)
	if err == nil && session != nil {
		state := generateState()
		states[session.Marshal()] = state
		query.Set("state", state)
	} else {
		LogErr.Println("Not setting state for session. Session not found.")
		if session == nil {
			LogErr.Println(" >> Session is null")
		}
		if err != nil {
			LogErr.Println(" >> Error found", err)
		}
		http.Redirect(response, request, "/login", http.StatusPermanentRedirect)
		return
	}

	url.RawQuery = query.Encode()

	http.Redirect(response, request, url.String(), http.StatusTemporaryRedirect)
}

func AuthCallback(response http.ResponseWriter, request *http.Request) {
	url := request.URL

	{
		if hd := url.Query().Get("hd"); hd != HostedDomain {
			LogStd.Println("Hosted domain did not match. Got", hd)
			gothic.Logout(response, request)
			http.Redirect(response, request, RedirectFailure, http.StatusForbidden)
			return
		}
	}

	{
		session, err := getSession(request)
		if session != nil && err == nil {
			state := states[session.Marshal()]
			queryState := url.Query().Get("state")
			if queryState != state {
				LogStd.Printf(`State did not match.
Expected: %s
     Got: %s`, state, queryState)
				gothic.Logout(response, request)
				http.Redirect(response, request, RedirectFailure, http.StatusForbidden)
				return
			}
		} else {
			LogStd.Println("Not checking for state. State for session not found.")
		}
	}

	_, err := gothic.CompleteUserAuth(response, request)

	if err != nil {
		gothic.Logout(response, request)
		http.Redirect(response, request, RedirectFailure, http.StatusForbidden)
		return
	}

	http.Redirect(response, request, RedirectSuccess, http.StatusPermanentRedirect)
}

func GetUser(response http.ResponseWriter, request *http.Request) (goth.User, error) {
	session, err := getSession(request)
	if err != nil || session == nil {
		return goth.User{}, err
	}

	return provider.FetchUser(session)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	if err := gothic.Logout(response, request); err != nil {
		LogErr.Println("Failed: session is null")
	}
	http.Redirect(response, request, RedirectSuccess, http.StatusPermanentRedirect)
}
