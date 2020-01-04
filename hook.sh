#!/bin/bash
echo "----------start pull update!!!----------------"
git pull
echo "----------install update!!!-------------------"
go install
echo "----------restart blog!!!---------------------"
ps -ef | grep -v grep | grep /root/go/bin/blog  | awk '{print $2}' | xargs kill -9
nohup /root/go/bin/blog > /root/go/bin/run.log 2>&1 &
echo "----------auto deploy complate!!!!---------------------"
