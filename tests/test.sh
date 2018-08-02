#!/bin/bash
for i in `ls -d */`; do 
    cd $i
    bats .
    [ $? -eq 0 ] || exit 1
    cd -
done
