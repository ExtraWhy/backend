#!/bin/bash
#poor man test bot
for (( ;; )) ; do 
	for ((i = 1 ; i <= 10 ; i++ )); do 
	curl -X GET http://localhost:8081/players/$i/bet/10
	done
done
