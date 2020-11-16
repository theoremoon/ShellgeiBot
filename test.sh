#!/bin/bash
for f in poc/.[!\.]* poc/*; do
    echo "test $f"
    mkdir testdir
    ./ShellgeiBot -test test_config.json "$f" &&
    [[ -z "$(ls -A testdir)" ]] && echo OK || echo NG
    rm -r testdir
    echo -e "==============================================\n"
done
