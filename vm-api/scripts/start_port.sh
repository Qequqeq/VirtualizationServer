#!/bin/sh
start_port(){
	port=$1
	tuna tcp $port --log /root/VirtualizationServer/database/tuna_ports >/dev/null 2>&1 &
}


case "$1" in 
	start)
		start_port "$2"
		;;
	*)
		echo "Using: $0 start <port>"
		exit 1
esac
