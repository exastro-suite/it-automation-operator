BASE_DIR=$(cd $(dirname $0); pwd)

kubectl delete -f ${BASE_DIR}/pv.yml
