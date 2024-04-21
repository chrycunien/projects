#!/bin/bash

# brew install md5sha1sum
FHASH=`md5sum $1`
while true; do
  NHASH=`md5sum $1`
  if [[ "$NHASH" != "$FHASH" ]]; then
    ./mdp -file $1
    FHASH=$NHASH
  fi
done