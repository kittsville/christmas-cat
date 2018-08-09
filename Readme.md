# Business Cat Bot

A Slack bot to purrcedurally generate Business Cat's daily stand-up message e.g.
```
Y: Curiously pawed client’s exports
T: Disrupt literally anyone’s workflow
```

Inspired by [this](https://twitter.com/kittsville/status/983623739421220864) adorable kitty.

## Setup

1. Create a [Slack App](https://api.slack.com/apps) called Business Cat
2. Create a bot user and install them to your relevant slack channel (e.g. \#standup or \#dev)
3. Note down the Webhook URL e.g. `https://hooks.slack.com/services/T394FRT3/MS45PO2/LvSkPS90We2RAQDt`
4. Set up a server/cloud event (we use AWS Lambda) to run Business Cat when you usually have your stand-up
5. Set `SLACK_TOKEN` in your environment as the Webhook from earlier

If Business Cat's running as a Lambda, you can re-deploy him like so:

1. Build with `GOOS=linux go build main.go`
2. Package with `zip handler.zip ./main`
3. Go to Business Cat's Lambda config
4. Upload zip under _Function code_ header
