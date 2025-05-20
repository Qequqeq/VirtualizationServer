#!/bin/sh
qemu_count=$(ps | grep "qemu-system-x86_64" | grep -v "grep" | wc -l)
if [ "$qemu_count" -eq 0 ]; then
	echo 1 > /root/database/last_id
fi
