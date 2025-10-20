#!/bin/bash
TOKEN="8357335831:AAHVuOVbYhoNDtwum5bAlhkdl6H5Ft21-rs"
curl -s "https://api.telegram.org/bot$TOKEN/getUpdates" | jq '.result[] | .message.chat.id'
