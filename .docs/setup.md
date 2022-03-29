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

## Initializing the configuration
In order to correctly initialize IBCJuno you need to run the following command: 

```shell
$ ibcjuno init
```

This will create `.ibcjuno` directory with `config.yaml` file inside.  
**Note:** If you want to change the default directory used by IBCJuno you can do this using the `--home` flag:

```shell
$ ibcjuno init --home /path/to/my/folder
```

Once the file is created, you are required to edit it and update database and tokens section. To do this you can run

```shell
$ nano ~/.ibcjuno/config.yaml
```

For a better understanding of what each section and field refers to, please read the [config reference](config.md).

## Running IBCJuno
Once the configuration file has been setup, you can run IBCJuno using the following command:

```shell
$ ibcjuno start
```

If you are using a custom directory for the configuration file, please specify it using the `--home` flag:


```shell
$ ibcjuno start --home /path/to/my/config/folder
```

We highly suggest you running IBCJuno as a system service so that it can be restarted automatically in the case it stops. To do this you can run:

```shell
$ sudo tee /etc/systemd/system/ibcjuno.service > /dev/null <<EOF
[Unit]
Description=IBCJuno aggregator
After=network-online.target

[Service]
User=$USER
ExecStart=$GOPATH/bin/ibcjuno start
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF
```

Then you need to enable and start the service:

```shell
$ sudo systemctl enable ibcjuno
$ sudo systemctl start ibcjuno
```

Then you can check the status of ibc service:

```shell
$ sudo systemctl status ibcjuno
```
