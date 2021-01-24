#!/usr/bin/env bash

ROOT_DIR=$(git rev-parse --show-toplevel)
DOCKER_COMPOSE_DIR='/docker/docker-compose'
DOCKER_COMPOSE_FILE='docker-compose.test.yaml'
DATABASES=('mongo' 'redis')
RETRIES=5
RETRY_WAIT=3
WAIT=2

function launch_container() {
  declare -r database=${1}
  declare -r retries=${2}
  declare -r retry_wait=${3}
  declare -r wait=${4}
  declare -i s=0

  if docker-compose -f "${DOCKER_COMPOSE_FILE}" ps | awk 'NR > 2 { print $1 }' | grep "${database}" > /dev/null; then
    echo "Container ${database} is already running. All set..."
  else
    # Container is not running, start it.
    echo "Launching ${database}..."
    docker-compose -f "${DOCKER_COMPOSE_FILE}" up "${database}" > /dev/null 2>&1 &
    sleep "${wait}"

    # Wait till it's ready ie. tcp port is open.
    port="$(docker-compose -f "${DOCKER_COMPOSE_FILE}" ps "${database}" | awk 'NR > 2 { print $NF }' | sed 's/^.*:\(.*\)->.*/\1/')"
    printf "Waiting for %s at port %s to be ready for connection" "${database}" "${port}"
    for _ in $(seq 1 "${retries}"); do
      nc -v -w 1 localhost "${port}" > /dev/null 2>&1 && s=0 && break || s=${?} && printf '.' && sleep "${retry_wait}";
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
