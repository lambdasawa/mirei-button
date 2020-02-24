package twitter

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var (
	callback = os.Getenv("MB_TWTITER_CALLBACK")

	client = oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
		Credentials: oauth.Credentials{
			Token:  os.Getenv("MB_TWITTER_CONSUMER_KEY"),
			Secret: os.Getenv("MB_TWITTER_CONSUMER_SECRET"),
		},
	}

	sessionKey = "session"

	tempCredKey   = "temp-cred"
	tokenCredKey  = "token-cred"
	screenNameKey = "screen-name"

	sessionOption = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	screenName = os.Getenv("MB_TWITTER_SCREENNAME")
)

func init() {
	gob.Register(&oauth.Credentials{})
}

func SignIn(c echo.Context) error {
	tempCred, err := client.RequestTemporaryCredentials(nil, callback, nil)
	if err != nil {
		return fmt.Errorf("get temporary credential: %v", err)
	}

	sess, _ := session.Get(sessionKey, c)
	sess.Options = sessionOption
	sess.Values[tempCredKey] = tempCred
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("save temporary credential: %v", err)
	}

	return c.Redirect(http.StatusFound, client.AuthorizationURL(tempCred, nil))
}

func Callback(c echo.Context) error {
	sess, err := session.Get(sessionKey, c)
	if err != nil {
		return fmt.Errorf("get session: %v", err)
	}

	tempCred, ok := sess.Values[tempCredKey].(*oauth.Credentials)
	if !ok {
		return fmt.Errorf("get temporary credential")
	}

	tokenCred, _, err := client.RequestToken(nil, tempCred, c.Request().FormValue("oauth_verifier"))
	if err != nil {
		return fmt.Errorf("fetch request token: %v", err)
	}

	sn, err := fetchScreenName(c, &client.Credentials, tokenCred)
	if err != nil {
		return fmt.Errorf("fetch screen name: %v", err)
	}
	if sn != screenName {
		return fmt.Errorf("invalid screen name")
	}

	delete(sess.Values, tempCredKey)
	sess.Values[tokenCredKey] = tokenCred
	sess.Values[screenNameKey] = screenName
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("delete temporary credential & save token credential: %v", err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func SignOut(c echo.Context) error {
	sess, _ := session.Get(sessionKey, c)
	sess.Options = sessionOption
	sess.Values = map[interface{}]interface{}{}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("delete temporary credential & delete token credential: %v", err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func Status(c echo.Context) error {
	sess, _ := session.Get(sessionKey, c)
	sess.Options = sessionOption
	screenName, ok := sess.Values[screenNameKey].(string)
	if !ok {
		return fmt.Errorf("find screen name")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": screenName,
	})
}

func fetchScreenName(c echo.Context, consumer, access *oauth.Credentials) (string, error) {
	httpClient := oauth1.NewConfig(client.Credentials.Token, client.Credentials.Secret).Client(oauth1.NoContext, oauth1.NewToken(access.Token, access.Secret))
	twitterClient := twitter.NewClient(httpClient)

	user, resp, err := twitterClient.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{
		IncludeEntities: twitter.Bool(false),
		SkipStatus:      twitter.Bool(true),
		IncludeEmail:    twitter.Bool(false),
	})
	if err != nil {
		return "", fmt.Errorf("verify credential: %v", err)
	}
	defer resp.Body.Close()

	return user.ScreenName, nil
}
