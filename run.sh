#!/bin/bash
export $(cat .env | xargs) && ./tmp/rinha
