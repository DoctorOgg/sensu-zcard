#!/bin/bash

function redis-cli() {
  /opt/homebrew/bin/redis-cli -h localhost $@
}

redis-cli zadd fubar:farts:222 1 "fart1"
redis-cli zadd fubar:farts:222 2 "fart2"
redis-cli zadd fubar:farts:222 3 "fart3"
redis-cli zadd fubar:farts:222 4 "fart4"
redis-cli zadd fubar:farts:222 5 "fart5"
redis-cli zadd fubar:farts:222 6 "fart6"

redis-cli zadd fubar:farts:43 1 "fart1"
redis-cli zadd fubar:farts:43 2 "fart2"
redis-cli zadd fubar:farts:43 3 "fart3"



