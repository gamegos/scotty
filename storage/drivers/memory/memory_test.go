package memory

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/gamegos/scotty/storage"
)

var appID = "fritestapp"
var channelID = "someRandomchannelID"
var subscriberIDs = []string{"sub_foo", "sub_bar"}

var stg storage.Storage

func TestMain(m *testing.M) {
	stg = New()
	os.Exit(m.Run())
}

func TestPutApp(t *testing.T) {

	app := &storage.App{
		ID: appID,
		GCM: storage.GCMConfig{
			APIKey:    "apikey",
			ProjectID: "projectid",
		},
	}

	err := stg.PutApp(app)

	if err != nil {
		t.Error(err)
	}
}

func TestGetApp(t *testing.T) {

	_, err := stg.GetApp(appID)

	if err != nil {
		t.Error(err)
	}
}

func TestAddChannel(t *testing.T) {

	err := stg.AddChannel(appID, channelID)

	if err != nil {
		t.Error(err)
	}
}

func TestAddSubscriber(t *testing.T) {

	err := stg.AddSubscriber(appID, channelID, subscriberIDs)

	if err != nil {
		t.Error(err)
	}
}

func TestGetChannelSubscribers(t *testing.T) {

	receivedSubscribers, err := stg.GetChannelSubscribers(appID, channelID)

	if err != nil {
		t.Error(err)
	}

	sort.Strings(receivedSubscribers)
	sort.Strings(subscriberIDs)

	if !reflect.DeepEqual(receivedSubscribers, subscriberIDs) {
		t.Error("Subscriber IDs does not match.")
	}
}

func TestAddSubscriberDevice(t *testing.T) {

	device := storage.Device{
		Platform:  "gcm",
		Token:     "footoken",
		CreatedAt: 1436546411,
	}

	for _, subscriberID := range subscriberIDs {
		err := stg.AddSubscriberDevice(appID, subscriberID, &device)

		if err != nil {
			t.Error(err)
		}
	}
}

func TestUpdateDeviceToken(t *testing.T) {

	for _, subscriberID := range subscriberIDs {
		err := stg.UpdateDeviceToken(appID, subscriberID, "footoken", "bartoken")

		if err != nil {
			t.Error(err)
		}
	}
}

func TestGetSubscriberDevices(t *testing.T) {

	expectedDevice := storage.Device{
		Platform:  "gcm",
		Token:     "bartoken",
		CreatedAt: 1436546411,
	}

	devices, err := stg.GetSubscriberDevices(appID, subscriberIDs[0])

	if err != nil {
		t.Error(err)
	}

	device := devices[0]

	if !reflect.DeepEqual(expectedDevice, *device) {
		t.Error("Device data does not match.")
	}
}

func TestDeleteChannel(t *testing.T) {

	err := stg.DeleteChannel(appID, channelID)

	if err != nil {
		t.Error(err)
	}
}
