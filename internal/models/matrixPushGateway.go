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

type RequestBody struct {
	Notification Notification `json:"notification"`
}
type Notification struct {
	Content           Content `json:"content,omitempty"`
	Counts            Counts  `json:"counts,omitempty"`
	Devices           []Device
	EventID           string `json:"event_id,omitempty"`
	Prio              string `json:"prio,omitempty"`
	RoomAlias         string `json:"room_alias,omitempty"`
	RoomID            string `json:"room_id,omitempty"`
	RoomName          string `json:"room_name,omitempty"`
	Sender            string `json:"sender,omitempty"`
	SenderDisplayName string `json:"sender_display_name,omitempty"`
	Type              string `json:"type,omitempty"`
	UserIsTarget      bool   `json:"user_is_target,omitempty"`
}

type Counts struct {
	MissedCalls int `json:"missed_calls,omitempty"`
	Unread      int `json:"unread,omitempty"`
}

type Device struct {
	AppId     string `json:"app_id"`
	Data      PusherData
	PushKey   string `json:"pushkey"`
	PushKeyTS int    `json:"pushkey_ts,omitempty"`
	Tweaks    Tweaks `json:"tweaks,omitempty"`
}

type PusherData struct {
}
type Tweaks struct {
	Sound string `json:"sound,omitempty"`
}

type Content struct {
	Body    string `json:"body,omitempty"`
	MSGType string `json:"msgtype,omitempty"`
}

type MatrixResponse struct {
	Rejected []string `json:"rejected,omitempty"`
}
