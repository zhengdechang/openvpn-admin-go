#!/bin/bash
set -e

export NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}

pm2 start ./pm2.json --no-daemon
