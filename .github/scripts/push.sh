#!/usr/bin/env bash

source ./.github/common.sh


docker push lucasegp/simple-$project:$version
docker tag lucasegp/simple-$project:$version lucasegp/simple-$project:latest