#!/usr/bin/env bash

if [[ -n "$CONTAINER_APP" ]]; then
  echo "Using $CONTAINER_APP"
  APP=$CONTAINER_APP
elif command -v podman &>/dev/null; then
  APP=podman
elif command -v docker &>/dev/null; then
  APP=docker
else
  echo "Could't find podman or docker; please install one before running this script again."
  exit 1
fi

echo ""
echo "Using $APP"

CONTAINER_NAME="goldap-test"

LDAP_ADMIN_DN="cn=admin,dc=goldap,dc=sh"
LDAP_ADMIN_PASSWORD="admin123"

LDAP_BASE_DN="dc=goldap,dc=sh"

HOST_LDAP_PORT=8389
HOST_SLDAP_PORT=8636
LDAP_URL="ldap://127.0.0.1:$HOST_LDAP_PORT"

RUNNING_CONTAINERS=$( $APP ps | grep $CONTAINER_NAME)
if [ -z "$RUNNING_CONTAINERS" ]; then
  echo "Starting container with $APP"
  echo -e "------------------------------\n"

  $APP \
    run \
    --rm \
    -d \
    --name $CONTAINER_NAME \
    --hostname "${CONTAINER_NAME}.goldap.sh" \
    -p "127.0.0.1:$HOST_LDAP_PORT:389" \
    -p "127.0.0.1:$HOST_SLDAP_PORT:636" \
    -e LDAP_ORGANISATION="goLDAP TUI" \
    -e LDAP_DOMAIN="goldap.sh" \
    -e LDAP_ADMIN_PASSWORD="${LDAP_ADMIN_PASSWORD}" \
    -e LDAP_TLS_VERIFY_CLIENT="never" \
    docker.io/osixia/openldap:1.5.0

  sleep 5
  
  SOURCE_PATH="$( dirname ${BASH_SOURCE[0]} )"
  ldapadd -H $LDAP_URL \
    -D $LDAP_ADMIN_DN \
    -w $LDAP_ADMIN_PASSWORD \
    -f ${SOURCE_PATH}/0-ous.ldif

  ldapadd -H $LDAP_URL \
    -D $LDAP_ADMIN_DN \
    -w $LDAP_ADMIN_PASSWORD \
    -f ${SOURCE_PATH}/1-uids.ldif
else
    echo $CONTAINER_NAME already running ...
    echo ""
fi

ldapsearch -H $LDAP_URL \
  -x \
  -b $LDAP_BASE_DN \
  -D $LDAP_ADMIN_DN \
  -w $LDAP_ADMIN_PASSWORD

echo -e "\n+---------------------------------+"
echo "| LDAP URL: $LDAP_URL |"
echo -e "+---------------------------------+\n"

$APP ps --filter "name=$CONTAINER_NAME"
