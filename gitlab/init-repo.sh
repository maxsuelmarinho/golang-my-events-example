#!/bin/bash -x

privateToken=$1
groupName="$2"
repoName="$3"
apiPort=$4

echo "Verifying if group $groupName already exists..."
groupId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/groups?search=$groupName" | jq '.[0].id')"

echo "Verifying if repository $repoName already exists..."
projectId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/projects/?search=$repoName" | jq '.[0].id')"

if [[ $groupId != "null" && $projectId != "null" ]]; then
  echo "Group '$groupName' and repository '$repoName' already exist. Skiping..."
  return 0;
fi

# Create a group
if [[ $groupId == "null" ]]; then
  echo "Group not found. Creating '$groupName' group..."
  curl --header "PRIVATE-TOKEN: $privateToken" -X POST "http://localhost:$apiPort/api/v4/groups?name=$groupName&path=$groupName" | jq
  groupId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/groups?search=$groupName" | jq '.[0].id')"
fi

# Create a repo
if [[ $projectId == "null" ]]; then
  echo "Repository not found. Creating '$repoName' repository..."
  curl --header "PRIVATE-TOKEN: $privateToken" -X POST "http://localhost:$apiPort/api/v4/projects?name=$repoName" | jq
  projectId="$(curl -s --header "PRIVATE-TOKEN: $privateToken" "http://localhost:$apiPort/api/v4/projects/?search=$repoName" | jq '.[0].id')"
fi

if [[ $groupId != "null" || $projectId != "null" ]]; then
  echo "Ops! Something went wrong! Group '$groupName' and repository '$repoName' not exist."
  return 1;
fi

# transfer the repo to the group
echo "Transfering '$repoName' repository to '$groupName' group..."
curl --header "PRIVATE-TOKEN: $privateToken" -X POST "http://localhost:$apiPort/api/v4/groups/$groupId/projects/$projectId" | jq
