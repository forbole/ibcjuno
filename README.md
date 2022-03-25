# IBCJuno

IBCJuno is IBC price aggregator and for Cosmos [IBC module](https://github.com/cosmos/ibc).
IBCJuno fetches the latest price of IBC tokens and stores it inside a [PostgreSQL](https://www.postgresql.org/) database on top of which [GraphQL](https://graphql.org/) APIs can be created using [Hasura](https://hasura.io/).

## Install IBCJuno 
Run inside ibcjuno directory: 

```shell
$ make install
```

## GraphQL integration
If you want to know how to run a GraphQL server that allows to expose the parsed data, please refer to the following guides: 

- [PostgreSQL setup with GraphQL](.docs/postgres-graphql-setup.md)

## Contributing
Contributions are welcome! Please open an Issues or Pull Request for any changes.

## License
[CCC0 1.0 Universal](https://creativecommons.org/share-your-work/public-domain/cc0/)