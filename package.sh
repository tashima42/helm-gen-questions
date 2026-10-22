# SPDX-FileCopyrightText: SUSE LLC
# SPDX-License-Identifier: Apache-2.0
#!/bin/sh

COMMIT=$(git rev-parse --short HEAD)
VERSION=${VERSION:-"0.0.0+snapshot-${COMMIT}"}

# Wrapper needed for macos
function is_gnu_sed() {
  sed --version >/dev/null 2>&1
}

function sed_i_wrapper() {
  if is_gnu_sed; then
    $(which sed) "$@"
  else
    a=()
    for b in "$@"; do
      [[ $b == '-i' ]] && a=("${a[@]}" "$b" "") || a=("${a[@]}" "$b")
    done
    $(which sed) "${a[@]}"
  fi
}

sed_i_wrapper -i "s|^version: \".*\"|version: \"${VERSION}\"|" plugin.yaml
grep -q "^version: \"${VERSION}\"$" plugin.yaml

rm -rf dist
mkdir -p dist/stage
cp -r plugin.yaml LICENSE README.md bin dist/stage/
helm plugin package dist/stage --sign=false --destination dist
cd dist && (sha256sum questions-${VERSION}.tgz || shasum -a 256 questions-${VERSION}.tgz) >checksums.txt
