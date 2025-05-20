#!/bin/bash

BASE_IMAGE="/root/VirtualizationServer/image/alpine_base.qcow2"
VM_DIR="/root/VirtualizationServer/image/vms"
SSH_PORT_START=10001
ID_FILE="/root/VirtualizationServer/database/last_id"
USER_FILE="/root/VirtualizationServer/database/users"

is_port_avaliable(){
	local port=$1
	! ss -tuln | grep -q ":${port} "
}

get_next_id(){
    [ ! -f "$ID_FILE" ] && echo "0" > "$ID_FILE"
    exec 200>>"$ID_FILE"
    flock -x 200
    current_id=$(<"$ID_FILE")
    next_id=$((current_id + 1))
    
    echo "$next_id" > "$ID_FILE"
    exec 200>&-
    echo "$current_id"
}


create_vm(){
    USERNAME="$1"
    USER_IMAGE="${VM_DIR}/${USERNAME}.qcow2"

    exec 300>>"$USER_FILE"
    flock -x 300

    if ! grep -q "^${USERNAME}$" "$USER_FILE"; then 
	 [ ! -f "$USER_FILE" ] && touch "$USER_FILE"
	 qemu-img create -f qcow2 -b "$BASE_IMAGE" -F qcow2 "$USER_IMAGE" 5G >/dev/null 2>&1 || {
         echo "Failed to create disk image" >&2
         return 1
    }
    echo "$USERNAME" >> "$USER_FILE"
    fi
    exec 300>&-

    USER_ID=$(get_next_id)
    SSH_PORT=$((SSH_PORT_START + USER_ID))
    if ! is_port_avaliable "$SSH_PORT"; then 
	    while ! is_port_avaliable "$SSH_PORT"; do
		    SSH_PORT=$((SSH_PORT + 1))
	    done
    fi

    qemu-system-x86_64 \
        -daemonize \
        -drive "file=${USER_IMAGE},format=qcow2,if=virtio" \
        -nic "user,hostfwd=tcp::${SSH_PORT}-:22" \
        -m 512 \
	> /dev/null 2>&1 || {
        echo "Failed to start VM" >&2
        return 1
	}
    echo $SSH_PORT
    /root/vm-api/scripts/start_port.sh start $SSH_PORT
}

case "$1" in
    create)
        create_vm "$2"
        ;;
    *)
        echo "Using: $0 create <username>"
        exit 1
esac
