package cas_connector

import (
	"embed"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/apache/answer-plugins/util"
	"github.com/apache/answer/plugin"
	"github.com/nizil/answer-plugin-cas/i18n"
)

//go:embed info.yaml
var Info embed.FS

type CasConnector struct {
	Config *CasConnectorConfig
}

func init() {
	plugin.Register(&CasConnector{
		Config: &CasConnectorConfig{},
	})
}

func (c *CasConnector) Info() plugin.Info {
	info := &util.Info{}
	info.GetInfo(Info)

	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    info.SlugName,
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      info.Author,
		Version:     info.Version,
		Link:        info.Link,
	}
}

// ConnectorLogoSVG returns the logo in SVG format
func (c *CasConnector) ConnectorLogoSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
		<circle cx="12" cy="12" r="10"/>
	</svg>`
}

// ConnectorName returns the name of the connector
func (c *CasConnector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.InfoName)
}

// ConnectorSlugName returns the slug name of the connector
func (c *CasConnector) ConnectorSlugName() string {
	return "cas_connector"
}

// ConnectorSender handles the start endpoint of the connector
func (c *CasConnector) ConnectorSender(ctx *plugin.GinContext, receiverURL string) (redirectURL string) {
	base := strings.TrimRight(c.Config.ServerURL, "/")
	values := url.Values{}
	values.Set("service", receiverURL)
	return fmt.Sprintf("%s/login?%s", base, values.Encode())
}

// ConnectorReceiver handles the callback endpoint of the connector
func (c *CasConnector) ConnectorReceiver(ctx *plugin.GinContext, receiverURL string) (userInfo plugin.ExternalLoginUserInfo, err error) {
	ticket := ctx.Query("ticket")
	if ticket == "" {
		return userInfo, fmt.Errorf("cas: missing ticket in callback")
	}

	base := strings.TrimRight(c.Config.ServerURL, "/")
	values := url.Values{}
	values.Set("service", receiverURL)
	values.Set("ticket", ticket)
	validateURL := fmt.Sprintf("%s/serviceValidate?%s", base, values.Encode())

	resp, err := http.Get(validateURL)
	if err != nil {
		return userInfo, fmt.Errorf("cas: contacting server: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return userInfo, fmt.Errorf("cas: reading response: %w", err)
	}

	var casResp casServiceResponse
	if err := xml.Unmarshal(body, &casResp); err != nil {
		return userInfo, fmt.Errorf("cas: parsing response: %w", err)
	}

	if casResp.Failure != nil {
		return userInfo, fmt.Errorf("cas: authentication failed (%s): %s",
			casResp.Failure.Code, strings.TrimSpace(casResp.Failure.Message))
	}
	if casResp.Success == nil || casResp.Success.User == "" {
		return userInfo, fmt.Errorf("cas: no authenticationSuccess in response")
	}

	username := casResp.Success.User
	email := casResp.Success.attribute("mail", "email", "emailAddress")
	displayName := casResp.Success.attribute("displayName", "cn", "fullName")
	if displayName == "" {
		displayName = username
	}

	userInfo = plugin.ExternalLoginUserInfo{
		// CAS usernames are unique on the CAS server, so they're a safe
		// external ID -- Answer just needs *something* stable per user.
		ExternalID:  username,
		DisplayName: displayName,
		Username:    username,
		// Only set Email if your CAS server's attribute release can be
		// trusted as verified -- see the README before enabling this in
		// a real deployment.
		Email:    email,
		MetaInfo: string(body),
	}
	return userInfo, nil
}
