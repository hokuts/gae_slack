package slack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const accessToken = "<PUT YOUR ACCESS TOKEN HERE>"
const signingSecret = "<PUT YOUR SIGNING SECRET HERE>"

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/slack/channels", listChannel).Methods("GET")
	r.HandleFunc("/slack/channels", createChannel).Methods("POST")
	r.HandleFunc("/slack/channels/{channel}/messages", postMessage).Methods("POST")
	r.HandleFunc("/slack/event_endpoint", handleEvent).Methods("POST")
	http.Handle("/", r)
}

func listChannel(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	params := url.Values{
		"token": {accessToken},
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
		"token": {accessToken},
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
		"token":   {accessToken},
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
	hm := hmac.New(sha256.New, []byte(signingSecret))
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
