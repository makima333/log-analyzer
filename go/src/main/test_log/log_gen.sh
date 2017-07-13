#!/bin/bash

for i in `seq 1 $1`
do
  echo 10.2.3.`jot -r 1 1 30` - - [`jot -r -w %02d 1 1 30`/$(cat ~/go/src/main/test_log/month.txt |shuf -n 1)/`jot -r 1 2016 2017`:`jot -r -w %02d 1 00 23`:`jot -r -w %02d 1 00 59`:47 +0900] "\"GET / HTTP/1.1\" 200 854 \"-\" \"Mozilla/4.0 (compatible; MSIE 5.5; Windows 98)\""
done
