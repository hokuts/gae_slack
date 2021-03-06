package slack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/slack").Subrouter()

	ch := s.PathPrefix("/channels").Subrouter()
	ch.HandleFunc("", createChannel).Methods("POST")
	ch.HandleFunc("", listChannel).Methods("GET")
	s.HandleFunc("/{channel}/messages", postMessage).Methods("POST")

	oauth := s.PathPrefix("/oauth").Subrouter()
	oauth.HandleFunc("/auth", handleAuth).Methods("GET")
	oauth.HandleFunc("/token", handleAccessToken).Methods("GET")

	evt := s.PathPrefix("/event").Subrouter()
	evt.HandleFunc("", handleEvent).Methods("POST")

	http.Handle("/", r)
}

func listChannel(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	params := url.Values{
		"token": {os.Getenv("ACCESS_TOKEN")},
	}
	res, err := client.Get("https://slack.com/api/channels.list?" + params.Encode())
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "%q", res.Body) // TODO: Parse JSON
}

func createChannel(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	name := r.PostFormValue("name")
	params := url.Values{
		"token": {os.Getenv("ACCESS_TOKEN")},
		"name":  {name},
	}
	resp, err := client.PostForm("https://slack.com/api/channels.create", params)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "created channel `%s`\n%d\n%q", name, resp.StatusCode, resp.Body)
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	vars := mux.Vars(r)
	channel := vars["channel"]
	text := r.PostFormValue("text")
	params := url.Values{
		"token":   {os.Getenv("ACCESS_TOKEN")},
		"channel": {channel},
		"text":    {text},
	}
	res, err := client.PostForm("https://slack.com/api/chat.postMessage", params)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "posted message to `%s`\nStatus=%d\n%q", channel, res.StatusCode, res.Body)
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	url, err := getSlackAuthURL(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func handleAccessToken(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	code := r.FormValue("code")
	state := r.FormValue("state")
	token, err := getSlackOAuthAccessToken(ctx, code, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, token)
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	reqBody := buf.String()

	// Signature verification
	hm := hmac.New(sha256.New, []byte(os.Getenv("SIGNING_SECRET")))
	hm.Write([]byte("v0:" + r.Header.Get("X-Slack-Request-Timestamp") + ":" + reqBody))
	sig := "v0=" + hex.EncodeToString(hm.Sum(nil))
	if sig != r.Header.Get("X-Slack-Signature") {
		status := http.StatusBadRequest
		w.WriteHeader(status)
		fmt.Fprintf(w, http.StatusText(status))
		return
	}

	// You must respond to `url_verification` event
	_type := gjson.Get(reqBody, "type")
	if _type.Exists() && _type.String() == "url_verification" {
		log.Infof(ctx, "Event received:\n%s", reqBody)
		w.Header().Set("Content-type", "text/plain")
		fmt.Fprintf(w, gjson.Get(reqBody, "challenge").String())
		return
	}

	////// Handle your events below
	log.Infof(ctx, "Event received:\n%s", reqBody)
}
