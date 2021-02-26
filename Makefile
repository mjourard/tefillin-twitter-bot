.PHONY: build clean deploy test remove localinvoke
region ?= us-east-1
stage ?= dev
sls ?= ./node_modules/.bin/sls

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/tweet main.go

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
	$(sls) invoke local -f bot --env TWITTER_CONSUMER_KEY=a \
                               --env TWITTER_CONSUMER_SECRET=a \
                               --env TWITTER_ACCESS_TOKEN_KEY=a \
                               --env TWITTER_ACCESS_TOKEN_SECRET=a \
                               --env STANDARD_TWEET=ayylmao \
                               $(other_env)
