#!/bin/sh
export $(cat .env | xargs) && ./rinha
