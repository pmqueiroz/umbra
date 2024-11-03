: "${VERSION:=development}"

build_umbra() {
  local platform=$1
  local arch=$2

  if [ -z "$platform" ] || [ -z "$arch" ]; then
    echo "Error: platform and arch must be specified"
    exit 1
  fi

  rm -rf build/${platform}-${arch}

  echo "Building umbra-${VERSION} binary to ${platform}-${arch}"

  go build -ldflags="-X main.Version=${VERSION}" -o build/${platform}-${arch}/bin/umbra

  echo "Copying built in libs"

  mkdir -p build/${platform}-${arch}/lib/
  cp -r lib/* build/${platform}-${arch}/lib/

  echo "Copying LICENSE, README"

  cp LICENSE README.md build/${platform}-${arch}/
}

build_umbra $@
