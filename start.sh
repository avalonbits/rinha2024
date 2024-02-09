#!/bin/sh
export $(cat .env | xargs) && /app/rinha
