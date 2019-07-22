#!/bin/bash
go build ./src/rest_api/cmd/app
sudo systemctl start restapi
sudo systemctl enable restapi