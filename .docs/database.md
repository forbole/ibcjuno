# Database
IBCJuno is designed to run with [PostgreSQL](https://www.postgresql.org/) database to store the latest IBC tokens prices. You should install the latest version on your machine following the [official documentation](https://www.postgresql.org/download/)

Once installed you need to create a new database and a new user that is going to read and write data inside it.  
Then, once that's done, you need to run the SQL query that you can find inside the [`db/schema` folder](../db/schema).  

Once that's done, you are ready to [continue the setup](setup.md).