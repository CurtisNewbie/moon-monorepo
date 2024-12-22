#!/bin/bash

#set -ex
set -e

if ! ping "curtisnewbie.com" -q -t 1; then
    echo "This script is written for the maintainer only, you are not supposed to run it"
    exit 1
fi

# run ./deploy/deploy_all.sh at root, for my dev env only.
(
    if [ ! -d backend ]; then
        exit -1
    fi

    for r in $(ls ./backend);
    do
        p="./backend/$r"
        echo "$p"
        (
            cd "$p"
            if [ -f "deploy" ]; then
                ./deploy
            fi
        )
    done

    cd "./frontend/moon"
    if [ -f "deploy-nginx.sh" ]; then
        ./deploy-nginx.sh
    fi
)
