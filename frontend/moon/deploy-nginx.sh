#!/bin/bash

# For my personal development environment only!!! Do not run this 

remote="alphaboi@curtisnewbie.com"
remote_path="/home/alphaboi/services/nginx/html/bolobao/"

NODE_OPTIONS=--openssl-legacy-provider ng build --prod;
scp -r ./dist/moon/* "${remote}:${remote_path}"




