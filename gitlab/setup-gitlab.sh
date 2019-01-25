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

sh -c "/init-repo.sh $GITLAB_PRIVATE_TOKEN $GITLAB_GROUP_NAME $GITLAB_REPO_NAME"
