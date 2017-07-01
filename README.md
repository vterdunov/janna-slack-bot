# Janna Slack Bot

This Slack bot is client for [Janna API](https://github.com/vterdunov/janna)

### Setup Slack
- Go to https://YOUR_SLACK_TEAM.slack.com/apps/A0F7YS25R-bots to get to the Bots app page
- Press "Add Configuration"
- Give your bot a name. E.g: "@janna"
- Remember your API Token. It will be used to connect to the Slack chat

### Quick Start
```
docker pull vterdunov/janna-slack-bot
docker run -d --rm \
  --name=janna-slack-bot \
  --restart=always \
  -e SLACK_TOKEN=YOUR_SLACK_TOKEN \
  -e JANNA_API_ADDRESS=http://janna.example.com:4567 \
  vterdunov/janna-slack-bot
```

### Build Docker image
Install Docker 17.05+. Because a multi-stage builds is used.  
`docker build -t janna-slack-bot .`

### Development
Install `go dep` tool:  
`go get -u github.com/golang/dep/cmd/dep`  

Install dependencies:  
`make dep`  

Compilation
`make compile`

Export environment variables:  
```bash
export SLACK_TOKEN=XXXX-123456789123-XXXXXXXXXXXXXXXXXXXXXXXX
export JANNA_API_ADDRESS=http://janna.example.com:4567
```
