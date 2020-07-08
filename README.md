# DP Chain Hub ::sparkle::

[![version](https://img.shields.io/github/tag/tokenchain/ixo-blockchain.svg)](https://github.com/tokenchain/ixo-blockchain/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tokenchain/ixo-blockchain)](https://goreportcard.com/report/github.com/tokenchain/ixo-blockchain)
[![LoC](https://tokei.rs/b1/github/tokenchain/ixo-blockchain)](https://github.com/tokenchain/ixo-blockchain)

This is the official repository for the Internet of Impact Relayer Hub (DP Hub)

> This document will have all details necessary to help getting started with DP Hub

## Documentation
- Guide for setting up a Relayer on the Darkpool Test Network: [here](https://github.com/tokenchain/docs/blob/master/developer-tools/test-networks/join-a-test-network.md)
- Modules specification: look into `x/<module>/spec`

### Other documentations
- [Whitepaper](https://github.com/tokenchain/dp-hub/blob/master/doc/whitepaper.md)
- [Chinese Developer Guide](https://github.com/tokenchain/dp-hub/blob/master/doc/commands.md)


### Scripts
Quick-start:


```bash
cd ./scripts/
sh installnewchain.sh
sh clean_build.sh && bash run_with_some_data.sh  # Option 1
sh clean_build.sh && bash run_with_all_data.sh   # Option 2
```

To run without resetting data:
```bash
cd ./scripts/
sh run_only.sh
```

(Optional) Once the chain has started, run one of the following:

Add more data and activity:
```bash
cd ./scripts/
sh add_dummy_testnet_data.sh
```

Demos:
```bash
cd ./scripts
sh demo_bonds.sh              # Option 1
sh demo_bonds_swapper.sh      # Option 2
sh demo_project.sh            # Option 3
sh demo_tx_broadcast_rest.sh  # Option 4
sh demo_tx_broadcast_rpc.sh   # Option 5
```

Nginx setup
To expose ports on nginx server
```shell script

server {
        listen 1317;
        listen [::]:1317;
        server_name cli.darkpool.vip;
        add_header 'Access-Control-Allow-Origin' '*';
        add_header 'Access-Control-Allow_Credentials' 'true';
        add_header 'Access-Control-Allow-Headers' 'Authorization,Accept,Origin,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range';
        add_header 'Access-Control-Allow-Methods' 'GET,POST,OPTIONS,PUT,DELETE,PATCH';

        location / {
                if ($request_method = 'OPTIONS') {
                        add_header 'Access-Control-Max-Age' 1728000;
                        add_header 'Content-Type' 'text/plain charset=UTF-8';
                        add_header 'Content-Length' 0;
                        return 204;
                }
                proxy_redirect off;
                proxy_set_header host $host;
                proxy_set_header X-real-ip $remote_addr;
                proxy_set_header X-forward-for $proxy_add_x_forwarded_for;
                proxy_pass http://localhost:1317;
        }
}

server {
        listen 26657;
        listen [::]:26657;
        server_name demon.darkpool.vip;
        add_header 'Access-Control-Allow-Origin' '*';
        add_header 'Access-Control-Allow_Credentials' 'true';
        add_header 'Access-Control-Allow-Headers' 'Authorization,Accept,Origin,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,Cache-Control,Content-Type';
        add_header 'Access-Control-Allow-Methods' 'GET,POST,OPTIONS,PUT,DELETE,PATCH';

        location / {
                if ($request_method = 'OPTIONS') {
                        add_header 'Access-Control-Max-Age' 1728000;
                        add_header 'Content-Type' 'text/plain charset=UTF-8';
                        add_header 'Content-Length' 0;
                        return 204;
                }
                proxy_redirect off;
                proxy_set_header host $host;
                proxy_set_header X-real-ip $remote_addr;
                proxy_set_header X-forward-for $proxy_add_x_forwarded_for;
                proxy_pass http://localhost:26657;
        }
}
```


#### API Doc
Please find the api document to be located at `*:1317/swagger-ui/` at the LCD.
