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

package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/clients"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/logger"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/metrics"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/models"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"path/filepath"
)

// RealFirebaseClient is the real implementation of FirebaseClient using messaging.Client.
type RealFirebaseClient struct {
	Client *messaging.Client
}

var ClientFirebase FirebaseClient

// FirebaseClient interface includes the methods you use from messaging.Client.
type FirebaseClient interface {
	Send(ctx context.Context, msg *messaging.Message) (string, error)
	// Add other methods you use from messaging.Client
}

// Send is the implementation of the Send method for RealFirebaseClient.
func (c *RealFirebaseClient) Send(ctx context.Context, msg *messaging.Message) (string, error) {
	// Your real implementation here
	response, err := c.Client.Send(ctx, msg)
	if err != nil {
		return "", err
	}
	return response, nil
}

// MockMessagingClient is a mock implementation of FirebaseClient interface.
type MockMessagingClient struct{}

// Send is a mock implementation for sending FCM messages.
func (c *MockMessagingClient) Send(_ context.Context, _ *messaging.Message) (string, error) {
	return "projects/timref/messages/1638949066784362/", nil
}

func SetupFirebase() (*firebase.App, context.Context, *messaging.Client, error) {
	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("./develop/serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	//Messaging client
	client, _ := app.Messaging(ctx)

	return app, ctx, client, err
}

func CreateFBCClient() error {
	app, _, _, err := SetupFirebase()
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	ClientFirebase = client
	return err
}

func SendToToken(body models.RequestBody, c models.Config) (models.MatrixResponse, error) {
	var err error
	var response string
	var rejected models.MatrixResponse
	ctx := context.Background()
	message := &messaging.Message{}

	logger.Logger.Debug("Got this", zap.Any("Body:", body))
	var NotificationTitle string
	if body.Notification.RoomName != "" {
		NotificationTitle = body.Notification.RoomName
	} else {
		NotificationTitle = "TIMessenger"
	}

	for newDevice := range body.Notification.Devices {
		if clients.GetEnabledClients(body.Notification.Devices[newDevice].AppId, c) {

			registrationToken := body.Notification.Devices[newDevice].PushKey

			message = &messaging.Message{
				Notification: &messaging.Notification{
					Title: NotificationTitle,
					Body:  body.Notification.Sender,
				},
				Data: map[string]string{
					"event_id":            body.Notification.EventID,
					"room_id":             body.Notification.RoomID,
					"type":                body.Notification.Type,
					"sender":              body.Notification.Sender,
					"sender_display_name": body.Notification.SenderDisplayName,
					"room_name":           body.Notification.RoomName,
					"room_alias":          body.Notification.RoomAlias,
					"prio":                body.Notification.Prio,
				},
				Android: &messaging.AndroidConfig{
					Priority: "high",
					Notification: &messaging.AndroidNotification{
						Title: body.Notification.RoomName,
						Body:  body.Notification.RoomAlias,
						Sound: body.Notification.Devices[newDevice].Tweaks.Sound,
					},
				},
				APNS: &messaging.APNSConfig{
					Payload: &messaging.APNSPayload{
						Aps: &messaging.Aps{
							Sound:            body.Notification.Devices[newDevice].Tweaks.Sound,
							MutableContent:   true,
							ContentAvailable: true,
						},
					},
				},
				Token: registrationToken,
			}

			response, err = ClientFirebase.Send(ctx, message)
			metrics.TotalFirbaseSentRequests.Inc()
			if err != nil {
				metrics.TotalFirebaseRejectedRequestsMetric.Inc()
				rejected.Rejected = append(rejected.Rejected, registrationToken)
				logger.Logger.Error(err.Error())
			}
		} else {
			logger.Logger.Warn("ApplicationID not recognized", zap.String("ApplicationId", body.Notification.Devices[newDevice].AppId))
			metrics.TotalWrongDeviceIdsProcessedMetric.Inc()
			registrationToken := body.Notification.Devices[newDevice].PushKey
			rejected.Rejected = append(rejected.Rejected, registrationToken)
		}
	}

	metrics.TotalFirebaseRequestsMetric.Inc()
	logger.Logger.Info("Successfully sent message:", zap.String("Response:", response), zap.Strings("Rejected:", rejected.Rejected))
	return rejected, nil

}
