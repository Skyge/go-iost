#!/bin/bash

NAME=$(git rev-parse --short HEAD)
if [ -n "$1" ]
then
    NAME=$1
fi

cd build/k8s/

echo "Generate iserver config"
cd iserver-config/
go run genconfig.go -c $NAME -m 7 -s 2
cd -

echo "Generate iserver config"
cd itest-config/
export DYLD_LIBRARY_PATH=${GOPATH}/src/github.com/iost-official/go-iost/vm/v8vm/v8/libv8/_darwin_amd64
go run genconfig.go -c $NAME -s "7,8"
cd -

echo "Create test cluster $NAME in k8s"
kubectl create namespace $NAME
kubectl create configmap iserver-config --from-file=iserver-config -n $NAME
kubectl create -f iserver.yaml -n $NAME
kubectl create configmap itest-config --from-file=itest-config -n $NAME
kubectl create -f itest.yaml -n $NAME
