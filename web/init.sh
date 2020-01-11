#!/bin/bash
set -e

sed -i "s|API_BASE|$API_BASE|g" /etc/nginx/web/main.js
sed -i "s|API_BASE|$API_BASE|g" /etc/nginx/conf.d/web.conf
nginx