#!/bin/bash

# set -ex
set -e

last=$(git describe --tags --abbrev=0)
changes=$(git diff --raw "$last" backend/**/schema/v*.sql)
regex="^:.+[[:space:]]A[[:space:]]+(.*)$"

if [ -z "$changes" ];then
    echo "No schema changes found since $last"
else
    echo "Schema changes found since $last"
    echo ""

    c=0
    while IFS= read -r l; do
        if [[ "$l" =~ $regex ]]; then
            ((c++))
            echo "$c. [${BASH_REMATCH[1]}](../${BASH_REMATCH[1]})"
        fi
        # echo "line-> $l"
    done <<< "$changes"
    echo ""
fi