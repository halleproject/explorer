<h2 align="center">    
  Halle Chain Explorer     
</h2>

 ## Overview 
 This project is a BlockChain Explorer of Halle Chain. The program supports talented developers and researchers in creating free and open-source software that would enable new innovations and businesses in the crypto community.
    
## Prerequisite

- Requires [Go 1.14+](https://golang.org/dl/)

- Run a node of Halle Chain

- Prepare PostgreSQL Database

## Folder Structure

    /
    |- chain-exporter
    |- mintscan
    |- mintscan-binance-dex-frontend

#### Chain Exporter

chain-exporter watches a full node of Halle Chain and export chain data into PostgreSQL database.

#### Mintscan

mintscan provides any necesarry custom APIs.

#### mintscan-binance-dex-frontend

mintscan-binance-dex-frontend show blocks and transactions.

## Configuration

For configuration, it uses human readable data-serialization configuration file format called [YAML](https://en.wikipedia.org/wiki/YAML).

To configure `chain-exporter` | `mintscan`, you need to configure  `config.yaml` file in each folder. Reference `example.yaml`.

To configure `mintscan-binance-dex-frontend`, you need access the config.js file directly under the src/ directory and specify your backend dev and prod apis.

**_Note that the configuration needs to be passed in via `config.yaml` file, so make sure to change the name to `config.yaml`._**

## Install

#### Git clone this repo
```shell
git clone https://github.com/halleproject/explorer.git    
```

#### Build by Makefile
```shell
cd explorer/chain-exporter
make build

cd explorer/mintscan
make build

cd explorer/mintscan-binance-dex-frontend
#comment out the following line in src/Root.js
#import "./firebase"
yarn dev  
yarn build:dev  
```    

## Database 

This project uses [Golang ORM with focus on PostgreSQL features and performance](https://github.com/go-pg/pg). Once `chain-exporter` begins to run, it creates the following database tables if not exist already.

- Block
- PreCommit
- Transaction
- Validator

## Contributing    
 We encourage and support an active, healthy community of contributors — any contribution, improvements, and suggestions are always welcome!     
    
### Note before I'm bothered to actually write the guide 
```    
I'm very conscious of how much more work could be done to make this project    
- a very general term - but just better.    
A lot of the code (with great reluctance) are even in my possibly abysmal standards    
'not up to par', I still have nightmares of the fact that I didn't adhere    
to the 'rule of hooks' in many occasions.    
(I beg of you to put down your pitchforks after reading the myriads of warning messages    
spewing out of this    
```    
   
###### _Ironically that single monstrosity of a file took up about 30% of the total time that I worked on this project_ 
```    
Please feel free to help clear up this displeasant stream of code    
before it becomes the hopeless mess it is most certainly destined to become    
without your most awaited upon help.    
    
yours sincerely, with a grain of salt *wink*    
```    
 ## License    
 Released under the Apache 2.0 License