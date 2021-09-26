#!/bin/bash

# USAGE: ./install.sh [version]
# will install the latest tool executable to your /usr/local/bin

set -e

echo "Installing..."

TMPDIR=${TMPDIR:-"/tmp"}
pushd "$TMPDIR" > /dev/null
  mkdir -p complexity
  pushd complexity > /dev/null;
    distro=$(if [[ "$(uname -s)" == "Darwin" ]]; then echo "osx"; else echo "linux"; fi)
    if [ -n "$1" ]
    then
      echo "Will download and install v$1"
      curl -sSL --fail -o complexity.zip "https://github.com/apiiro/code-complexity/releases/download/v$1/complexity-$1-$distro.zip"
    else
      curl -s --fail https://api.github.com/repos/apiiro/code-complexity/releases/latest | grep "browser_download_url.*$distro.zip" | cut -d : -f 2,3 | tr -d \" | xargs curl -sSL --fail -o complexity.zip
    fi
    unzip complexity.zip
    chmod +x complexity-*
    cp -f complexity-* /usr/local/bin/complexity
  popd > /dev/null
  rm -rf complexity
popd > /dev/null

echo "Done: $(complexity -v)"
