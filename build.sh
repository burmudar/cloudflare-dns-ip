#!/bin/bash

TIME=$(date)
VERSION=""
COMMIT=$(git log --oneline --no-notes | awk '{ print $1; }' | head -n 1)

go build -ldflags="-X 'cloudflare-dns/cmd.BuildTime=${TIME}' -X 'cloudflare-dns/cmd.BuildVersion=${VERSION}' -X 'cloudflare-dns/cmd.BuildCommit=${COMMIT}'"