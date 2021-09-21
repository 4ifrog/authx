#!/bin/bash

ROOT_DIR=$(git rev-parse --show-toplevel)
DOCKER_COMPOSE_DIR='/docker/docker-compose'
DOCKER_COMPOSE_FILE='docker-compose.test.yaml'
DATABASES=('mongo')
RETRIES=10
RETRY_WAIT=1
WAIT=3

function get_port() {
  declare -r service=${1}
  declare -ir retries=${2}
  declare -i retry_wait=${3}
  declare -i s=0
  declare -i inc=0

  for _ in $(seq 1 "${retries}"); do
    port="$(docker-compose -f "${DOCKER_COMPOSE_FILE}" ps "${service}" | awk 'NR > 2 { print $NF }' | sed 's/^.*:\(.*\)->.*/\1/')"
    if [ "${port}" == '' ]; then
      sleep "${retry_wait}"
      (( retry_wait += inc++ ))
    else
      break
    fi
  done;

  echo "${port}"
}

function launch_container() {
  declare -r database=${1}
  declare -ir retries=${2}
  declare -i retry_wait=${3}
  declare -ir wait=${4}
  declare -i s=0
  declare -i inc=0

  if docker-compose -f "${DOCKER_COMPOSE_FILE}" ps | awk 'NR > 2 { print $1 }' | grep "${database}" > /dev/null; then
    echo "Container ${database} is already running. All set..."
  else
    # Container is not running, start it.
    echo "Launching ${database}..."
    docker-compose -f "${DOCKER_COMPOSE_FILE}" up "${database}" > /dev/null 2>&1 &
    sleep "${wait}"

    # Wait till it's ready ie. tcp port is open.
    echo "Getting port for ${database}"
    port="$(get_port "${database}" ${RETRIES} ${RETRY_WAIT})"
    if [ "${port}" == '' ]; then
      exit 1
    fi

    printf "Waiting for %s at port %s to be ready for connection" "${database}" "${port}"
    for _ in $(seq 1 "${retries}"); do
      if nc -v -w 1 localhost "${port}" > /dev/null 2>&1; then
        s=0
        break
      else
        s=${?}
        printf '.'
        sleep "${retry_wait}"
        (( retry_wait += inc++ ))
      fi
    done;

    echo
    if [ ${s} -eq 0 ]; then
      echo "Port ${port} is open and ${database} is ready..."

      # One final wait, sometimes server may need extra time to initialize after port is open.
      sleep "${wait}"
      docker-compose -f "${DOCKER_COMPOSE_FILE}" ps "${database}"
    else
      exit ${s}
    fi
  fi
}

function main() {
  cd "$(dirname "${ROOT_DIR}${DOCKER_COMPOSE_DIR}")" > /dev/null || exit 1

  for database in "${DATABASES[@]}"; do
    launch_container "${database}" ${RETRIES} ${RETRY_WAIT} ${WAIT}
  done

  cd - > /dev/null || exit 1
}

main
