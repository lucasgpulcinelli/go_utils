#!/usr/bin/env bash

source ./.github/common.sh


docker build -t lucasegp/simple-$project:$version $project \
    -f Dockerfile --build-arg PROJECT=$project