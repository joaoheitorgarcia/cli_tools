#!/usr/bin/env bash

exportFolder=$1
if [ "$exportFolder" == "" ]; then
  exportFolder="$HOME/.local/bin/"
fi

dirs=(*/)
for folder in "${dirs[@]}"; do
  if [ "$folder" != "build/" ]; then

    cd "$folder" || return
    go test || { echo ""; echo "Tests failed"; exit 1; }
    go build -o ../build || { echo ""; echo "Build Failed"; exit 1; }
    cd ..

  fi
done

echo ""; echo "ALL BUILDS FINISHED COPYING INTO $exportFolder"

buildFolder=(build/*)
for bin in "${buildFolder[@]}"; do
  cp "$bin" "$exportFolder"
done

echo "ALL BINARIES COPIED TO $exportFolder"