#!/bin/bash

IN="c2f_admin_webapp"
log="/var/log"
pid="/var/run"
export GOPATH=$HOME/work/c2f_admin_webapp/
export PATH=$PATH:$HOME/work/c2f_admin_webapp/bin
arr=$(echo $IN | tr ";" "\n")
dt=$(date '+%d-%m-%YT%H:%M:%S')
red='\033[0;31m'
green='\033[0;32m'
NC='\033[0m'
if [ $# == 1 ]
then
        for x in $arr
        do
                if [ $1 == "start" ]
                then
                        no_subsys=0
                        if [ -f $pid/$x.pid ]
                        then
                                subsys_pid="$(<$pid/$x.pid)"
                                no_subsys=`ps -ewwo pid | grep -wc "$subsys_pid"`
                        fi
                        if [ $no_subsys -gt 0 ]
                        then
                                echo -e "${red}$x ${green}is Already Running, first stop it${NC}"
                        else
                                echo -e "${red}Starting Process ${NC}--> ${green}$x${NC}"
                                $x >>$log/$x.log 2>&1 &
                                echo $! >$pid/$x.pid
                        fi
                elif [ $1 == "stop" ]
                then
                        if [ ! -d "$GOPATH/tmp/" ]; then
                                mkdir $GOPATH/tmp/
                        fi
                        if [ -f "$log/$x.log" ]; then
                                mv $log/$x.log $GOPATH/tmp/$x$dt.log
                        fi
                        subsys_pid="$(<$pid/$x.pid)"
                        echo -e "${red}Stopping Process ${NC}--> ${green}$x:$subsys_pid${NC}"
                        kill $subsys_pid
                        rm -f $pid/$x.pid
                elif [ $1 == "status" ]
                then
                        if [ -f "$pid/$x.pid" ]; then
                                subsys_pid="$(<$pid/$x.pid)"
                                echo -e "${red}$x${NC}:${green}$subsys_pid${NC}"
                        else
                                echo -e "${red}$x${NC}:${green}No process running on this name :)${NC}"
                        fi
                elif [ $1 == "help" ]
                then
                        echo -e "${red}start             ${green} Option to Start System.${NC}"
                        echo -e "${red}stop              ${green}Option to Stop System.${NC}"
                        echo -e "${red}status            ${green}Option to check status System.${NC}"
                else
                        echo -e "${red}wrong argument. Use help option for more clarification.${NC}"
                        break
                fi
                sleep 0.2
done
else
        echo "Wrong number of arguments"
fi
