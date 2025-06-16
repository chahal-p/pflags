#!/usr/bin/env bash

[ -d "/tmp/pflags" ] || git clone https://github.com/chahal-p/pflags.git /tmp/pflags
pushd /tmp/pflags > /dev/null
git pull

make build
sudo cp ./out/pflags /usr/local/bin/

popd > /dev/null