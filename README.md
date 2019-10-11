# Project Design

This is an example for Go project design.

The design is based on layered design, which contains several layers and components.

## Designing Project For Industrial Programming

What is industrial programming? [Peter Bourgon](https://peter.bourgon.org/go-for-industrial-programming/) explain the terms in a very good way.

At least, in my understanding, it is consists of theese:

```text
- In a startup or corporate environment.
- Within a team where engineers come and go.
- On code that outlives any single engineer.
- Serving highly mutable business requirements.
```

## The Design

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

## Index

1. [Project Design](/docs/project_design.md)
2. [Service Boundaries](/docs/service_boundaries.md)
3. [Tips](/docs/tips.md)
