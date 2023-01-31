# gh-issue-notifier

Selective notifications for GitHub issues

Get notified when issues with a specific label are created.

## How to run

### Docker

Following environment variables are required, a good idea would be to put them in an `.env` file.

```text
SHOUTRRR_URL=<your shoutrrr url>
INTERVAL=10s
OWNER=<repository owner>
REPOSITORY=<repository name>
LABELS=good first issue
```

You can use the provided docker image either with the compose file found in this repository or start the container from the command line via `docker run --rm --env-file .env -v $(pwd)/issues.txt:/issues.txt ghcr.io/4350pchris/gh-issue-notifier`
