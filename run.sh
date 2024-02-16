#!/bin/bash
export $(cat .env | xargs) && ./tmp/rinha -port 9999
