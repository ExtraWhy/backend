#!/bin/bash



curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 1, "money" : 10000, "name" : "Lubaka"}'

curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 2, "money" : 10000, "name" : "Kalniq"}'

curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 3, "money" : 10000, "name" : "Gandalf"}'

curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 4, "money" : 10000, "name" : "Krasena"}'

curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 5, "money" : 10000, "name" : "Ekstramena"}'

curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 6, "money" : 10000, "name" : "Shto?"}'

curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 7, "money" : 10000, "name" : "Kucheto"}'

curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 8, "money" : 10000, "name" : "Bonbonev"}'

curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 9, "money" : 999999999, "name" : "Skalata"}'

curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 10, "money" : 10000, "name" : "Robota"}'

