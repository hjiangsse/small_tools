#!/bin/bash
## This script use scp download or upload directory between local mechine and
## remote mechine
source ~/.fastdir.conf
hostinfo="${username}@${remotehost}"

dir=''

function get_user_choice() {
    IFS=': '
    read -r -a dirs <<< "$2"
    
    echo -n "The "
    echo -n $1
    echo " directories: "

    echo ----------------------------------------
    for index in "${!dirs[@]}"
    do
        echo "$((index+1)): ${dirs[index]}"
    done
    echo ----------------------------------------
    echo -n "Choice the one your want to interact[e.g. 1]: "
    read choice

    # test if ${rmchoice} is number
    re='^[0-9]+$'
    if ! [[ $choice =~ $re ]] ; then
        echo "error: Not a number" >&2; exit 1
    fi

    # test if ${rmchoice} is in bound 
    if [ "$choice" -gt "$((index+1))" ]; then
        echo "error: Index overflow" >&2; exit 1
    fi

    dir=${dirs[$((choice-1))]}
    echo
}

get_user_choice "remote" "${remotedirs}"
rmdir=$dir
get_user_choice "local" "${localdirs}"
lcdir=$dir

cmd="scp -r ${hostinfo}:${rmdir}/* ${lcdir}"
eval "${cmd}"

