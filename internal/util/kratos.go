package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ory/kratos-client-go/client"
)

var kratosPublicURL, _ = url.Parse("http://localhost:4433")
var kratosPublic = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{kratosPublicURL.Scheme}, Host: kratosPublicURL.Host, BasePath: kratosPublicURL.Path})

func GetKratosPublic() *client.OryKratos {
	return kratosPublic
}

type kratoUiMessage struct {
	Text string `json:"text"`
}
type kratoUi struct {
	Messages []kratoUiMessage `json:"messages"`
}
type kratoFailRes struct {
	Id   string  `json:"id"`
	Type string  `json:"type"`
	Ui   kratoUi `json:"ui"`
}

func KratosSelfService(action string, values url.Values, flow string) error {
	kratosUrlString := kratosPublicURL.String() + "/self-service/" + action
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
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if res.StatusCode > 399 {
		var data kratoFailRes
		err = decoder.Decode(&data)
		if err == nil && len(data.Ui.Messages) > 0 {
			return errors.New(data.Ui.Messages[0].Text)
		}
		return errors.New(action + " failed!")
	}
	return nil
}
