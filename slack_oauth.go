package slack

import (
	"context"
	"encoding/json"
	"log"

	"golang.org/x/oauth2"
)

const slackAppClientID = "<PUT YOUR CLIENT ID HERE>"
const slackAppClientSecret = "<PUT YOUR CLIENT SECRET HERE>"
const redirectURL = "<PUT YOUR REDIRECT URL HERE>"
const slackWorkspaceID = ""

func getSlackAuthURL(ctx context.Context) (string, error) {
	stateString := generateState()
	return slackOAuthConfig().AuthCodeURL(stateString), nil
}

func getSlackOAuthAccessToken(ctx context.Context, code string, stateString string) (string, error) {

	err := verifyState(stateString)
	if err != nil {
		return "", err
	}

	token, err := slackOAuthConfig().Exchange(ctx, code)
	if err != nil {
		log.Printf("cannot get oauth access token: %v", err)
		return "", err
	}

	t, err := json.Marshal(token)
	if err != nil {
		log.Printf("cannot generate oauth token json: %v", err)
		return "", err
	}

	return string(t), nil
}

func slackOAuthConfig() *oauth2.Config {
	var authURL string
	if slackWorkspaceID != "" {
		authURL = "https://" + slackWorkspaceID + ".slack.com/oauth/authorize"
	} else {
		authURL = "https://slack.com/oauth/authorize"
	}
	return &oauth2.Config{
		ClientID:     slackAppClientID,
		ClientSecret: slackAppClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: "https://slack.com/api/oauth.access",
		},
		Scopes:      []string{"client"},
		RedirectURL: redirectURL,
	}
}

//// Example (You can add expire date, etc...)
// type OAuthState struct {
// 	UID         string
// 	StateString string
// }

func generateState() string {
	return ""

	//// Example
	// user := getCurrentUser()
	// uuid, err := uuid.NewRandom()
	// if err != nil {
	// 	return "", err
	// }
	// stateString := uuid.String()
	// state := OAuthState{
	// 	UID:         user.UID,
	// 	StateString: stateString,
	// }
	// key := datastore.NewKey(ctx, "OAuthState", stateString, 0, nil)
	// datastore.Put(ctx, key, &state)
	// return stateString
}

func verifyState(stateString string) error {
	return nil

	//// Example
	// user := getCurrentUser()
	// key := datastore.NewKey(ctx, "OAuthState", stateString, 0, nil)
	// var state OAuthState
	// err := datastore.Get(ctx, key, &state)
	// if err == datastore.ErrNoSuchEntity {
	// 	log.Print(err.Error())
	// 	return errors.New("invalid oauth state")
	// } else if err != nil {
	// 	log.Print(err.Error())
	// 	return err
	// } else if state.UID != user.UID {
	// 	log.Printf("UID did not match: `%s`, `%s`", state.UID, user.UID)
	// 	return errors.New("invalid oauth state")
	// }
	// return nil
}
