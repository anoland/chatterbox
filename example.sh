#!/bin/bash
#go get github.com/secondbit/wendy || exit ${LINENO}
#[ -f peer-rand.go ] || curl -ksO https://gist.githubusercontent.com/nicerobot/9753246/raw/peer-rand.go
rm 80*.out
#go build peer-rand.go && 

  ./peer-rand 8080 >8080.out 2>&1 & pids=(${pids[*]} $!)
	echo "sleeping main node"	
  for i in 1 2 3 4; do
	sleep 3
    ./peer-rand 808${i} 8080 2>&1 & pids=(${pids[*]} $!)
  done

trap "kill ${pids[*]}; exit 1" 0 1 2 3
tail -f 8080.out & pids=(${pids[*]} $!)
sleep ${1:-300}
for p in ${pids[*]}; do sleep ${2:-5}; echo "kill ${p}"; kill ${p}; done

