#!/bin/sh
if [ -d $GOPATH/src/go-migrations/template/ ]
then
  cp -a $GOPATH/src/go-migrations/template/. ./
  exit 1
fi
if [ -d vendor/go-migrations/template ]
then
  cp -a vendor/src/go-migrations/template/. ./
  exit 1
fi
echo "Dependency path not found"
exit 0