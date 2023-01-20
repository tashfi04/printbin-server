#! /bin/sh

export PRINTBIN_CONFIG_NAME="config-local"
export PRINTBIN_CONFIG_SECRET_NAME="config-secret-local"
export PRINTBIN_CONFIG_PATH="."

export GO111MODULE=on
#export GOARCH="amd64"
#export GOOS="linux"
#export CGO_ENABLED=1

export GOPROXY="https://proxy.golang.org"
export GOSUMDB=sum.golang.org

cmd=$1
sub_cmd=$2

package="github.com/tashfi04/printbin-server"
binary="printbin-server"

rm ${binary}

echo "Building printbin API ..."
GOFLAGS=-mod=vendor go build

echo $cmd
echo $sub_cmd

if [ "$cmd" = "serve" ]; then
  echo "Executing run command..."
  ./${binary} serve -p 8001
  exit
fi

if [ "$cmd" = "migration" ]; then
  if [ "$sub_cmd" = "up" ]; then
    echo "Executing run command..."
    ./${binary} migration up
    exit
  fi

  if [ "$sub_cmd" = "down" ]; then
    echo "Executing run command..."
    ./${binary} migration down
    exit
  fi

  if [ "$sub_cmd" = "reset" ]; then
    echo "Executing run command..."
    ./${binary} migration reset
    exit
  fi
fi

echo "No command specified!"
