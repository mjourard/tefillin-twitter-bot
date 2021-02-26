package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"os"
	"time"
)

const (
	EnvTwitterConsumerKey       = "TWITTER_CONSUMER_KEY"
	EnvTwitterConsumerSecret    = "TWITTER_CONSUMER_SECRET"
	EnvTwitterAccessTokenKey    = "TWITTER_ACCESS_TOKEN_KEY"
	EnvTwitterAccessTokenSecret = "TWITTER_ACCESS_TOKEN_SECRET"
	EnvStandardTweet            = "STANDARD_TWEET"
	EnvDateOverride             = "DATE_OVERRIDE"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (string, error) {
	consumerKey := os.Getenv(EnvTwitterConsumerKey)
	consumerSecret := os.Getenv(EnvTwitterConsumerSecret)
	accessKey := os.Getenv(EnvTwitterAccessTokenKey)
	accessSecret := os.Getenv(EnvTwitterAccessTokenSecret)
	tweetContent := os.Getenv(EnvStandardTweet)

	if consumerKey == "" || consumerSecret == "" || accessKey == "" || accessSecret == "" || tweetContent == "" {
		return "", errors.New("missing environment variables")
	}

	curTime := time.Now()
	if dateOverride := os.Getenv(EnvDateOverride); dateOverride != "" {
		t, err := time.Parse("2006-01-02", dateOverride)
		if err != nil {
			return "", errors.New(fmt.Sprintf("DATE_OVERRIDE set but the date value of %s was not valid. Must be of the format yyyy-mm-dd", dateOverride))
		} else {
			curTime = t
		}
	}

	if !shouldTweetToday(curTime) {
		return "Not going to tweet today", nil
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessKey, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	tweet, resp, err := client.Statuses.Update(tweetContent, nil)
	if err != nil {
		return "Failed to send tweet", err
	}
	if resp != nil {
		if resp.StatusCode != 200 {
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			var newStr string
			if err != nil {
				newStr = "Unable to read response body"
			} else {
				newStr = buf.String()
			}
			return "Non 200 status code returned", errors.New(newStr)
		}
	}

	return fmt.Sprintf("Created tweet with id %d", tweet.ID), nil
}

func shouldTweetToday(curTime time.Time) bool {
	switch curTime.Weekday() {
	case time.Friday:
		return false
	case time.Saturday:
		return false
	}

	return true
}

func main() {
	lambda.Start(Handler)
}
