#!/usr/bin/env bash
git add . && git stash
git pull
./start-webhook.sh