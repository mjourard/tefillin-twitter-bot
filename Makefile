.PHONY: build clean deploy test remove localinvoke
region ?= us-east-1
stage ?= dev
sls ?= ./node_modules/.bin/sls

define GetFromPkg
$(shell node -p "require('./dev.env.json').$(1)")
endef
tweet0 ?= $(call GetFromPkg,STANDARD_TWEET_0)
tweet1 ?= $(call GetFromPkg,STANDARD_TWEET_1)
tweet2 ?= $(call GetFromPkg,STANDARD_TWEET_2)
tweet3 ?= $(call GetFromPkg,STANDARD_TWEET_3)
tweet4 ?= $(call GetFromPkg,STANDARD_TWEET_4)
key := $(call GetFromPkg,TWITTER_CONSUMER_KEY)
secret := $(call GetFromPkg,TWITTER_CONSUMER_SECRET)
accessKey := $(call GetFromPkg,TWITTER_ACCESS_TOKEN_KEY)
accessSecret := $(call GetFromPkg,TWITTER_ACCESS_TOKEN_SECRET)

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/tweet main.go twittererrorcodes.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
ifeq ("$(wildcard $(stage).env.json)","")
	@echo "$(stage).env.json does not exist. Create it before deploying"
else
	$(sls) deploy --region $(region) --stage $(stage) --verbose
endif

test:
	go test

remove:
	$(sls) remove

localinvoke:
	$(sls) invoke local -f bot --env TWITTER_CONSUMER_KEY="$(key)" \
                               --env TWITTER_CONSUMER_SECRET="$(secret)" \
                               --env TWITTER_ACCESS_TOKEN_KEY="$(accessKey)" \
                               --env TWITTER_ACCESS_TOKEN_SECRET="$(accessSecret)" \
                               --env STANDARD_TWEET_0="$(tweet0)" \
                               --env STANDARD_TWEET_1="$(tweet1)" \
                               --env STANDARD_TWEET_2="$(tweet2)" \
                               --env STANDARD_TWEET_3="$(tweet3)" \
                               --env STANDARD_TWEET_4="$(tweet4)" \
                               --log \
                               $(other_env)
