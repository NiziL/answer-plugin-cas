package cas_connector

import (
	"encoding/json"

	"github.com/apache/answer-plugins/cas-connector/i18n"
	"github.com/apache/answer/plugin"
)

// ConnectorConfig is persisted by Answer (as JSON) whenever an admin saves
// the plugin's settings page, and reloaded into Connector.Config on startup.
type CasConnectorConfig struct {
	ServerURL string `json:"server_url"`
}

// ConfigFields describes the settings form shown under
// Admin -> Plugins -> CAS Connector.
func (c *CasConnector) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "server_url",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigServerURLTitle),
			Description: plugin.MakeTranslator(i18n.ConfigServerURLDescription),
			Required:    true,
			Value:       c.Config.ServerURL,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
		},
	}
}

// ConfigReceiver is called by Answer whenever the settings form above is
// saved, with the submitted values marshalled back into ConnectorConfig's
// JSON shape.
func (c *CasConnector) ConfigReceiver(config []byte) error {
	newConfig := &CasConnectorConfig{}
	if err := json.Unmarshal(config, newConfig); err != nil {
		return err
	}
	c.Config = newConfig
	return nil
}
