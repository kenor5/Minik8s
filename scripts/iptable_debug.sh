#!/bin/bash

allChain=$(iptables -L -t nat    | grep "Chain" | awk '{print $2}' | tac )

for i in $allChain
do
    iptables -t nat -I $i 1 -j LOG --log-prefix "Debug: $i"
done