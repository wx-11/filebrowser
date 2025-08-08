#\!/bin/bash
echo 'Hello from demo script\!'
date
echo 'System info:'
uname -a
echo 'Free memory:'
cat /proc/meminfo | head -3
echo 'Demo completed\!'
