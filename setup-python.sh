#!/bin/sh

if [ ! -d .venv ]; then
    virtualenv --distribute --python /usr/local/bin/python2 .venv
    . .venv/bin/activate
    pip install -r requirements.txt
fi
