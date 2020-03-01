#!/usr/bin/env bash

REMOTE_NAME=gitlab
REMOTE_HOST=peon
REMOTE_PORT=30080
GROUP_NAME=golang
REPO_NAME=golang-my-events-example

if [ $(git remote | grep "${REMOTE_NAME}") != "${REMOTE_NAME}" ]; then
    # Add remote url
    git remote add ${REMOTE_NAME} ssh://git@${REMOTE_HOST}:${REMOTE_PORT}/${GROUP_NAME}/${REPO_NAME}.git
fi
