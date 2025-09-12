#!/bin/sh

go install github.com/chahal-p/pflags/cli/pflags@latest

installation_path="$HOME/.local/bin"

binpath=""

test -n "$GOBIN" && binpath="$GOBIN" || {
  test -n "$GOPATH" && binpath="$GOPATH/bin" || binpath="$HOME/go/bin"
}

test -z "$binpath" && { echo "Could not identify Go bin directory."; exit 1; }

use_sudo="false"

for opt do
  value=$(expr "x$opt" : 'x[^=]*=\(.*\)')
  case "$opt" in
  --installation-path=*) installation_path="$value"
  ;;
  --sudo) use_sudo="true"
  ;;
  *) echo "Unknown argument: $opt"; exit 1
  ;;
  esac
done
test -e "$installation_path" || mkdir -p "$installation_path"
test -d "$installation_path" || { echo "'$installation_path' installation path should be a directory"; exit 1; }

if [ "$use_sudo" = "true" ]; then
  sudo cp "$binpath/pflags" "$installation_path"
else
  cp "$binpath/pflags" "$installation_path"
fi