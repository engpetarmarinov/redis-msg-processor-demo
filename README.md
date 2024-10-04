# Distributed golang app using Redis pubsub channels and streams with consumer groups

## Run locally
```shell
sudo docker compose up -d
```

## Run publisher
```shell
export $(grep -v '^#' .env | xargs) && python3 scripts/publisher.py
```
