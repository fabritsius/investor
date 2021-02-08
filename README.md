# investor

Helper tool for managing your stock assets on [Tinkoff Investment Platform](https://www.tinkoff.ru/invest/).

I wanted to track more data about my portfolio than there is in the app so I decided to make my own tracker. The bank has [an API](https://tinkoffcreditsystems.github.io/invest-openapi/) which I am using here. Feel free to use my [referral link](https://www.tinkoff.ru/sl/3tqgECf6gYa) if you also want to invest some of your fortune =)

**This is an early stage for the project. Nothing is final.**

## Usage

1. Clone the repo with `git clone https://github.com/fabritsius/investor`
2. Go to the project root with `cd investor/`
3. Ensure you have [Cassandra](https://cassandra.apache.org/) running on standard local ports (or run `make db`)
3. Init the DB with `investor.accounts_by_user` table (see the code, better docs will come later)
4. insert a user where `account` is `tinkoff` and `key` is `YOUR_TINKOFF_TOKEN`
5. Run `make tinkoff` to start tinkoff plugin
6. Run `make aggregate` from a separate console to test the plugin

## TODO

- [x] Get basic data about stock portfolio from Tinkoff
- [x] Use gRPC to connect plugins to the aggregator
- [ ] Store portfolio history in a database
- [ ] Use a [Telegram Bot](https://core.telegram.org/bots) as a user interface
- [ ] Create a webpage with the project description
- [ ] Add an advanced prediction features for market analysis
