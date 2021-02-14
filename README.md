# investor

Helper tool for managing your stock assets on [Tinkoff Investment Platform](https://www.tinkoff.ru/invest/).

I wanted to track more data about my portfolio than there is in the app so I decided to make my own tracker. The bank has [an API](https://tinkoffcreditsystems.github.io/invest-openapi/) which I am using here. Feel free to use my [referral link](https://www.tinkoff.ru/sl/3tqgECf6gYa) if you also want to invest some of your fortune =)

**This is an early stage for the project. Nothing is final.**

## Usage

1. Clone the repo with `git clone https://github.com/fabritsius/investor`
2. Go to the project root with `cd investor/`
3. Run `make db` if you have [docker](https://www.docker.com/) or make sure you have [cassandra](https://cassandra.apache.org/) running locally
4. Run `make keyspace` to create project keyspace (requires [cassandra](https://cassandra.apache.org/) running in [docker](https://www.docker.com/))
5. Create a new [.env](.env) file using [.env.example](.env.example) as an example
6. Run `make tinkoff` to start tinkoff plugin
7. Run `make aggregate` from a separate console to test the plugin

## TODO

- [x] Get basic data about stock portfolio from Tinkoff
- [x] Use gRPC to connect plugins to the aggregator
- [ ] Store portfolio history in a database
- [ ] Use a [Telegram Bot](https://core.telegram.org/bots) as a user interface
- [ ] Create a webpage with the project description
- [ ] Add an advanced prediction features for market analysis
