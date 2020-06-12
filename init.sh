#!/bin/sh
if [ -d $GOPATH/github.com/ShkrutDenis/go-migrations/template/ ]
then
  cp -a $GOPATH/github.com/ShkrutDenis/go-migrations/template/. ./
  exit 1
fi
if [ -d vendor/github.com/ShkrutDenis/go-migrations/template ]
then
  cp -a vendor/github.com/ShkrutDenis/go-migrations/template/. ./
  exit 1
fi
echo "Dependency path not found"
exit 0