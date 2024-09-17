# Go Web Server

This is a WIP web server project meant to learn the net/http package after Go 1.22.

View the git Commit History in order to see the lesson objectives.

> [!IMPORTANT]
> This is not meant to be used in any official capacity, so the `.env` file is included for simplicity.





## Installation

1. Clone this repository:
    ```terminal
    $ git clone https://github.com/your-username/web-server.git
    $ cd web-server
    ```
2. Compile the web server. We can either compile or run directly with Go, or build a Docker image:

    - (option) Build docker image
        ```terminal
        $ docker build -t web-server .
        ```
        > Container port 8080 is exposed



## Usage

- `make run`    -> runs the web-server
- `make debug`  -> will wipe the DB and run the web-server
- `make test`   -> will run `make debug` and run some cURL commands against the web api




## TODO

- [ ] Create API documentation
