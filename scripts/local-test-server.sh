#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<EOF
Usage: $0 <start|status|stop>

  start   Start and seed the LDAP server, or show its status if already running
  status  Show the LDAP server container status
  stop    Stop the LDAP server and remove its container (--rm)
  -h, --help  Show this help
EOF
}

if [[ $# -ne 1 ]]; then
  usage >&2
  exit 1
fi

case "$1" in
  -h|--help)
    usage
    exit 0
    ;;
  start|status|stop) ;;
  *)
    usage >&2
    exit 1
    ;;
esac

if [[ -n "${CONTAINER_APP:-}" ]]; then
  APP=$CONTAINER_APP
elif command -v podman &>/dev/null; then
  APP=podman
elif command -v docker &>/dev/null; then
  APP=docker
else
  echo "Couldn't find podman or docker; please install one before running this script again." >&2
  exit 1
fi

echo "Using $APP"

CONTAINER_NAME="goldap-test"

LDAP_ADMIN_DN="cn=admin,dc=goldap,dc=sh"
LDAP_ADMIN_PASSWORD="admin123"

HOST_LDAP_PORT=8389
HOST_SLDAP_PORT=8636
LDAP_URL="ldap://127.0.0.1:$HOST_LDAP_PORT"

show_status() {
  "$APP" ps --filter "name=$CONTAINER_NAME"
}

case "$1" in
  status)
    show_status
    exit 0
    ;;
  stop)
    "$APP" stop "$CONTAINER_NAME"
    echo "Container stopped, bye!"
    exit 0
    ;;
esac

RUNNING_CONTAINERS=$("$APP" ps --filter "name=$CONTAINER_NAME" --format '{{.Names}}')
if grep -Fxq "$CONTAINER_NAME" <<< "$RUNNING_CONTAINERS"; then
  show_status
else
  if ! command -v ldapadd &>/dev/null; then
    echo "Couldn't find ldapadd; install LDAP client tools before starting the server." >&2
    exit 1
  fi

  echo "Starting container with $APP"

  "$APP" \
    run \
    --rm \
    -d \
    --name "$CONTAINER_NAME" \
    --hostname "${CONTAINER_NAME}.goldap.sh" \
    -p "127.0.0.1:$HOST_LDAP_PORT:389" \
    -p "127.0.0.1:$HOST_SLDAP_PORT:636" \
    -e LDAP_ORGANISATION="goLDAP TUI" \
    -e LDAP_DOMAIN="goldap.sh" \
    -e LDAP_ADMIN_PASSWORD="${LDAP_ADMIN_PASSWORD}" \
    -e LDAP_TLS_VERIFY_CLIENT="never" \
    docker.io/osixia/openldap:1.5.0

  sleep 5
  
  SOURCE_PATH="$(dirname "${BASH_SOURCE[0]}")"
  ldapadd -H "$LDAP_URL" \
    -D "$LDAP_ADMIN_DN" \
    -w "$LDAP_ADMIN_PASSWORD" \
    -f "${SOURCE_PATH}/0-ous.ldif"

  ldapadd -H "$LDAP_URL" \
    -D "$LDAP_ADMIN_DN" \
    -w "$LDAP_ADMIN_PASSWORD" \
    -f "${SOURCE_PATH}/1-uids.ldif"

  URL_INFO="LDAP URL: $LDAP_URL"
  length=${#URL_INFO}
  printf -v LINE "%$((${#URL_INFO} + 2))s"
  TOP="${LINE// /─}"
  printf "\n┎${TOP}┐\n"
  echo "┃ LDAP URL: $LDAP_URL │"
  BTM="${LINE// /━}"
  printf "┗${BTM}┙\n"
  show_status
fi
