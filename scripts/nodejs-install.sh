#!/usr/bin/env bash

VERSION=11

sudo apt-get update -y
sudo apt-get install build-essential libssl-dev gcc g++ make -y
curl -sL https://deb.nodesource.com/setup_$VERSION.x | sudo -E bash -
sudo apt-get install -y nodejs
sudo apt autoremove

echo "Node Version: $(node -v)"
echo "NPM Version: $(npm -v)"

sudo npm install -g typescript
sudo npm install -g webpack
sudo npm install -g webpack-cli
sudo npm install -g http-server



