#!/bin/sh

go install github.com/chahal-p/pflags/cli/pflags@latest

installation_path="$HOME/.local/bin"

binpath=""

test -n "$GOBIN" && binpath="$GOBIN" || {
  test -n "$GOPATH" && binpath="$GOPATH/bin" || binpath="$HOME/go/bin"
}

test -z "$binpath" && { echo "Could not identify Go bin directory."; exit 1; }

for opt do
  value=$(expr "x$opt" : 'x[^=]*=\(.*\)')
  case "$opt" in
  --installation-path=*) installation_path="$value";
  ;;
  *) echo "Unknown argument: $opt"; exit 1;
  ;;
  esac
done
test -e "$installation_path" || mkdir -p "$installation_path"
test -d "$installation_path" || { echo "'$installation_path' installation path should be a directory"; exit 1; }

cp "$binpath/pflags" "$installation_path"