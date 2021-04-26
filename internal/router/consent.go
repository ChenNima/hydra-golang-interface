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

func consentPage(c *gin.Context) {
	consentChallenge := c.Query("consent_challenge")
	hydraAdmin := util.GetHydraAdmin()
	loginRequest, err := hydraAdmin.Admin.GetConsentRequest(&admin.GetConsentRequestParams{
		ConsentChallenge: consentChallenge,
		Context:          context.Background(),
	})
	if err != nil {
		logger.GetLogger().Error(err)
	}
	body := loginRequest.Payload
	if body.Skip {
		hydraAdmin.Admin.AcceptConsentRequest(&admin.AcceptConsentRequestParams{
			ConsentChallenge: consentChallenge,
			Context:          context.Background(),
			Body: &models.AcceptConsentRequest{
				GrantScope:               body.RequestedScope,
				GrantAccessTokenAudience: body.RequestedAccessTokenAudience,
			},
		})
		c.Redirect(http.StatusFound, body.Client.RedirectUris[0])
		return
	}
	c.HTML(http.StatusOK, "consent.html", gin.H{
		"Challenge":      consentChallenge,
		"User":           body.Subject,
		"Client":         body.Client,
		"Action":         "/consent",
		"RequestedScope": body.RequestedScope,
	})
}

func doConsent(c *gin.Context) {
	hydraAdmin := util.GetHydraAdmin()
	challenge := c.Request.FormValue("challenge")
	submit := c.Request.FormValue("submit")
	grantScope := c.Request.PostForm["grant_scope"]
	remember := c.Request.FormValue("remember") == "1"

	if submit == "Deny access" {
		res, _ := hydraAdmin.Admin.RejectConsentRequest(&admin.RejectConsentRequestParams{
			Context:          context.Background(),
			ConsentChallenge: challenge,
			Body: &models.RejectRequest{
				Error:            "access_denied",
				ErrorDescription: "The resource owner denied the request",
			},
		})
		c.Redirect(http.StatusFound, *res.Payload.RedirectTo)
	}
	body, _ := hydraAdmin.Admin.GetConsentRequest(&admin.GetConsentRequestParams{
		ConsentChallenge: challenge,
		Context:          context.Background(),
	})
	res, _ := hydraAdmin.Admin.AcceptConsentRequest(&admin.AcceptConsentRequestParams{
		ConsentChallenge: challenge,
		Context:          context.Background(),
		Body: &models.AcceptConsentRequest{
			GrantScope:               grantScope,
			GrantAccessTokenAudience: body.GetPayload().RequestedAccessTokenAudience,
			Remember:                 remember,
			RememberFor:              3600,
		},
	})
	c.Redirect(http.StatusFound, *res.GetPayload().RedirectTo)
}
