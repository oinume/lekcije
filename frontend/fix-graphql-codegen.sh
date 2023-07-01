#!/bin/sh

perl -pi -e 's!^import { RequestInit }!//import { RequestInit }!g' src/graphql/generated.ts
