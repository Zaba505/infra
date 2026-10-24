#!/usr/bin/env bash
# Checks that every pin of the Dagger version agrees with the Dagger CLI pin
# in .github/workflows/build.yaml: engineVersion in every tracked
# dagger-module.toml, and the version the Getting started page tells the
# operator to install. Build runs it, and so does verify in
# .claude/backlog.json, so a change that moves one pin and not the others
# fails in both places.
#
# Run it from the repository root.
set -euo pipefail

workflow=.github/workflows/build.yaml
page=docs/content/getting-started/prerequisites.md

want=$(sed -n -E 's/^[[:space:]]*DAGGER_VERSION:[[:space:]]*"?([^"[:space:]]+)"?.*/\1/p' "$workflow" | head -1)
if [ -z "$want" ]; then
  echo "::error file=$workflow::cannot read the DAGGER_VERSION pin"
  exit 1
fi

manifests=$(git ls-files -- 'dagger-module.toml' '**/dagger-module.toml')
if [ -z "$manifests" ]; then
  echo "::error::found no tracked dagger-module.toml to check the Dagger pin against"
  exit 1
fi

status=0
while IFS= read -r manifest; do
  pinned=$(sed -n -E 's/^engineVersion[[:space:]]*=[[:space:]]*"v?([^"]+)".*/\1/p' "$manifest")
  if [ "$pinned" != "$want" ]; then
    echo "::error file=$manifest::engineVersion is '$pinned' but Build installs the Dagger CLI at '$want'. Move them together."
    status=1
  fi
done <<<"$manifests"

if ! grep -qF "DAGGER_VERSION=$want " "$page"; then
  echo "::error file=$page::the page does not tell the operator to install the Dagger CLI at '$want'. Move it with the other pins."
  status=1
fi

if [ "$status" -eq 0 ]; then
  echo "every Dagger pin is $want"
fi
exit "$status"
