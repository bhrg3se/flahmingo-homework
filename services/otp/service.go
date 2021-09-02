package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/bhrg3se/flahmingo-homework/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func startService(config utils.Config) {

	//initialise pub sub client
	absPath, err := filepath.Abs("/etc/flahmingo/key.json")
	if err != nil {
		logrus.Errorf("google key not found: %v", err)
	}
	opt := option.WithCredentialsFile(absPath)
	psClient, err := pubsub.NewClient(context.Background(), config.GoogleCloud.ProjectID, opt)
	if err != nil {
		panic(err)
	}

	sub := psClient.Subscription("verification-sub")

	logrus.Info("waiting for PubSub messages")
	// handle received message
	err = sub.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		receiverNumber := message.Attributes["PHONE_NUMBER"]
		otp := message.Attributes["OTP"]
		msg := fmt.Sprintf("Your one time passoword is: %s", otp)

		//send receive message through twilio sms api
		err = sendSMS(config.Twilio.AccountSID, config.Twilio.AuthToken, msg, config.Twilio.PhoneNumber, receiverNumber)
		if err != nil {
			logrus.Error(err)
		}
		message.Ack()
	})

	if err != nil {
		logrus.Error(err)
	}

}

// sendSMS sends SMS using Twilio's API
func sendSMS(accountSID, authToken, msg, twilioPhoneNumber, phoneNumber string) error {

	v := url.Values{}
	v.Set("To", phoneNumber)
	v.Set("From", twilioPhoneNumber)
	v.Set("Body", msg)
	rb := *strings.NewReader(v.Encode())
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID)

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send sms: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("got failed response from twilio: %s", resp.Status)
	}
	return nil
}
