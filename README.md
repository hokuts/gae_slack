Small example of slack API with GAE/Go

## Environment

TODO: GAE, GOPATH etc

```
$ go get github.com/gorilla/mux
```

## Set up your slack App

1. Sign in to slack workspace as a admin
1. Create App
    1. Go to https://api.slack.com/apps
    1. Click `Create New App`
    1. Input App name and select workspace
    1. Click save button
1. Set up App (If you started here, go to https://api.slack.com/apps and select your App)
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
```
$ goapp serve .
```

Call API
```
$ curl http://localhost:8080/slack/channels
$ curl -X POST http://localhost:8080/slack/channels -d "name=test_channel"
$ curl -X POST http://localhost:8080/slack/channels/general/messages -d "text=MY_MESSAGE"
```
