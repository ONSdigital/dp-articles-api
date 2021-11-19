#!/bin/bash -eux

pushd dp-articles-api
  make build
  cp build/dp-articles-api Dockerfile.concourse ../build
popd
