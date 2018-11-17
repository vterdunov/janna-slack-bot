#!/bin/bash

PROJECT_NAME='janna-bot'

set -e

# Do not rebuild/retest image that we already have.
if [ -n "$TRAVIS_TAG" ]; then
  echo "Found release tag"
  docker pull ${DOCKER_USERNAME}/${PROJECT_NAME}:${COMMIT}
  docker tag ${DOCKER_USERNAME}/${PROJECT_NAME}:${COMMIT} ${DOCKER_USERNAME}/${PROJECT_NAME}:${TRAVIS_TAG}
  docker push ${DOCKER_USERNAME}/${PROJECT_NAME}:${TRAVIS_TAG}
  exit 0
fi

# Run linters, tests and build docker image
make

if [ "$TRAVIS_BRANCH" == "master" ]; then
  make push
  make push TAG=latest
fi
