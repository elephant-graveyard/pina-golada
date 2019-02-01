#!/bin/bash

set -euo pipefail

find . -type f -name '*.go' -print0 | xargs -0 gofmt -s -e -l -w #-s to try to simpifly code, -e to print errors, -w to write improved version to actual file
echo "Formatted all go files using gofmt"

shfmt -s -w -i 2 -ci scripts/*.sh #-s to try to simplify code -w to print it to the file -i 2 to specify intend lenght and -ci to indent switch
echo "Formatted all shell scripts"
