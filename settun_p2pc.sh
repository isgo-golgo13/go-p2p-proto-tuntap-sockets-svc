#! /bin/sh

TUN_DEV_DEFAULT=tun0
TUN_DEV=$1

if [ $# -eq 0 ]; then
    sudo ip tuntap add ${TUN_DEV_DEFAULT} mode tun
    sudo ip link set dev ${TUN_DEV_DEFAULT} up
    sudo ip addr add 10.1.0.10/24 dev ${TUN_DEV_DEFAULT}
else
    sudo ip tuntap add ${TUN_DEV} mode tun
    sudo ip link set dev ${TUN_DEV} up
    sudo ip addr add 10.1.0.10/24 dev ${TUN_DEV}
fi

