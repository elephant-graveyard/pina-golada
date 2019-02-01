#!/usr/bin/env bash

# Copyright © 2019 The Homeport Team
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

set -euo pipefail

Debug() {
    echo "\\033[37m${1}\\033[0m"
}

Error() {
    echo "\\033[31m${1}\\033[0m"
}

WORK_DIR="./target"
if [ -d "${WORK_DIR}" ]; then
    rm -r -f "${WORK_DIR}"
    Debug "Deleted previous working directory"
fi

ASSET_DIR="./assets"
if ! [ -d "${ASSET_DIR}" ]; then
    Error "Could not find asset directory ${ASSET_DIR}"
    exit 1
fi 

mkdir "${WORK_DIR}"
Debug "Created new working directory"

cd "${ASSET_DIR}"
for DIR in $(find . -mindepth 1 -maxdepth 1 -type d); do 
    DIR=${DIR//\.\//}
    tar -zcf "../target/${DIR}.tar.gz" "${DIR}" 
    echo "\\033[1;31m${ASSET_DIR}/${DIR}\\033[0m \\033[93m▶▶ ${WORK_DIR}/${DIR}.tar.gz\\033[0m"
done


