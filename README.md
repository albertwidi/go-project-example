# Go Project Example

This is an example for Go project.

The motivation behind this project is to learn and broaden my limited knowledge about programming, project design, concepts, algorithm, and especially, Go itself. In this project, I will try to implement business logic/flow into Go program for various use-cases.

## Designing Project For Industrial Programming

What is industrial programming? [Peter Bourgon](https://peter.bourgon.org/go-for-industrial-programming/) explain the terms, as:

- In a startup or corporate environment.
- Within a team where engineers come and go.
- On code that outlives any single engineer.
- Serving highly mutable business requirements.

## The Project

The project theme is `Property`. I will try to build a Property application, where people able to search and book the property.

### Use-cases

1. Users were able to register and log in.
2. Users were able to register their Properties.
    - Register the property detail
    - Upload the property image
3. Users were able to book a Property.
4. Users were able to receive notifications and have a notification inbox.

### Project Stack

This project is using:

1. PostgreSQL for the main database
2. Redis for k/v, user session management, caching. 

## Getting Started

This is guide to get started with this project and installing dependencies required for running this project locally

### Requirements

1. Make
2. Docker
3. Soda CLI from [gobuffalo](https://gobuffalo.io/en/docs/db/toolbox). You can install this by using `make install-deps`

### Create Database And Migrate

To create and migrate the database, we will use `soda CLI` created by `gobuffalo`. The command is wrapper by this [script](/database/setup.sh).

Use this command to fully create and migrate the databse schema:

`make dbup`

### Flags

The following flags is avaiable to help the project configuration and debug parameters.

- `--config_file` define where the configuration file is located
- `--env_file` define where environment variable file is located, this is a helper file to set environment variable.
- `--debug` define the debug configuration for the project. With flag key:
    - `-debug=server=1` to turn on the debug server.
    - `-debug=testconfig=1` to test the configuration of the project.
    - example: `--debug=-debugserver=1-testconfig=1` 

### Configuration

For applicatoin configuration, config-file is used and located on the root project directory. The [configuration-file](./project.config.toml) is written in `toml` grammar.

The configuration value in the `configuration-file` is embed as environment-variable value. The value will be replaced by the environment-variable value in runtime or when program started. To help the process, the project use the help of [environment-variable](./project.env.toml) file. The `environment-variable` files is choosen because it is simpler as we can have multiple files and we can hide/ignore the file for specific use-case, for example, secret value of `client_id` and `client_secret` of some cloud vendor.

The mixed of `configuration-variable` and `environment-variable` is used to help people in project to see what configuration structure is exists within the project, and able to dynamically changed depends on the environment variables value.

Configuration Structure:

- Servers `[object]`:
    - Main `[object]`:
        - Address: address of the main server, for example `localhost:8000`
    - Admin `[object]`:
        - Address: adress of the admin server, for example `localhost:5726`
    - Debug `[object]`:
        - Address `[string]`: address of debug server, for example `localhost:9000`

- Log
    - level: level of the log, `debug|info|warn|error|fatal`
    - file: file location to store the log
    - color: print with color on the terminal

- Resources
    - Object Storage `[array]`
        - [Object Storage Object]
            - name `[string]`: name of the object storage, for example `image`
    - Database `[object]`
        - Connect `[array]`
            - [Connect Object]
                - name: name of the database, for example `user`
                - driver: the driver of database, `mysql|postgres`
                - leader `object`
                    - dsn: the dsn of database leader, for example 
                - replica `object`
                    - dsn: the dsn of database replica, for example
    - Redis `[object]`:
        - Connect `[array]`:
            - [Connect Object]
                - Name `[string]`: name of redis, for example `session`
                - Address `[string]`: address of redis server, for example `localhost:6379`

Resources configuration allows the project to easily add and remove resources. Because, as the project grow, we might need to add more connection to more postgres, redis or other type of database. Instead of handling the connection manually inside the code, a [wrappeer](./internal/kothak/kothak.go) is added to hold all the connection to resources

### Environment State

The project have no environment state. Different flags and configuration value is used in different environment.

Environment state like `dev`, `staging`, and `production` is usually used to check in what environment the program/application is running. From experience, this considered harmful for the program itself, as developer tempted to abuse the state for many things. Developer tempted to abuse the state because the function is available, and sometimes it is the easiest way to accomplish some goals. By using the state, people in the project are cutting edges and create conditional expression for various use-cases. This leads to broken mental model, bugs, and edge-cases to the product which make life harder for the maintainers.

For example, in code:

```go
if env.IsDevelopment() {
    // do something only in dev
}

if env.IsStaging() {
    // do something only in staging
}

if env.IsProduction() {
    // do something only in production
}
```

This environment state, sometimes also used for configuration directive. When the configuration directive is gated by the environment state, another problem occur. Because, configuration for each environment might have different variables and value, and different configuration file can mean different things.

For example, in configuration with spec `project.{environment_name}.config.toml`:

- project.dev.config.toml
- project.staging.config.toml
- project.production.config.toml

Or, imagine if you have many different configurations(with various reasons/decision) using this kind of directive:

- project.dev.config1.toml
- project.staging.config1.toml
- project.production.config1.toml
- project.dev.config2.toml
- project.staging.config2.toml
- project.production.config2.toml

Multiple configurations with environment state directive, usually used to address different configurations in each environments. For example, when a database is pointing to one instance in `dev` but not in `staging`, which completely different. Or, when doing some migration we want to get rid of some configuration variables in some environment. This all are valid use-cases, and the given solution by using the environment state for configuration directive works. Usually, until the configuration is become too long and different for each environments, then turning into problems for the maintainers.

As sometimes we need to run with some special configuration in non-production or in production environment, this might be able to achieved by using the combination of flags and configuration-file. Variables from flags and configuration is more clear and straightforward than `IsEnvrionment`, and can be used to checked the design choices, do we have too many hacks? Why? For whatever reason, the flags/configuration variables between environments should stay the same, to maintain consistency.

But, in the end, it depends/back to each project policies and governance.

## Project Structure

### Cmd

All Go main programs is located in `go-project-example/cmd/*` folder.

### Internal

To be added

## Code Structure

As stated at the top of this document, the design contains several layers and components and very much similar to onion ring or clean architecture attempt.

But, let's talk about the components first.

### Components

1. Server
2. Usecase
3. Repository
4. Entity

#### Server

Is where all the `http` handler exists. This layer is responsible to hold all the `http` handlers and request validation.

**Main Server**

Main server for serving main application. All business use-case handler exists within this server.

**Admin Server**

Admin server is a server for administrational purpose. On by default and should not open to public.

Use-case for admin server:

- `/metrics` endpoint
- `/resource/status` endpoint
- pprof endpoint
- check current configuration value

**Debug Server**

Debug server is a server for experimental purpose, and should be enabled with a spesific flag. This server should not be triggered in production environment.

Use-case for debug server:

- Login bypass
- Serve fileserver for local object storage

#### Usecase

To be added

#### Repository

To be added

#### Entity

To be added

#### Layers

To be added

### Microservice Structure Inside Monolith

When a company become larger and the number of people and teams increase rapidly, it is make sense to adopt microservice architecture for the company... **to be added**

**NOTE** 

Microservice itself is not a silver bullet to begin with, it solves organizational scale problem, teams become more independent and responsible of their own products, etc. But, it also introduces a lot of problem, for example network, state, and other distributed system problems. Before taking this path, its best to understand what problems that we want to solve with microservice gain more knowledge around microservice, and how we want to tackle the problems inside the organization.

## Error Handling

In Go, `errors` are value. This is one of the of the Go's [proverbs](https://go-proverbs.github.io/). So, instead of throwing error, Go program will treat the `error` as value and the `error` should be checked or returned to the caller. 

> Don't just check errors, handle them gracefully - Go's proverbs

> Check the error or return, but not both. - Dave Channey

While handling `error` in Go is very easy and straightforward, the `error` itself is sometimes lacking of context. 

To be continued.

## Misc

### Set Default Timezone

To set default timezone in go programs, we can set an environment variable named `TZ` when running or in our programs before `time` package get called. The value of `TZ` is the location of Timezone. In my case, it is `Asia/Jakarta`.

For example, in code:

```go
os.SetEnv("TZ", "Asia/Jakarta")
```

or, during execution:

```shell
TZ=Asia/Jakarta ./yourbin
```

**What about more than one timezone in a single application?**

You can always use `time` package function called `time.LoadLocation` to create spesific timezone/location time object.

For example:

```go
loc, err := time.LoadLocation("Asia/Jakarta")
if err != nil {
    // handle the error
}

// this function produce time.Now specified to loc or Asia/Jakarta
nowInLoc := time.Now().In(loc)
```

## References

### Go Errors

**Articles**

- https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
- https://blog.golang.org/go1.13-errors

**Videos**

- [GopherCon 2019: Marwan Sulaiman - Handling Go Errors](https://www.youtube.com/watch?v=4WIhhzTTd0Y)

**Code**
- Upspin errors [package](https://github.com/upspin/upspin/blob/master/errors/errors.go)