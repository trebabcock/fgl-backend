#!/bin/bash

CGO_ENABLED=0 go build

sshpass -p "fgldbpass" ssh root@134.209.25.57 <<'EOF'
kill -9 $(cat save_pid.txt)
echo $(cat blah.txt)
rm save_pid.txt fgl-backend my.log
exit
EOF

sshpass -p "fgldbpass" sftp root@134.209.25.57 <<EOF
put /mnt/c/Users/Tre/Documents/Projects/Go/fgl-database/fgl-backend/fgl-backend
exit
EOF

sshpass -p "fgldbpass" ssh root@134.209.25.57 <<'EOF'
nohup ./fgl-backend HwXMQawGhx > my.log &
echo $! > save_pid.txt
EOF
