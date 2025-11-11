#!/usr/bin/env bash

APP=docker
CONTAINER_NAME="goldap-test"

LDAP_ADMIN_DN="cn=admin,dc=goldap,dc=sh"
LDAP_ADMIN_PASSWORD="admin123"

LDAP_BASE_DN="dc=goldap,dc=sh"

LDAP_URL="ldap://127.0.0.1:389"

RUNNING_CONTAINERS=$( $APP ps | grep $CONTAINER_NAME)
if [ -z "$RUNNING_CONTAINERS" ]; then
  echo "Starting container with $APP"
  $APP \
    run \
    --rm \
    -d \
    --name $CONTAINER_NAME \
    --hostname "${CONTAINER_NAME}.goldap.sh" \
    -p "127.0.0.1:389:389" \
    -p "127.0.0.1:636:636" \
    -e LDAP_ORGANISATION="goLDAP TUI" \
    -e LDAP_DOMAIN="goldap.sh" \
    -e LDAP_ADMIN_PASSWORD="${LDAP_ADMIN_PASSWORD}" \
    -e LDAP_TLS_VERIFY_CLIENT="never" \
    docker.io/osixia/openldap:1.5.0
else
    echo $CONTAINER_NAME already running ...
    echo ""
fi

$APP ps --latest

ldapsearch -H ldap://127.0.0.1:389 \
  -x \
  -b $LDAP_BASE_DN \
  -D $LDAP_ADMIN_DN \
  -w $LDAP_ADMIN_PASSWORD

