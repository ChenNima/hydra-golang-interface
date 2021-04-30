package util

import (
	"net/url"

	"github.com/ory/kratos-client-go/client"
)

var kratosPublicURL, _ = url.Parse("http://localhost:4433")
var kratosPublic = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{kratosPublicURL.Scheme}, Host: kratosPublicURL.Host, BasePath: kratosPublicURL.Path})

func GetKratosPublic() *client.OryKratos {
	return kratosPublic
}

func GetKratosPublicUrl() *url.URL {
	return kratosPublicURL
}
