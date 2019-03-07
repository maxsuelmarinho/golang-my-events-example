#!/bin/bash -x

echo "================================"
echo "Wait for Gitlab warms up..."
echo "================================"

while true; do
  echo "waiting Gitlab warms up...";
  echo "$GITLAB_PRIVATE_TOKEN $GITLAB_GROUP_NAME $GITLAB_REPO_NAME"
  /opt/gitlab/bin/gitlab-rails -v
  result=$?
  if [[ $result == 0 ]] && [[ ! -z $(/opt/gitlab/bin/gitlab-rails r "user=User.where(id: 1).first; print user") ]]; then
    echo "User found.";
    break;
  fi

  sleep 30s;
done

echo "Creating access token..."
/opt/gitlab/bin/gitlab-rails r "token=PersonalAccessToken.new(user: User.where(id: 1).first, name: 'token teste', token: '$GITLAB_PRIVATE_TOKEN', scopes: Gitlab::Auth::available_scopes); token.save!"

privateToken=$GITLAB_PRIVATE_TOKEN
groupName=$GITLAB_GROUP_NAME
repoName=$GITLAB_REPO_NAME
apiPort=$GITLAB_API_PORT

userId=""
groupId=""
projectId=""
while true; do
  echo "Waiting http server warms up...";
  httpStatus=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$apiPort/-/health)
  if [[ $httpStatus != 200 ]]; then
    sleep 30s;
    continue;
  fi

  keysCount="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/users/$userId/keys" | jq length)";
  if [[ -z $keysCount || $keysCount != "1" ]]; then
    echo "Add SSH key to root user";
    userId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/users?username=root" | jq '.[0].id')";
    curl -s --header "PRIVATE-TOKEN: $privateToken" -X POST -F "title=at vagrant" -F "key=$(cat /root/.ssh/id_rsa.pub)" "http://localhost:$apiPort/api/v4/users/$userId/keys" | jq;
  fi

  echo "Verifying if group $groupName already exists...";
  groupId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/groups?search=$groupName" | jq '.[0].id')"
  if [[ -z $groupId || $groupId == "null" ]]; then
    echo "Group not found. Creating '$groupName' group...";
    curl -s --header "PRIVATE-TOKEN: $privateToken" -X POST "http://localhost:$apiPort/api/v4/groups?name=$groupName&path=$groupName" | jq;
    groupId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/groups?search=$groupName" | jq '.[0].id')";
  fi

  echo "Verifying if repository $repoName already exists...";
  projectId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/projects/?search=$repoName" | jq '.[0].id')";
  if [[ -z $projectId || $projectId == "null" ]]; then
    echo "Repository not found. Creating '$repoName' repository...";
    curl -s --header "PRIVATE-TOKEN: $privateToken" -X POST "http://localhost:$apiPort/api/v4/projects?name=$repoName" | jq;
    projectId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/projects/?search=$repoName" | jq '.[0].id')";
  fi

  if [[ ! -z $keysCount && $keysCount == "1" && ! -z $groupId && $groupId != "null" && ! -z $projectId && $projectId != "null" ]]; then
    break;
  fi

  sleep 30s;
done;

# transfer the repo to the group
echo "Transfering '$repoName' repository to '$groupName' group..."
curl -s --header "PRIVATE-TOKEN: $privateToken" -X POST "http://localhost:$apiPort/api/v4/groups/$groupId/projects/$projectId" | jq

runnersToken="$(curl -s --header "PRIVATE-TOKEN: $privateToken" http://localhost:$apiPort/api/v4/projects/$projectId | jq -r '.runners_token')"

#runnerId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" --request POST "http://localhost:$apiPort/api/v4/runners" --form "token=$runnersToken" --form "description=docker-runner" --form "tag_list=docker" --form "run_untagged=true" --form "locked=false") | jq -r '.id'"

if [[ ! -z $projectId && $projectId != "null" ]]; then
  curl -s --header "PRIVATE-TOKEN: $privateToken" --request POST "http://localhost:$apiPort/api/v4/projects/$projectId/variables" --form "key=DOCKER_USERNAME" --form "value=${DOCKER_USERNAME}" | jq
  curl -s --header "PRIVATE-TOKEN: $privateToken" --request POST "http://localhost:$apiPort/api/v4/projects/$projectId/variables" --form "key=DOCKER_PASSWORD" --form "value=${DOCKER_PASSWORD}" | jq
  curl -s --header "PRIVATE-TOKEN: $privateToken" --request POST "http://localhost:$apiPort/api/v4/projects/$projectId/variables" --form "key=KUBE_TOKEN" --form "value=${KUBE_TOKEN}" | jq
  curl -s --header "PRIVATE-TOKEN: $privateToken" --request POST "http://localhost:$apiPort/api/v4/projects/$projectId/variables" --form "key=KUBE_CA_CERT" --form "value=${KUBE_CA_CERT}" | jq
  curl -s --header "PRIVATE-TOKEN: $privateToken" --request POST "http://localhost:$apiPort/api/v4/projects/$projectId/variables" --form "key=KUBE_SERVER" --form "value=${KUBE_SERVER}" | jq
  curl -s --header "PRIVATE-TOKEN: $privateToken" --request POST "http://localhost:$apiPort/api/v4/projects/$projectId/variables" --form "key=KUBE_CLUSTER" --form "value=${KUBE_CLUSTER}" | jq
fi
