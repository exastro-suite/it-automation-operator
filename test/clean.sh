#!/bin/bash -u

BASE_DIR=$(cd $(dirname $0); pwd)
TARGET=$1

${BASE_DIR}/ita-undeploy.sh ${TARGET}
${BASE_DIR}/volume-delete.sh ${TARGET}
