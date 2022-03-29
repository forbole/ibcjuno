## Configuration
The default `config.yaml` file should look like the following:

<details>

<summary>Default config.yaml file</summary>

```yaml
database:
  host: localhost
  max_idle_connections: 0
  max_open_connections: 0
  name: database-name
  password: password
  port: 5432
  schema: public
  ssl_mode: 
  user: user

tokens:
    token:
        - name: Desmos
          units:
            - denom: udsm
              exponent: 0
              ibc_denom: ibc/EA4C0A9F72E2CEDF10D0E7A9A6A22954DB3444910DB5BE980DF59B05A46DAD1C
            - denom: dsm
              exponent: 6
              price_id: desmos
        - name: Atom
          units:
            - denom: uatom
              exponent: 0
              ibc_denom: ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2
            - denom: atom
              exponent: 6
              price_id: cosmos
        - name: Osmosis
          units:
            - denom: uosmo
              exponent: 0
              ibc_denom: uosmo
            - denom: osmo
              exponent: 6
              price_id: osmosis
```

</details>

Let's see what each section refers to:

- [`database`](#database)
- [`tokens`](#tokens)

## `database`
This section contains all different configuration related to the PostgreSQL database where IBCJuno will write the data.

| Attribute | Type | Description | Example |
| :-------: | :---: | :--------- | :------ |
| `host` | `string` | Host where the database is found | `localhost` | 
| `port` | `integer` | Port to be used to connect to the PostgreSQL instance | `5432` |
| `name` | `string` | Name of the database to which connect to | `ibcjuno` | 
| `user` | `string` | Name of the user to use when connecting to the database. This user must have read/write access to all the database. | `ibcjuno` | 
| `password` | `string` | Password to be used to connect to the database instance | `password` | 
| `schema` | `string` | Schema to be used inside the database (default: `public`) | `public` | 
| `ssl_mode` | `string` | [PostgreSQL SSL mode](https://www.postgresql.org/docs/9.1/libpq-ssl.html) to be used when connecting to the database. If not set, `disable` will be used. | `verify-ca` |
| `max_idle_connections` | `integer` | Max number of idle connections that should be kept open (default: `1`) | `10` |
| `max_open_connections` | `integer` | Max number of open connections at any time (default: `1`) | `15` | 
