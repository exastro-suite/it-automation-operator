BASE_DIR=$(cd $(dirname $0); pwd)

oc apply -f ${BASE_DIR}/pv.yml
