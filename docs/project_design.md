# Design

![design](/docs/images/design.png)

## Services

Services contains a group of service that exists in a repository. For example one `server` might have several kind of services: `user service`, `order service`, `payment service`

The services itself can be related to each others or not at all. For example `order service` might need `user service` help to check wether the user is active or not. Or maybe `payment service` need `order service` to check the status of the order itself.

### Service

The dependencies in a service managed by `interfaces` exists within the service itself. For example the `order service` has `IsUserActive` method for `UserService` interface:

```go
type UserService interface {
	IsUserActive(int64) (bool, error)
}
```

This happens because `order service` needs some help/method from `user service` to check wether the user is active or not. But because the `order service` need `user service`, it doesn't mean the `order service` needs to know all APIs available in the `user service`. `Order service` might only need to know some of the APIs and is abstracted by using the `interface`.

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

When creating a resource, it is recommended to inject the resource dependencies. For example database or redis. It might be true in `microservice` project, we don't need to throw `resource dependencies`(for example db) to the `resource`, and open the dependencies inside the `resource` instead. But this approach is not good for a monolithic project.

A monolithic project usually share dependencies, and if the dependencies is not injected into the `resource`, all `resource` that need the same dependencies to open the same connection again and again. For example `order` and `payment` using the same `redis` to store some data, `order` and `payment` need to connect to the same `redis` independently if the connection is not injected.

### Do

Server:

```go
func Main() error {
    // open connection to redis
    r := redis.Connect(addr)
    orderres := orderresource.New(r)
    paymentres := paymentresource.New(r)
}
```

Order Service:

```go
package order

func New(r redis.Redis) OrderResource {
    // some code
    return ordersvc
}
```
Payment Service:

```go
package payment

func New(r redis.Redis) PaymentResource {
    // some code
    return paymentsvc
}
```

### Instead Of

Server:

```go
func Main() error {
    orderres := orderresource.New()
    paymentres := paymentresource.New()
}
```

Order Service:

```go
package order

func New(r redis.Redis) OrderResource {
    // open connection to redis
    r := redis.Connect(addr)
    // some code
    return orderres
}
```
Payment Service:

```go
package payment

func New(r redis.Redis) PaymentResource {
    // open connection to redis
    r := redis.Connect(addr)
    // some code
    return paymentres
}
```

As we can see in above example, `order service` and `payment service` has the same redis connection. Reusing the connection object would be much simpler instead of re-creating the connection.

## API

Api is a bridge between service and http/grpc service.

In `api`, contains all remote API declaration that available in the application. The service used in the `api` is declared in `api/interface.go`. It is expected to call APIs `Init()` as `api` is using global function and `service` object need to be injected to the `api` in order to use the global function.

## Pkg

Contains all shared libraries used by the repository.

## Server

Is where the `server` application begin.

This is the main `user space`.

## Cmd

Keeps `main.go` simple. It is only executing `server.Main()` for now.

Separating the `main` program itself(maybe flags declaration and recover) from the `server` program.