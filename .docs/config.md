## Configuration
The default `config.yaml` file should look like the following:

<details>

<summary>Default config.yaml file</summary>

```yaml
database:
    name: database-name
    host: localhost
    port: 5432
    user: user
    password: password
    schema: public
    max_open_connections: 1
    max_idle_connections: 1
tokens:
    token:
        - name: Desmos
          units:
            - denom: dsm
              ibc_denom:
                - denom: udsm
                  src_chain: desmos
                  dst_chain: osmosis
                  channel: channel-1
                  ibc_denom: ibc/EA4C0A9F72E2CEDF10D0E7A9A6A22954DB3444910DB5BE980DF59B05A46DAD1C
              exponent: 6
              price_id: desmos
```

</details>

Let's see what each section refers to:

- [`database`](#database)
- [`tokens`](#tokens)

## `database`
This section contains all different configurations related to the PostgreSQL database where IBCJuno will write the data.

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

## `tokens`
This section contains the details of the tokens that IBCJuno will fetch at the latest prices.

| Attribute | Type | Description | Example |
| :-------: | :---: | :--------- | :------ |
| `token` | `object` | Contains token configuration data | | 

### token
This section contains the info about token names & units 

| Attribute | Type | Description | Example |
| :-------: | :---: | :--------- | :------ |
| `name` | `string` | Default denom name | `Desmos` | 
| `units` | `object` | Contains the details about token units | | 

### units
This section contains the details about token units 
| Attribute | Type | Description | Example |
| :-------: | :---: | :--------- | :------ |
| `denom` | `string` | Denom unit | `dsm` |
| `exponent` | `integer` | Denom unit exponent value | `6` |
| `ibc_denom` | `string` | IBC denom unit | `ibc/EA4C0A9F72E2CEDF10D0E7A9A6A22954DB3444910DB5BE980DF59B05A46DAD1C` |
| `price_id` | `string` | Coingecko ID | `desmos` |

