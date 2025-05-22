#!/bin/bash

curl -L -b /tmp/cookies.txt \
"https://drive.usercontent.google.com/download?id=1Mjlt5c9IfywGFYvySt7aGffxG3s33rQF&export=download&confirm=t&uuid=6f0cd54f-f138-4cf6-93d3-1ed2c20aa1fc" \
-o alpine_base.qcow2
mv alpine_base.qcow2 image/
