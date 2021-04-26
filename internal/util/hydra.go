package util

import (
	"net/url"

	"github.com/ory/hydra-client-go/client"
)

var adminURL, _ = url.Parse("http://localhost:4445")
var hydraAdmin = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path})

func GetHydraAdmin() *client.OryHydra {
	return hydraAdmin
}
