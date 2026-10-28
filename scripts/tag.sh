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
V=$1

TAGNAME="v${V}-${RELEASE}"

echo "Running go format, vet, tests ..."
go fmt ./... && go vet ./... && go test ./...

echo "Setting git tag to $TAGNAME"
git tag -a $TAGNAME -m "relase version set to $TAGNAME" 

echo pushing to origin
git push origin $TAGNAME

if ! command -v gh >/dev/null 2>&1; then
    echo "gh not found, please install GitHub CLI to upload release"
fi

echo "Building static binaries ..."
rm -rf dist && mkdir dist

for target in linux/amd64 linux/arm64 windows/amd64; do
  os=${target%/*}; arch=${target#*/}
  name=goldap-$V-$os-$arch
  ext=""; [ "$os" = windows ] && ext=.exe

  mkdir -p dist/$name
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch \
    go build -trimpath -ldflags "-s -w" -o dist/$name/goldap$ext .
  cp LICENSE README.md dist/$name/
done

file dist/*/goldap*

cd dist
tar -czf goldap-$V-linux-amd64.tar.gz goldap-$V-linux-amd64
tar -czf goldap-$V-linux-arm64.tar.gz goldap-$V-linux-arm64
zip -qr  goldap-$V-windows-amd64.zip  goldap-$V-windows-amd64

shasum -a 256 *.tar.gz *.zip > checksums.txt

cd ..

gh release create "$TAGNAME" \
  dist/*.tar.gz dist/*.zip dist/checksums.txt \
  --title "$TAG" \
  --generate-notes

gh release view "$TAG"
