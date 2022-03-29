# IBCJuno Setup
Setting up IBCJuno is pretty straightforward. It requires three things to be done:
1. Install IBCJuno.
1. Initialize the configuration.
2. Start the parser.

## Installing Juno
In order to install IBCJuno you are required to have [Go 1.17+](https://golang.org/dl/) installed on your machine. Once you have it, the first thing to do is to clone the GitHub repository. To do this you can run

```shell
$ git clone https://github.com/MonikaCat/ibcjuno.git
```

Then, you need to install the binary. To do this, run

```shell
$ make install
```

This will put the `ibcjuno` binary inside your `$GOPATH/bin` folder. You should now be able to run `ibcjuno` to make sure it's installed:

```shell
$ ibcjuno
IBCJuno is a IBC price aggregator and exporter. It queries the latest IBC tokens prices
and stores it inside PostgreSQL database. IBCJuno is meant to run with a GraphQL layer on top
to ease the ability for developers and downstream clients to query the latest price of any IBC token.

Usage:
  IBCJuno [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize configuration files
  start       Start IBCJuno price aggregator
  version     Print the version of IBCJuno

Flags:
  -h, --help          help for IBCJuno
      --home string   Set the home folder of the application, where all files will be stored (default "/Users/monikapusz/.IBCJuno")

Use "IBCJuno [command] --help" for more information about a command.
```
