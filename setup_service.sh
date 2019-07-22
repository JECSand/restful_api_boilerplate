#!/bin/bash
echo "[Unit]
Description=rest api initialization

[Service]
PIDFile=/tmp/restapi.pid-4040
User=api_user
Group=api_user
WorkingDirectory="$1"
ExecStart=/bin/bash -c '"$1"/app'

[Install]
WantedBy=multi-user.target" >> /usr/lib/systemd/system/restapi.service