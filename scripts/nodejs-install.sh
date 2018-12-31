#!/usr/bin/env bash

sudo apt-get update -y
sudo apt-get install build-essential libssl-dev -y
curl -sL https://deb.nodesource.com/setup_11.x | sudo -E bash -
sudo apt-get install -y nodejs