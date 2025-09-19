#!/bin/sh

installation_path="$HOME/.local/bin"

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

curl -s -o /tmp/pflags https://raw.githubusercontent.com/chahal-p/pflags/refs/heads/main/pflags
chmod +x /tmp/pflags
if [ "$use_sudo" = "true" ]; then
  sudo mv /tmp/pflags "$installation_path"
else
  mv /tmp/pflags "$installation_path"
fi