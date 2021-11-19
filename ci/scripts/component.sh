#!/bin/bash -eux

pushd dp-articles-api
  make test-component
popd
