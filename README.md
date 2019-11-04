# Go Project Example

This is an example for Go project.

The motivaition of this project is for me to learn and widen my limited knowledge about programming, project design, and concepts implementation. In this project, we will try to implement business logic/flow into Go program for various use-cases.

Some of them might not follow existing specs/standards, feel free to open an issue, and please let me know.

## The Project

The project theme is `Property`. We will try to build a Property application, where people able to search and book the property.

### Use-cases

1. Users were able to register and log in.
2. Users were able to register their Properties.
3. Users were able to book a Property.
4. Users were able to receive notifications and have a notification inbox.

## Designing Project For Industrial Programming

What is industrial programming? [Peter Bourgon](https://peter.bourgon.org/go-for-industrial-programming/) explain the terms, as:

- In a startup or corporate environment.
- Within a team where engineers come and go.
- On code that outlives any single engineer.
- Serving highly mutable business requirements.

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
- `--log` define the log configuration for the project. The flag contains comma separated value:
    - `--log=file=./to_your_file.log` the location of log file.
    - `--log=level=debug|info|warning|error|fatal` the level of log.
    - `--log=color=1` to set if color to console/terminal is enabled.
- `--debug` define the debug configuration for the project. The flag contains comma separated value:
    - `--debug=devserver=1` to turn on the dev server.
    - `--debug=testconfig=1` to test the configuration of the project.

### Configuration

For applicatoin configuration, config-file is used and located on the root project directory. The [configuration-file](./project.config.toml) is written in `toml` grammar.

The configuration value in the `configuration-file` is embed as environment-variable value. The value will be replaced by the environment-variable value in runtime or when program started. To help the process, the project use the help of [environment-variable](./project.env.toml) file.

The mixed of `configuration-variable` and `environment-variable` is used to help people in project to see what configuration structure is exists within the project, and also able to dynamically changed based on the environment variables.

### Environment State

The project have no environment state. Different flags and configuration value is used in different environment.

Some projects is using environment state, for example `dev`, `staging`, and `production`. From the experience, this considered harmful for the program itself as people tend to abuse the state for many things. By using the state, people in the project are cutting edges and make many conditional expression for many use-cases. In the end, it leads to bugs and edge-cases to the whole product and harder to maintain.

For example in code:

```go
if env.IsDevelopment() {
    // do something here
}
```

Worse, some projects might use the state for configuration. When the configuration is gated by the environment state. It might introduce another problem, because configuration for each environment might have different parameters and value.

For example in configuration with spec `project.{environment_name}.config.toml`:

- project.dev.config.toml
- project.staging.config.toml
- project.production.config.toml

The approach above usually used to address different configuration needs in each environment. When database configuration is totally different from `dev` and other environment, doing some migration and some configuration is no longer needed, or special configuration that only exists within an environment. This all valid use-cases and the given solution works. Until the configuration is become too long and different for each environments, and turning to problems for the maintainers.

As sometimes we need to run with some special configuration in non-production or in production environment, this might be able to achieved by using the combination of flags and configuration-file. Value from flags and configuration is more clear and straightforward than `IsEnvrionment`, and can be used to check whether we have the right design choices, do we have too many hacks? Why? For whatever reason, the flags/configuration between environments should stay the same, to maintain program consistency.

But, in the end, it depends on each project policies and governance.

## Project Structure

### Cmd

All Go main programs is located in `go_project_example/cmd/*` folder.

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
- check current configuration value

**Dev Server**

Dev server is a server for experimental purpose, and should be enabled with a spesific flag. This server should not be triggered in production environment.

Use-case for dev server:

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
