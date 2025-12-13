#!/bin/bash

rm -f database

declare -A position_by_key # hash map
declare -i last_index=0 # integer

db_set() {
    echo "$1,$2" >> database
    last_index=$(( last_index + 1 ))
    position_by_key[$1]=$last_index
}

db_get() {
    position=${position_by_key[$1]}
    awk 'NR=='$position database | sed -e "s/^$1,//" | tail -n 1
}

db_set key value
db_set key2 value2
db_set key valuebis
db_set key3 value3
db_set key3 value3bis
db_set key4 value4
db_set key4 value4bis