# Go Project Example

This is an example for Go project.

The motivaition of this project is for me to learn and widen my limited knowledge about programming, project design, and concepts implementation.

## The Project

This project is about implementing business logic/flow into the Go program for various usecases.

Some of them might not follow existing specs/standards, feel free to open an issue and please let me know.

Project theme: `Property`

Items covered in this project:

- User reigstration
- User authentication
- Property registration
  - Uploading property image
- Booking And Order creation
- Notification inbox

## Designing Project For Industrial Programming

What is industrial programming? [Peter Bourgon](https://peter.bourgon.org/go-for-industrial-programming/) explain the terms in a very good way.

At least, in my understanding, it is consists of theese:

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

#### Usecase

To be added

#### Repository

To be added

#### Entity

To be added

#### Layers

To be added
