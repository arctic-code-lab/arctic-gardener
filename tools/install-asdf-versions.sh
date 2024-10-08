#!/bin/bash
if [ ! -f .tool-versions ]; then
  echo ".tool-versions file not found!"
  exit 1
fi

cat .tool-versions | while read -r version; do
  ret=$(asdf where $version)
  ex=$?
  if [ $ex -eq 0 ]; then
    echo "$version is already installed. Reshimming"
    asdf reshim $version
  else
    echo "$version is not installed. Trying to install..."
    asdf install $version
    ex=$?
    if [ $ex -ne 0 ]; then
      echo "Error installing $version"
      exit $ex
    fi
  fi
done
