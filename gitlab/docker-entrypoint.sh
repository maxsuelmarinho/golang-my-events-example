#!/bin/bash -x

echo "========================"
echo "GITLAB INITIALIZING"
echo "========================"

sh -c "/setup-gitlab.sh" & sh -c "/assets/wrapper"
