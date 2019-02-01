#!/bin/bash

set -euo pipefail

# We set the exit status to return 1 so that the analysis script will error if the code does not follow the conventions implied by go
GO111MODULE=on golint --set_exit_status ./...
echo "All found go files passed golint tests"

GO111MODULE=on go vet ./...
echo "All found go files passed govet tests"

find . -type f -name '*.go' -print0 | xargs -0 misspell -error
echo "All found go files passed the misspell tests"

find . -type f -name '*.md' -print0 | xargs -0 misspell -error
echo "All found md files passed the misspell tests"

shellcheck scripts/*.sh
echo "All found script files passed the shellcheck tests"
