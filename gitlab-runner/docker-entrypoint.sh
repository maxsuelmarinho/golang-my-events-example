#!/bin/bash -x

echo "====================="
echo "TEST"
echo "====================="

#sh -c "/usr/bin/setup-gitlab-runner.sh"
#sh -c "/setup-gitlab-runner.sh" & sh -c "run --user=gitlab-runner --working-directory=/home/gitlab-runner"

exec "/setup-gitlab-runner.sh" & exec /entrypoint "$@"
