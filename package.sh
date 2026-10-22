#!/bin/sh

COMMIT=$(git rev-parse --short HEAD)
VERSION=${VERSION:-"0.0.0+snapshot-${COMMIT}"}
sed -i '' "s/^version: \".*\"/version: \"${VERSION}\"/" plugin.yaml
grep -q "^version: \"${VERSION}\"$" plugin.yaml

rm -rf dist
mkdir -p dist/stage
cp -r plugin.yaml LICENSE README.md bin dist/stage/
helm plugin package dist/stage --sign=false --destination dist
cd dist && (sha256sum questions-${VERSION}.tgz || shasum -a 256 questions-${VERSION}.tgz) >checksums.txt
