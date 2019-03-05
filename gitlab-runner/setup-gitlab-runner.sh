#!/bin/bash -x

echo "========================"
echo "GITLAB RUNNER INITIALIZING"
echo "========================"

projectId=""
runnersToken=""
while true; do
  echo "Waiting http gitlab server warms up...";
  httpStatus=$(curl -s -o /dev/null -w "%{http_code}" $GITLAB_API_URL/-/health)
  if [[ $httpStatus != 200 ]]; then
    sleep 15s;
    continue;
  fi

  echo "Waiting initial project be created...";
  projectId="$(curl -s --header "PRIVATE-TOKEN: $GITLAB_PRIVATE_TOKEN" "$GITLAB_API_URL/api/v4/projects/?search=$GITLAB_REPO_NAME" | jq '.[0].id')"
  if [[ -z "$projectId" || "$projectId" == "null" ]]; then
    sleep 15s;
    continue;
  fi

  echo "Waiting runners token be generated...";
  runnersToken="$(curl -s --header "PRIVATE-TOKEN: $GITLAB_PRIVATE_TOKEN" $GITLAB_API_URL/api/v4/projects/$projectId | jq -r '.runners_token')"
  if [[ ! -z "$runnersToken" && "$runnersToken" != "null" ]]; then
    break;
  fi
  sleep 15s;
done;

gitlab-runner register --non-interactive \
    --url $GITLAB_API_URL \
    --registration-token "$runnersToken" \
    --executor "docker" \
    --docker-image ubuntu:16.04 \
    --description "Docker Runner" \
    --tag-list "docker" \
    --run-untagged \
    --locked="false" \
    --docker-volumes /var/run/docker.sock:/var/run/docker.sock
