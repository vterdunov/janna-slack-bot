# Janna Slack Bot

This Slack bot is client for [Janna API](https://github.com/vterdunov/janna)

### Setup Slack
- Go to https://YOUR_SLACK_TEAM.slack.com/apps/A0F7YS25R-bots to get to the Bots app page
- Press "Add Configuration"
- Give your bot a name. E.g: "@janna"
- Remember your API Token. It will be used to connect to the Slack chat

### Quick Start
```
docker pull vterdunov/janna-bot
docker run -d --rm \
  --name=janna-slack-bot \
  --restart=always \
  -e SLACK_TOKEN=YOUR_SLACK_TOKEN \
  -e JANNA_API_ADDRESS=http://janna.example.com \
  vterdunov/janna-bot
```

### Development
- Install Go v1.11+ environment.
- Copy `cp .env.example .env` and change env file.
- Compile and run
```
make run
```

Run `make help` to additional useful commands.

### Build Docker image
```
make docker
```
