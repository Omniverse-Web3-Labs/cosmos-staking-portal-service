# Cosmos staking portal service

We have developed a Staking product for blockchain projects built on CometBFT, and this repository contains the backend service code.

## Prerequisites
Go >= 1.23

## Install

Clone the repository

```
git clone https://github.com/Omniverse-Web3-Labs/cosmos-staking-portal-service.git
```

Install dependencies
```
go mod tidy
```

## Configuration

- databases
    - postgres
        - default: used to store user related data
        - subquery: used to store data indexed by subquery
- endpoint: the endpoint of CometBFT node
- redis: the redis service information

## Run

```
go run cmd/api/main.go
```