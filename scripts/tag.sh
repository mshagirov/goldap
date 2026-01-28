#!/usr/bin/env bash

HELP_TEXT=$(cat <<EOF
 Usage:\n
  \033[1;32m$0\033[0m <VERSION>\n
\n
E.g.: for v1.2.3-COMMIT_HASH\n
  \033[1;32m$0\033[0m 1.2.3\n
EOF
)

if [ "$#" -ne 1 ]; then
  echo "Error using release-tag.sh: one input arg is required." >&2
  echo -e $HELP_TEXT >&2
  exit 1
fi

if [[ $1 == "-h" || $1 == "--help" ]]; then
  echo -e $HELP_TEXT
  exit 0
fi

RELEASE=$(git rev-parse --short=10 HEAD)
VERSION_NUMBER=$1

TAGNAME="v${VERSION_NUMBER}-${RELEASE}"

echo "setting git tag to $TAGNAME"

git tag -a $TAGNAME -m "relase version set to $TAGNAME" 

