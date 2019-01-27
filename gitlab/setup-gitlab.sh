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

#bash /init-repo.sh $GITLAB_PRIVATE_TOKEN $GITLAB_GROUP_NAME $GITLAB_REPO_NAME $GITLAB_API_PORT
privateToken=$GITLAB_PRIVATE_TOKEN
groupName=$GITLAB_GROUP_NAME
repoName=$GITLAB_REPO_NAME
apiPort=$GITLAB_API_PORT

while true; do
  echo "Waiting http server warms up...";
  httpStatus=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$apiPort/-/health)
  if [[ $httpStatus == 200 ]]; then
    break;
  fi
  sleep 30s;
done;

echo "Add SSH key to root user"
userId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/users?username=root" | jq '.[0].id')"
curl -s --header "PRIVATE-TOKEN: $privateToken" -X POST -F "title=at vagrant" -F "key=$(cat /root/.ssh/id_rsa.pub)" "http://localhost:$apiPort/api/v4/users/$userId/keys" | jq

echo "Verifying if group $groupName already exists..."
groupId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/groups?search=$groupName" | jq '.[0].id')"

echo "Verifying if repository $repoName already exists..."
projectId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/projects/?search=$repoName" | jq '.[0].id')"

if [[ $groupId != "null" && $projectId != "null" ]]; then
  echo "Group '$groupName' and repository '$repoName' already exist. Skiping..."
fi

# Create a group
if [[ $groupId == "null" ]]; then
  echo "Group not found. Creating '$groupName' group..."
  curl -s --header "PRIVATE-TOKEN: $privateToken" -X POST "http://localhost:$apiPort/api/v4/groups?name=$groupName&path=$groupName" | jq
  groupId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/groups?search=$groupName" | jq '.[0].id')"
fi

# Create a repo
if [[ $projectId == "null" ]]; then
  echo "Repository not found. Creating '$repoName' repository..."
  curl -s --header "PRIVATE-TOKEN: $privateToken" -X POST "http://localhost:$apiPort/api/v4/projects?name=$repoName" | jq
  projectId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/projects/?search=$repoName" | jq '.[0].id')"
fi

if [[ $groupId == "null" && $projectId == "null" ]]; then
  echo "Ops! Something went wrong! Group '$groupName' and repository '$repoName' not exist."
fi

# transfer the repo to the group
echo "Transfering '$repoName' repository to '$groupName' group..."
curl -s --header "PRIVATE-TOKEN: $privateToken" -X POST "http://localhost:$apiPort/api/v4/groups/$groupId/projects/$projectId" | jq

runnersToken="$(curl -s --header "PRIVATE-TOKEN: $privateToken" http://localhost:$apiPort/api/v4/projects/$projectId | jq -r '.runners_token')"
runnerId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" --request POST "http://localhost:$apiPort/api/v4/runners" --form "token=$runnersToken" --form "description=docker-runner" --form "tag_list=docker" --form "run_untagged=true" --form "locked=false") | jq -r '.id'"