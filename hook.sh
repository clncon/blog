#!/bin/bash
echo "----------start pull update!!!----------------"
git pull
echo "----------install update!!!-------------------"
go install
echo "----------restart blog!!!---------------------"
ps -ef | grep -v grep | grep blog  | awk '{print $2}' | xargs kill -9
/root/go/bin/blog.sh
echo "----------auto deploy complate!!!!---------------------"