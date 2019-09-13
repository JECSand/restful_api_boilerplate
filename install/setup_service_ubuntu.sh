#!/bin/bash
echo "[Unit]
Description=rest api initialization

[Service]
PIDFile=/tmp/restapi.pid-4040
User="$2"
Group="$2"
WorkingDirectory="$1"
ExecStart=/bin/bash -c '"$1"/app'

[Install]
WantedBy=multi-user.target" >> /lib/systemd/system/restapi.service