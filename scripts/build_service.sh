#!/usr/bin/env bash

DAEMONNAME="ixod"
CLINAME="ixocli"
CHAIN_ID="darkpool-1x"
DAEMON=$GOBIN/$DAEMONNAME
CLI=$GOBIN/$CLINAME


su $USERNAME <<EOSU
$DAEMON init "Darkpool node"
EOSU

sleep 5


echo "---"
echo "Your peer ID:"
$DAEMON tendermint show-node-id
echo "---"

cat << EOF > /etc/systemd/system/$DAEMONNAME.service
# /etc/systemd/system/$DAEMONNAME.service
[Unit]
Description=$DAEMONNAME Node
After=network.target
 
[Service]
Type=simple
User=$USERNAME
WorkingDirectory=$HOME
ExecStart=$DAEMON start
Restart=on-failure
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
sleep 3
systemctl enable $DAEMONNAME.service
echo "Service created at /etc/systemd/system/$DAEMONNAME.service."
echo "Run 'systemctl start $DAEMONNAME.service' to start the node"
