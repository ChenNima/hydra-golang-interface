package router

import (
	"context"
	"net/http"

	"felix.chen/login/internal/logger"
	"felix.chen/login/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

func loginPage(c *gin.Context) {
	loginChallenge := c.Query("login_challenge")
	hydraAdmin := util.GetHydraAdmin()
	loginRequest, err := hydraAdmin.Admin.GetLoginRequest(&admin.GetLoginRequestParams{
		LoginChallenge: loginChallenge,
		Context:        context.Background(),
	})
	if err != nil {
		logger.GetLogger().Error(err)
	}
	body := *(loginRequest.Payload)
	if *body.Skip {
		hydraAdmin.Admin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
			LoginChallenge: loginChallenge,
			Context:        context.Background(),
			Body: &models.AcceptLoginRequest{
				Subject: body.Subject,
			},
		})
		c.Redirect(http.StatusFound, body.Client.RedirectUris[0])
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Challenge": loginChallenge,
	})
}

func doLogin(c *gin.Context) {
	loginChallenge := c.Request.FormValue("challenge")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	remember := c.Request.FormValue("remember") == "1"

	if email != "felix.chen@sap.com" && password != "123" {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Challenge": loginChallenge,
			"Error":     "The username / password combination is not correct",
			"Remember":  remember,
		})
		return
	}

	hydraAdmin := util.GetHydraAdmin()
	res, err := hydraAdmin.Admin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
		LoginChallenge: loginChallenge,
		Context:        context.Background(),
		Body: &models.AcceptLoginRequest{
			Subject:     &email,
			Remember:    remember,
			RememberFor: 3600,
		},
	})
	if err != nil {
		logger.GetLogger().Error(err)
	}
	c.Redirect(http.StatusFound, *res.Payload.RedirectTo)
}
