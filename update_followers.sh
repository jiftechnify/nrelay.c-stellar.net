#!/bin/bash -ex
cd $(dirname $0)/followers_updater

echo "start updating followers ($(date))"

node index.js
cp whitelist.txt ../resource/whitelist.txt
docker exec strfry_whitelisted touch /app/plugin/evsifter_whitelist
docker exec strfry_whitelisted_router touch /app/plugin/evsifter_import_dm_wl

echo "completed"
