#!/bin/bash

TIME=$(date)
VERSION=$(git tag | head -n 1)
COMMIT=$(git log --oneline --no-notes | awk '{ print $1; }' | head -n 1)

go build -ldflags="-X 'cloudflare-dns/cmd.BuildTime=${TIME}' -X 'cloudflare-dns/cmd.BuildVersion=${VERSION}' -X 'cloudflare-dns/cmd.BuildCommit=${COMMIT}'"