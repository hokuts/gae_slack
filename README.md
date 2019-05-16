# GAE Slack API Example

Small example of slack API using bot with GAE/Go

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

note: When you create channel via API, channel creator will not be the App bot, but the user who authorized the App.

## Test local server

Edit `app.yaml` and replace `<PUT YOUR ACCESS TOKEN HERE>` with token above.

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

1. Deploy this app.  
Or you can boot the app localy and use [ngrok](https://ngrok.com/).  
Event API needs public URL for this app.
1. Open your app settings
    1. Go to <https://api.slack.com/apps>
    1. Select your app
1. Set up your app
    1. Open `Event Subscriptions` from left menu
    1. Enable Events
    1. Input `Request URL`. It might be like `<YOUR APP ROOT>/slack/event`.
        * Slack will send verification request to your app.
    1. Add events you need to subscribe.
        * e.g. `message.channels`
        * You may have to add required scope.
    1. Add `App Unfurl Domains` if you need.
    1. Click `Save Changes`.
1. Prepare signature verification
    1. Open `Basic Information` from left menu
    1. Copy `Signing Secret`
    1. Edit `app.yaml` and replace `<PUT YOUR SIGNING SECRET HERE>` with signing secret above.

## Set up Slack OAuth

1. Open your app settings
    1. Go to <https://api.slack.com/apps>
    1. Select your app
1. Get Client ID and Client Secret
    1. Open `Basic Information` from left menu
    1. Copy `Client ID` and `Client Secret`
    1. Edit `app.yaml` and resplace `<PUT YOUR CLIENT ID HERE>` and `<PUT YOUR CLIENT SECRET HERE>` with above.  
    Make sure to set them as STRING. They may have to be in quotes, or they might be treated as numbers.
1. Set Redirect URL
    1. Open `OAuth & Permissions` from left menu
    1. Input `Redirect URLs`. It might be like `http://localhost:8080/slack/oauth/token`.
    1. Click `Save URLs`
    1. Edit `app.yaml` and replace `<PUT YOUR REDIRECT URL HERE>` with same URL as above.
1. Optional: Set Workspace ID
    1. Edit `app.yaml` and replace `<PUT YOUR WORKSPACE ID HERE>` with your slack workspace ID.  
    If you do this, you will not need to select the Workspace through authorization flow.

Now you can start authorization flow by accessing to <http://localhost:8080/slack/oauth/auth> with your web browser.
At end of the flow, you will see granted access token in the web browser.