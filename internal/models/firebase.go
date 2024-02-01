/*
 * Copyright (C) 2023 akquinet GmbH
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

type NotificationPayload struct {
	Title            string `json:"title,omitempty"`
	Body             string `json:"body,omitempty"`
	BodyLocKey       string `json:"body_loc_key,omitempty"`
	BodyLocArgs      string `json:"body_loc_args,omitempty"`
	Icon             string `json:"icon,omitempty"`
	Tag              string `json:"tag,omitempty"`
	Sound            string `json:"sound,omitempty"`
	Badge            string `json:"badge,omitempty"`
	Color            string `json:"color,omitempty"`
	ClickAction      string `json:"click_action,omitempty"`
	TitleLocKey      string `json:"title_loc_key,omitempty"`
	TitleLocArgs     string `json:"title_loc_args,omitempty"`
	AndroidChannelID string `json:"android_channel_id,omitempty"`
}

type Message struct {
	Data                  interface{}          `json:"data,omitempty"`
	To                    string               `json:"to,omitempty"`
	Notification          *NotificationPayload `json:"notification,omitempty"`
	Priority              string               `json:"priority,omitempty"`
	RegistrationIds       []string             `json:"registration_ids,omitempty"`
	MutableContent        bool                 `json:"mutable_content,omitempty"`
	Condition             string               `json:"condition,omitempty"`
	CollapseKey           string               `json:"collapse_key,omitempty"`
	ContentAvailable      bool                 `json:"content_available,omitempty"`
	RestrictedPackageName string               `json:"restricted_package_name,omitempty"`
	DryRun                bool                 `json:"dry_run,omitempty"`
	TimeToLive            int                  `json:"time_to_live,omitempty"`
}

type TokenDetails struct {
	Application      string `json:"application,omitempty"`
	Platform         string `json:"platform,omitempty"`
	AppSigner        string `json:"appSigner,omitempty"`
	AttestStatus     string `json:"attestStatus,omitempty"`
	AuthorizedEntity string `json:"authorizedEntity,omitempty"`
	ConnectionType   string `json:"connectionType,omitempty"`
	ConnectDate      string `json:"connectDate,omitempty"`
	StatusCode       int
	Error            string                                  `json:"error,omitempty"`
	Rel              map[string]map[string]map[string]string `json:"rel,omitempty"`
}

type Client struct {
	ApiKey  string
	Message *Message
	ApiFCM  string
	ApiIID  string
}
