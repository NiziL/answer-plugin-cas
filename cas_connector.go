/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cas_connector

import (
	"embed"
	"github.com/apache/answer-plugins/cas-connector/i18n"
	"github.com/apache/answer-plugins/util"
	"github.com/apache/answer/plugin"
)

//go:embed info.yaml
var Info embed.FS

type CasConnector struct {
	Config *CasConnectorConfig
}

type CasConnectorConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func init() {
	plugin.Register(&CasConnector{
		Config: &CasConnectorConfig{},
	})
}

func (g *CasConnector) Info() plugin.Info {
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
func (g *CasConnector) ConnectorLogoSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
		<circle cx="12" cy="12" r="10"/>
	</svg>`
}

// ConnectorName returns the name of the connector
func (g *CasConnector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.InfoName)
}

// ConnectorSlugName returns the slug name of the connector
func (g *CasConnector) ConnectorSlugName() string {
	return "cas_connector"
}

// ConnectorSender handles the start endpoint of the connector
func (g *CasConnector) ConnectorSender(ctx *plugin.GinContext, receiverURL string) (redirectURL string) {
	// TODO: Implement OAuth flow start
	// This is a Hello World example - implement your OAuth logic here
	return receiverURL + "?hello=world"
}

// ConnectorReceiver handles the callback endpoint of the connector
func (g *CasConnector) ConnectorReceiver(ctx *plugin.GinContext, receiverURL string) (userInfo plugin.ExternalLoginUserInfo, err error) {
	// TODO: Implement OAuth callback handling
	// This is a Hello World example - implement your OAuth callback logic here
	return plugin.ExternalLoginUserInfo{
		ExternalID:    "hello-world-user",
		DisplayName:   "Hello World User",
		Username:      "helloworld",
		Email:         "hello@example.com",
		Avatar:        "",
		MetaInfo:      "",
	}, nil
}

