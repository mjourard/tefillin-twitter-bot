package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	EnvTwitterConsumerKey       = "TWITTER_CONSUMER_KEY"
	EnvTwitterConsumerSecret    = "TWITTER_CONSUMER_SECRET"
	EnvTwitterAccessTokenKey    = "TWITTER_ACCESS_TOKEN_KEY"
	EnvTwitterAccessTokenSecret = "TWITTER_ACCESS_TOKEN_SECRET"
	EnvStandardTweetFormat      = "STANDARD_TWEET_%d"
	EnvDateOverride             = "DATE_OVERRIDE"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (string, error) {
	consumerKey := os.Getenv(EnvTwitterConsumerKey)
	consumerSecret := os.Getenv(EnvTwitterConsumerSecret)
	accessKey := os.Getenv(EnvTwitterAccessTokenKey)
	accessSecret := os.Getenv(EnvTwitterAccessTokenSecret)
	logger := getLogger(ctx)

	if consumerKey == "" || consumerSecret == "" || accessKey == "" || accessSecret == "" {
		return "", errors.New("missing environment variables")
	}
	logger.Info("Loading tweets")
	tweets := loadTweets(EnvStandardTweetFormat)
	if len(tweets) == 0 {
		return "", errors.New(fmt.Sprintf("could not load tweet list from environment. Env variable format: %s", EnvStandardTweetFormat))
	}

	curTime := time.Now()
	if dateOverride := os.Getenv(EnvDateOverride); dateOverride != "" {
		logger.Infof("Setting a date override to %s", dateOverride)
		t, err := time.Parse("2006-01-02", dateOverride)
		if err != nil {
			return "", errors.New(fmt.Sprintf("DATE_OVERRIDE set but the date value of %s was not valid. Must be of the format yyyy-mm-dd", dateOverride))
		} else {
			curTime = t
		}
	}

	logger.Info("Checking if I should tweet today")
	if !shouldTweetToday(curTime) {
		return "Not going to tweet today", nil
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessKey, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	statusParams := &twitter.StatusUpdateParams{
		Lat:                twitter.Float(38.8977),
		Long:               twitter.Float(77.0365),
		DisplayCoordinates: twitter.Bool(true),
	}

	for idx := 0; idx < len(tweets); {

		tweetContent := tweets[idx]
		idx++
		// Send a Tweet
		logger.Infof("Attempting to tweet with '%s'", tweetContent)
		tweet, resp, err := client.Statuses.Update(tweetContent, statusParams)
		if err != nil {
			if twitErr, ok := err.(twitter.APIError); ok {
				if twitErr.Errors[0].Code == TwitDuplicateStatus {
					logger.Infof("Got duplicate tweet warning with '%s', trying the next one...", tweetContent)
					continue
				}
				logger.Errorf("Unhandled twitter error: %s", twitErr.Error())
				return "Unhandled twitter error code", fmt.Errorf("unhandled twitter error code: %s", twitErr.Error())
			} else {
				return "Failed to send tweet", fmt.Errorf("failed to send tweet. Twitter error: %w", err)
			}
		}

		if resp != nil {
			logger.Infof("Twitter API responded with status of %d", resp.StatusCode)
			switch resp.StatusCode {
			case 200:
				return fmt.Sprintf("Created tweet with id %d. Contents: %s", tweet.ID, tweetContent), nil
			default:
				logger.Warnf("Twitter api responded with code %d", resp.StatusCode)
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(resp.Body)
				if err != nil {
					return "Non 200 status code returned", fmt.Errorf("non-200 status code: %d. unable to read response body: %w", resp.StatusCode, err)
				}
				return "Non 200 status code returned", errors.New(buf.String())
			}
		}
	}
	return "Ran out of tweets to attempt to attempt to tweet", fmt.Errorf("ran out of tweets to attempt to attempt to tweet. Attempted %d tweets", len(tweets))
}

func shouldTweetToday(curTime time.Time) bool {
	switch curTime.Weekday() {
	case time.Saturday:
		//don't tweet on Saturdays because of Shabbos
		return false
	}

	return true
}

func loadTweets(envVarFormat string) []string {
	tweets := make([]string, 0)
	idx := 0
	for {
		tweet := os.Getenv(fmt.Sprintf(envVarFormat, idx))
		if tweet == "" {
			break
		}
		tweets = append(tweets, tweet)
		idx++
	}
	return tweets
}

func getLogger(ctx context.Context) *log.Entry {
	lc, _ := lambdacontext.FromContext(ctx)
	logger := log.WithFields(log.Fields{
		"request_id": lc.AwsRequestID,
	})
	log.SetOutput(os.Stdout)
	return logger
}

func main() {
	lambda.Start(Handler)
}
