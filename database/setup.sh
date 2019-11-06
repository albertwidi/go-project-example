#!/bin/bash

export SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export DATABASE_SCHEMA_DIR=./schema

# we need to specify the database configuration for pop by sending additional parameter
# example: ./setup.sh create database.yml
create() {
    if [ -z $1 ]; then
        echo "database configuration cannot be empty"
        exit 1;
    fi

    soda create -a --config $1
    if [ $? -ne 0 ]; then
        echo "failed to create database with config $1";
        exit 1;
    fi
}

# example: ./setup.sh drop database.yml
drop() {
    if [ -z $1 ]; then
        echo "database configuration cannot be empty"
        exit 1;
    fi

    soda drop -a --config $1
    if [ $? -ne 0 ]; then
        echo "failed to drop database with config $1";
    fi
}

# example: ./setup.sh generate user users
generate() {
    if [ -z $1 ]; then
        echo "database configuration cannot be empty"
        exit 1;
    fi

    if [ -z $2 ]; then
        echo "database name cannot be empty"
        exit 1;
    fi

    if [ -z $3 ]; then 
        echo "migration name cannot be empty"
        exit 1;
    fi

    cd ${SCRIPT_DIR}
    soda generate sql -c $1 -e $2 -p ${DATABASE_SCHEMA_DIR}/$2 $3
    if [ $? -ne 0 ]; then
        echo "failed to generate new sql migration schema"
        cd -
        exit 1;
    fi
    cd -
}

migrate() {
    if [ -z $1 ]; then
        echo "database configuration cannot be empty"
        exit 1;
    fi

    if [ -z $2 ]; then
        echo "database name cannot be empty"
        exit 1;
    fi

    if [ -z $3 ]; then 
        echo "migration type cannot be empty. must be up/down"
        exit 1;
    fi

    cd ${SCRIPT_DIR}
    soda migrate $3 -d -c $1 -e $2 -p ${DATABASE_SCHEMA_DIR}
    if [ $? -ne 0 ]; then
        echo "failed to generate new sql migration schema"
        cd -
        exit 1;
    fi
    cd -
}

# command for database setup
case $1 in
    create)
        create $2
    ;;
    drop)
        drop $2
    ;;
    generate)
        generate $2 $3 $4
    ;;
    migrate)
        migrate $2 $3 $4
    ;;
esac 