#### HOWTO REDIS<>WS


#### Install , check and subscribe for redis
1. `sudo apt install redis`
2. `sudo systemctl start redis`
3. `redis-cli ping`
4. `redis-cli PUBLISH jackpot:updates '{"jackpot": 9999.99}'`


