#!/bin/bash
MAILTO="andrea@sqlitecloud.io,marco@sqlitecloud.io"
MAILFROM="unit-status-mailer@sqlitecloud.io"
UNIT=$1
SUBJECT=$2

EXTRA=""
for e in "${@:3}"; do
  EXTRA+="$e"$'\n'
done

UNITSTATUS=$(systemctl status $UNIT)

ssmtp -t <<EOF
From:$MAILFROM
To:$MAILTO
Subject:$SUBJECT: $UNIT

Status report for unit: $UNIT
$EXTRA

$UNITSTATUS
EOF

echo -e "Status mail sent to: $MAILTO for unit: $UNIT"