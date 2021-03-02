# Tefillin Twitter Bot

## Description
A Twitter bot that tweets once a day if you've laid tefillin on appropriate days.

Based on https://github.com/serverless/examples/tree/master/aws-node-twitter-joke-bot

## Installation
* Run `npm i` to install dev dependencies
* Run `go mod vendor` to install go dependencies 

## Setup
An environment file is required to contain the secret keys for the Twitter API and the joke API's URL. The naming convention for this file is `[STAGE].env.json`, for example `dev.env.json`.
Copy it from env.json.dist

## Testing
The bot can be invoked manually during development with the following command
```
# if you haven't built already...
make build 
make localinvoke
```

If you want to override what day the bot thinks it is, you can do so with 
```
make localinvoke other_env="--env DATE_OVERRIDE2021-02-24" 
```
Where the date you want to set it as is of the format `yyyy-mm-dd`

## Deployment
- Both the `STAGE` and `REGION` options can be used with this bot. If left out the bot will default to `dev` stage and `eu-west-1` region.

Deploy command
```
make deploy
```

Set the region and stage:
```
make region=ca-central-1 stage=prod deploy
```
