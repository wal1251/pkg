#!/bin/sh

which ssh-agent || ( apk add --no-cache openssh-client )
eval $(ssh-agent -s)
echo "$SSH_KEY" | tr -d '\r' | ssh-add - > /dev/null
mkdir -p ~/.ssh
chmod 700 ~/.ssh

touch ~/.ssh/known_hosts
chmod 644 ~/.ssh/known_hosts
echo "$KNOWN_HOSTS" > ~/.ssh/known_hosts

ssh -v "ladmin@$HOST" "mkdir -p ~/rmr-pkg"
scp $COMPOSE_FILE "ladmin@$HOST":~/rmr-pkg/docker-compose.yml
scp "$SECRETS" "ladmin@$HOST":~/rmr-pkg/.env

ssh "ladmin@$HOST" "docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY"
ssh -o ServerAliveInterval=30 "ladmin@$HOST" "docker-compose -f ~/rmr-pkg/docker-compose.yml pull"
ssh -o ServerAliveInterval=30 "ladmin@$HOST" "docker-compose -f ~/rmr-pkg/docker-compose.yml up -d"
ssh -o ServerAliveInterval=30 "ladmin@$HOST" "docker logout $CI_REGISTRY"