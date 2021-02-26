# Test Twitter Bot

## Description
A Twitter bot that tweets once a day if you've laid tefillin on appropriate days.

Based on https://github.com/serverless/examples/tree/master/aws-node-twitter-joke-bot

## Installation
Run `npm i` to install dependencies

## Setup
An environment file is required to contain the secret keys for the Twitter API and the joke API's URL. The naming convention for this file is `[STAGE].env.json`, for example `dev.env.json`.
Copy it from env.json.dist

## Testing
The bot can be invoked manually during development with the following command
```
sls invoke local -f bot
```

## Deployment
- Both the `STAGE` and `REGION` options can be used with this bot. If left out the bot will default to `dev` stage and `eu-west-1` region.

Deploy command
```
sls deploy --stage prod --region us-east-1
```
