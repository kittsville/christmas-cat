# Christmas Cat Bot

A Slack bot to purrcedurally generate ~~Business~~ Christmas Cat's daily stand-up message e.g.
```
Y: Curiously pawed client’s presents
T: Disrupt literally anyone’s carol singers
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

## Customise

Want work-specific terms to appear in standup messages? You can customise the subjects (e.g. client's, our, etc.)
or objects (e.g. exports, workflow, roadmap, etc.) referred to by Business Cat using the `EXTRA_SUBJECTS` and
`EXTRA_OBJECTS` environment variables. Both take a comma separated list like so:  
`EXTRA_OBJECTS=agenda,churn,intern`
