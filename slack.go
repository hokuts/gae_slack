package slack
import (
	"net/http"
	"net/url"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

const accessToken = "<PUT YOUR ACCESS TOKEN HERE>"

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/slack/channels", listChannel).Methods("GET")
	r.HandleFunc("/slack/channels", createChannel).Methods("POST")
	r.HandleFunc("/slack/channels/{channel}/messages", postMessage).Methods("POST")
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
		"name": {name},
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
		"token": {accessToken},
		"channel": {channel},
		"text": {text},
	}
	_, err := client.PostForm("https://slack.com/api/chat.postMessage", params)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "posted message to `%s`", channel)
}
