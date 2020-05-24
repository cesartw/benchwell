#!/bin/bash

qemu-system-x86_64 -m 16G -boot order=d -drive file=testing-vm.qemu,format=raw -cpu host -smp cores=3 -enable-kvm
