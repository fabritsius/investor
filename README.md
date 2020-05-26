# investor

Helper tool for managing your stock assets on [Tinkoff Investment Platform](https://www.tinkoff.ru/invest/).

I wanted to track more data about my portfolio than there is in the app so I decided to make my own tracker. The bank has [an API](https://tinkoffcreditsystems.github.io/invest-openapi/) which I am using here. Feel free to use my [referral link](https://www.tinkoff.ru/sl/3tqgECf6gYa) if you also want to invest some of your fortune =)

**This is an early stage for the project. Nothing is final.**

## Usage

1. Clone the repo with `git clone https://github.com/fabritsius/investor`
2. Go to the project root with `cd investor`
3. Fill in your API token into `'env.sh'` file
4. Set environment variables with `. env.sh`
5. Run the code with `go run main.go`

## TODO

- [x] Get basic data about stock portfolio from Tinkoff
- [ ] Store portfolio history in a database
- [ ] Use a [Telegram Bot](https://core.telegram.org/bots) as a user interface
- [ ] Create a webpage with the project description
- [ ] Add an advanced prediction features for market analysis
