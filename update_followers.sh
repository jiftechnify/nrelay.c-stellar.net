#!/bin/bash -exl
export PATH="$HOME/.local/share/mise/shims:$PATH"

cd $(dirname $0)/followers_updater

echo "start updating followers ($(date))"

node index.js
cp whitelist.txt ../resource/whitelist.txt
docker compose exec whitelisted_relay touch /app/plugin/evsifter_whitelist
docker compose exec whitelisted_router touch /app/plugin/evsifter_import_dm_wl

echo "completed"
