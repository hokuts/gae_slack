# GAE Slack API Example

Small example of slack API with GAE/Go

## Environment

TODO: GAE, GOPATH etc

```bash
go get github.com/gorilla/mux
go get github.com/tidwall/gjson
```

## Set up your slack App

1. Sign in to slack workspace as a admin
1. Create App
    1. Go to <https://api.slack.com/apps>
    1. Click `Create New App`
    1. Input App name and select workspace
    1. Click save button
1. Set up App (If you started here, go to <https://api.slack.com/apps> and select your App)
    1. Open `OAuth & Permissions` from left menu
    1. Select scopes below and save
        * `channels:read`
        * `channels:write`
        * `chat:write:bot`
1. Install App
    1. Open `Basic Information` from left menu
    1. Click `Install your app to your workspace`
    1. Authorize your App
1. Get access token
    1. Open `Install App` from left menu
    1. Copy `OAuth Access Token`

note: You don't need to create a Bot User.  
note: When you create channel via API, channel creator will not be the App bot, but the user who authorized the App.

## Test local server

Edit `slack.go` and replace `<PUT YOUR ACCESS TOKEN HERE>` with token above.

Boot local server

```bash
goapp serve .
```

Call API

```bash
curl http://localhost:8080/slack/channels
curl -X POST http://localhost:8080/slack/channels -d "name=test_channel"
curl -X POST http://localhost:8080/slack/channels/general/messages -d "text=MY_MESSAGE"
```

## Set up Event API

1. Deploy this project or use [ngrok](https://ngrok.com/).  
Below steps need public URL for this app.
1. Go to <https://api.slack.com/apps>
1. Select your app
1. Open `Event Subscriptions` from left menu
1. Enable Events
1. Input `Request URL`. It might be like `<YOUR APP ROOT>/slack/event_endpoint`.
    * Slack will send verification request to your app.
1. Add events you need to subscribe.
    * e.g. `message.channels`
    * You may have to add required scope.
1. Add `App Unfurl Domains` if you need.
1. Click `Save Changes`.