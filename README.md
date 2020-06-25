# DP Chain Hub

[![version](https://img.shields.io/github/tag/tokenchain/ixo-blockchain.svg)](https://github.com/tokenchain/ixo-blockchain/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tokenchain/ixo-blockchain)](https://goreportcard.com/report/github.com/tokenchain/ixo-blockchain)
[![LoC](https://tokei.rs/b1/github/tokenchain/ixo-blockchain)](https://github.com/tokenchain/ixo-blockchain)

This is the official repository for the Internet of Impact Relayer Hub (DP Hub)

> This document will have all details necessary to help getting started with DP Hub

## Documentation
- Guide for setting up a Relayer on the Darkpool Test Network: [here](https://github.com/tokenchain/docs/blob/master/developer-tools/test-networks/join-a-test-network.md)
- Modules specification: look into `x/<module>/spec`

## Scripts
Quick-start:
```bash
cd ixo-blockchain/scripts/
bash clean_build.sh && bash run_with_some_data.sh  # Option 1
bash clean_build.sh && bash run_with_all_data.sh   # Option 2
```

To run without resetting data:
```bash
cd ixo-blockchain/scripts/
bash run_only.sh
```

(Optional) Once the chain has started, run one of the following:

- Add more data and activity:
```bash
cd ixo-blockchain/scripts/
bash add_dummy_testnet_data.sh
```

- Demos:
```bash
cd ixo-blockchain/scripts
bash demo_bonds.sh              # Option 1
bash demo_bonds_swapper.sh      # Option 2
bash demo_project.sh            # Option 3
bash demo_tx_broadcast_rest.sh  # Option 4
bash demo_tx_broadcast_rpc.sh   # Option 5
```

- [Whitepaper](https://www.find.com)

