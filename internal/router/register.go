package router

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"felix.chen/login/internal/logger"
	"felix.chen/login/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/ory/kratos-client-go/client/public"
)

func registerPage(c *gin.Context) {
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
	c.HTML(http.StatusOK, "register.html", gin.H{
		"Challenge": loginChallenge,
	})
}

type Identity struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func register(values url.Values, flow string) error {
	kratosUrlString := util.GetKratosPublicUrl().String() + "/self-service/registration"
	kratosUrl, _ := url.Parse(kratosUrlString)
	q := kratosUrl.Query()
	q.Set("flow", flow)
	kratosUrl.RawQuery = q.Encode()
	client := &http.Client{}
	req, err := http.NewRequest("POST", kratosUrl.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(values.Encode())))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode > 399 {
		return errors.New("Register failed!")
	}
	return nil
}

func doRegister(c *gin.Context) {
	loginChallenge := c.Request.FormValue("challenge")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	remember := c.Request.FormValue("remember") == "1"

	kratosPublic := util.GetKratosPublic()

	registerFlow, err := kratosPublic.Public.InitializeSelfServiceRegistrationViaAPIFlow(&public.InitializeSelfServiceRegistrationViaAPIFlowParams{
		Context: context.Background(),
	})
	// err := kratosPublic.Public.InitializeSelfServiceRegistrationViaBrowserFlow(&public.InitializeSelfServiceRegistrationViaBrowserFlowParams{
	// 	Context:    context.Background(),
	// 	HTTPClient: &client,
	// })
	if err != nil {
		log.Error(err)
	}
	flow := string(*registerFlow.GetPayload().ID)

	// _, err = kratosPublic.Public.CompleteSelfServiceRegistrationFlowWithPasswordMethod(&public.CompleteSelfServiceRegistrationFlowWithPasswordMethodParams{
	// 	Context: context.Background(),
	// 	Flow:    &flow,
	// 	Payload: &Identity{
	// 		Email:    email,
	// 		Password: password,
	// 	},
	// })
	registerData := url.Values{}
	registerData.Add("traits.email", email)
	registerData.Add("password", password)
	registerData.Add("method", "password")
	err = register(registerData, flow)
	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"Challenge": loginChallenge,
			"Error":     err.Error(),
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
