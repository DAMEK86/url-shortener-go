#!/usr/bin/env bash

appName='url-shortener'
VERSION=0.0.1

export CGO_ENABLED=0
goflags=""
if [[ "${READ_ONLY:-false}" == "true" ]]; then
    echo "Running in readonly mode"
    goflags="-mod=readonly"
fi

## go-fmt: format go code
function task_go_fmt {
    go fmt ./...
}

## build [OS] [version]: builds the go executable with influx for container usage
function task_build {
  source_version=$(git rev-parse --short HEAD)
  os=${1:-linux}
  build_version=${2:-$VERSION}
  GOOS=${os} go build -a -o ${appName} \
  -trimpath \
  -ldflags="-s -w -X github.com/damek86/url-shortener-go/config.BuildVersion=${build_version} -X github.com/damek86/url-shortener-go/config.SourceVersion=${source_version}" \
   ${goflags} main.go
}

## build-container: builds the container image
function task_build_container {
    task_build "linux" ${VERSION}
    docker build -t damek/${appName}:${VERSION} .
}

## run-container: run the container image
function task_run_container() {
     docker run -p 8080:8080 damek/${appName}:${VERSION}
}

function task_usage {
    echo "Usage: $0"
    sed -n 's/^##//p' <$0 | column -t -s ':' |  sed -E $'s/^/\t/'
}

CMD=${1:-}
shift || true
RESOLVED_COMMAND=$(echo "task_"$CMD | sed 's/-/_/g')
if [ "$(LC_ALL=C type -t $RESOLVED_COMMAND)" == "function" ]; then
    pushd $(dirname "${BASH_SOURCE[0]}") >/dev/null
    $RESOLVED_COMMAND "$@"
else
    task_usage
fi