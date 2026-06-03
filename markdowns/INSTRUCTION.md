# Complex guide to make it work

1. Make sure you have latest distribution of golang (at least 1.21.)
2. Make sure you have local distribution of docker
3. Configure .env files: 
   - `./.env`
   - `./db/.env`
   
   For example: copy the content of .env.example files
4. go to the `db` directory: 
    `cd db`
5. Download images and run containers:
    `docker compose up -d`
6. Wait a litlle while and log in to the pgadmin container service: `localhost:5050`:
   Login is your `PGADMIN_EMAIL` and password is your `PGADMIN_PASSWORD` 
7. Estabilish connection to the new server with proper information from your `.env`:
   1. Click `Add New Server`
   2. Enter the name you would like to call your connection e.g `kitty`
   3. Switch to the `Connection` tab
   4. `Host name/address` : Important info: if you run containers as a localhost please enter the `db`value in this field
   5. `Port` : enter your `POSTGRES_PORT` value e.g. `5432`
   6. `Maintenance database`: enter your `POSTGRES_DB` value
   7. `Username`: enter your `POSTGRES_USER` value
   8. `Password`: enter your `kittyhub` value
   9. Click `save` button
8. Go your connection, database and schemas.
9.  Open new query tool for that database
10. Create new table using script from [init.sql](../db/init.sql)
11. Execute query
12. Open new terminal in the main KittyProtocol catalog
13. Run a hub via command:
    
    ```bash
    go run ./cmd/hub
    ```

14. Run cli clients via command:
       
    ```bash
    go run ./cmd/client_cli
    ```

15. Build gui clients .exe files via command:
    
    ```bash
        go build -x -o meowssenger.exe ./gui_src
    ```

    And the double click on your `mewssenger.exe` file visibile in file explorer.
