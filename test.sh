#!/bin/bash
for f in `ls -A poc/*`; do
    echo "test $f"
    mkdir testdir
    ./ShellgeiBot -test test_config.json "$f" &&
    [[ "$(ls -A testdir | wc -l)" -eq "0" ]] && echo OK || echo NG
    rm -r testdir
    echo -e "==============================================\n"
done
