#!/bin/bash
set -ex

# run ./deploy/deploy_all.sh at root, for my dev env only.
(
    if [ ! -d backend ]; then
        return -1
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