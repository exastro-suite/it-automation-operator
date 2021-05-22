#!/bin/bash -u

BASE_DIR=$(cd $(dirname $0); pwd)
TARGET=$1

kubectl apply -f ${BASE_DIR}/ita1.yml
kubectl apply -f ${BASE_DIR}/ita2.yml
