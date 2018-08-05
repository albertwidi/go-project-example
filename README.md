# Kothak

Kothak is a repository design project for Go.

The aim of the project is to design an Industrial Go Application

## Services

Services contains a group of services exists in a repository. For example one `server` might have several kind of services: `user service`, `order service`, `payment service`

The services itself can be related to each others or not at all. For example `order service` might need `user service` help to check wether the user is active or not. Or maybe `payment service` need `order service` to check the status of the order itself.

### Service

The dependencies in a service managed by `interfaces` exists within the service. For example the `order service` has `IsUserActive` method for `UserService` interface:

```go
type UserService interface {
	IsUserActive(int64) (bool, error)
}
```

This happens because `order service` needs some help/method from `user service` to check wether the user is active or not. But because the `order service` need `user service`, that doesn't mean the `order service` needs to know all APIs available in the `user service`. `Order service` might only need to know some of the APIs and is abstracted by using `interface{}`.

By using `interface` for dependencies to the `service` itself, it is easier to maintain what is needed. It is also easier to test without care much about the depdendencies. We can mock the dependencies or just using `docker` for database/redis test.

All `service` declaration is in `service.go` in every `service` directory.

## Resource

Resource is where dependencies calls happens. Database, redis, api, etc.

Usually, resource is an implementation of `service interface`. For example `order service` need resource to create an order, and have this interface in the `order service`:

```go
type Resource interface {
	CreateOrder(context.Context, Order) error
}
```

When creating a resource, it is recommended to inject the resource dependencies itself. For example database or redis. It might be true in a `microservice` project we not need to throw `resource dependencies`(for example db) to the `resource`, and open the dependencies inside. But this approach is not good for a big monolithic project.

A big monolithic project usually share dependencies, and if the dependencies is not injected into the `resource`, all `resource` that need the same dependencies to open the same connection again and again. For example `order` and `payment` using the same `redis` to store some data, `order` and `payment` need to connect to the same `redis` independently if the connection is not injected.

### Recommended

Server:

```go
func Main() error {
    // open connection to redis
    r := redis.Connect(addr)
    ordersvc := orderservice.New(r)
    paymentsvc := paymentservice.New(r)
}
```

Order Service:

```go
package order

func New(r redis.Redis) OrderService {
    // some code
    return ordersvc
}
```
Payment Service:

```go
package payment

func New(r redis.Redis) PaymentService {
    // some code
    return paymentsvc
}
```

### Instead Of

Server:

```go
func Main() error {
    ordersvc := orderservice.New()
    paymentsvc := paymentservice.New()
}
```

Order Service:

```go
package order

func New(r redis.Redis) OrderService {
    // open connection to redis
    r := redis.Connect(addr)
    // some code
    return ordersvc
}
```
Payment Service:

```go
package payment

func New(r redis.Redis) PaymentService {
    // open connection to redis
    r := redis.Connect(addr)
    // some code
    return paymentsvc
}
```

as we can see in above example, `order service` and `payment service` has the same redis connection, reusing the connectin object would be much simpler.

## Transport

## Pkgs

Contains all shared libraries used by the repository.

## Server

This is where the applications began

## Cmd

Kothak keeps `main.go` simple. It is only executing `server.Main()`.

Separating the `main` program itself(maybe flags declaration and recover) from the `server` program.