#!/bin/bash
#
#
#
set -euo pipefail

log() {
	echo "[$(basename "$0") $(date +%Y%m%d%H%I%M) ] $*" >> bootstrap.log
}

cd /home/core
log "Bootstrap starting as user $(id -u):$(id -g).."
log "Trying to start etcd2, all systemd units: $(systemctl)"
systemctl start etcd2
log "Started etcd2"

