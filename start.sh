#!/bin/sh
rm -fr /tmp/rinha*
export $(cat .env | xargs) && ./rinha
