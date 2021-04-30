package service

import (
	"encoding/base64"
	"regexp"
	"time"

	"felix.chen/login/internal/store"
	"github.com/gorilla/sessions"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	abrenderer "github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"
	"github.com/volatiletech/authboss/v3/otp/twofactor"
	"github.com/volatiletech/authboss/v3/otp/twofactor/totp2fa"
)

func GetAuthboss() *authboss.Authboss {
	ab := authboss.New()
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	sessionCookieName := "ab_blog"
	sessionStore := abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	cookieStore := abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((30 * 24 * time.Hour) / time.Second))
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore
	// ab.Config.Paths.RootURL = "http://localhost:8083"
	ab.Config.Paths.Mount = "/auth"

	// Set up our server, session and cookie storage mechanisms.
	// These are all from this package since the burden is on the
	// implementer for these.
	ab.Config.Storage.Server = store.NewMemStorer()
	// ab.Config.Storage.SessionState = sessionStore
	// ab.Config.Storage.CookieState = cookieStore

	// Another piece that we're responsible for: Rendering views.
	// Though note that we're using the authboss-renderer package
	// that makes the normal thing a bit easier.
	ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "ab_views")

	// We render mail with the authboss-renderer but we use a LogMailer
	// which simply sends the e-mail to stdout.
	ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "ab_views")

	// The preserve fields are things we don't want to
	// lose when we're doing user registration (prevents having
	// to type them again)
	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}

	// TOTP2FAIssuer is the name of the issuer we use for totp 2fa
	ab.Config.Modules.TOTP2FAIssuer = "ABBlog"
	ab.Config.Modules.ResponseOnUnauthed = authboss.RespondRedirect

	// Turn on e-mail authentication required
	ab.Config.Modules.TwoFactorEmailAuthRequired = false

	// This instantiates and uses every default implementation
	// in the Config.Core area that exist in the defaults package.
	// Just a convenient helper if you don't want to do anything fancy.
	defaults.SetCore(&ab.Config, false, false)

	// Here we initialize the bodyreader as something customized in order to accept a name
	// parameter for our user as well as the standard e-mail and password.
	//
	// We also change the validation for these fields
	// to be something less secure so that we can use test data easier.
	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]+`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength: 2,
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: false,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, nameRule},
			"recover_end": {passwordRule},
		},
		Confirms: map[string][]string{
			"register":    {"password", authboss.ConfirmPrefix + "password"},
			"recover_end": {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": {"email", "name", "password"},
		},
	}

	// Set up 2fa
	twofaRecovery := &twofactor.Recovery{Authboss: ab}
	if err := twofaRecovery.Setup(); err != nil {
		panic(err)
	}

	totp := &totp2fa.TOTP{Authboss: ab}
	if err := totp.Setup(); err != nil {
		panic(err)
	}

	// Initialize authboss (instantiate modules etc.)
	if err := ab.Init(); err != nil {
		panic(err)
	}

	return ab
}
