# Service Boundaries/Domain

In this section, I want to talk about service boundaries or we know as service domain. Especially `service` domain in Industrial Go application.

![service boundaries1](/docs/images/service_boundaries1.png)

## Services/Microservices

In computer services world, a single service may has their own `resource` and  interacting with other service like the picture above. This might also happens in a `microservice` world, where a service can have their own resource in a big cluster of machines, then forming a platform.

Services in picture above are interacting with each others. While it might sounds weird in `microservice` world, for two service to `depends` on each other, it really depends on the domain of the service itself. We need to know what is the domain of the service or what is the purpose of the service? How big the service is? Is it a flaw in the design? Etc.

![microservice interaction](/docs/images/microservice_interaction.png)

As the `microservice` grow, the interaction between microservice become more complex. Why `service B` need to contact `service C` to be able to reach `service F`, what is the purpose of `service C`? Is `service C` really needed? The domains can become really blurry sometimes, it becomes worse as we design it over network and no tools are able to validate our design.

External service dependencies are common, and sometimes service depends to each others. But, in `microservice` world, this can be a serious problem. Cyclic problem might coming from flaw when designing the `microservice` dependencies, and might lead to a very serious problem. For example, we have service `A`, `B, and `C`. And the conditions are:

1. Service `A` can be called from public network and depends on service `B`
2. Service `B` can be called from public network and depends on service `A` and `C`
3. Service `C` is a standalone service

Cyclic dependencies in microservice picture:

![microservice cyclic dependencies](/docs/images/microservice_cylic.png)

Things can go from good to bad very fast in this kind of dependencies. It is good now we have seperated responsibility between services, but things can go wrong when `service A` need to call `service B` and somehow `service B` need to call `service A` to answer `service A` or (A->B->A->A). 

This design leave so many questions for us. First, why `service A` and `service B` depends on each others? Its understandable if `service A` depends on `service B`or vice versa, but what data did `service B` needs from `service A`? If we need data from `service A`, why don't just call `service A` directly? Lastly, if something bad happens on one service, it can be bad for both service.

In this case, a design review needs to be done, and the the cyclic calls need to get removed immediately.

## Monolith Application

Things get interesting when developing a monolith Go application, because Go prevent us from doing a cyclic import.

![import_cyclic](/docs/images/import_cyclic.png)

Package import in the picture above is not allowed in Go, and will be recognized as error by the compiler.

But why this is a big deal in the first place? The point is, we can create a structure in our Go application and use the compiler feature as our tool to shape a better design. Import cycle is not just an error, it is a sign of something, a sign that we might miss something when we design our application.

![service package](/docs/images/service_package.png)

Structuring several `services` inside a Go application means we need to create a certain package domain for each `service`. But if one `service` package only do a certain `domain`, how do they interact? How the dependencies happens?

In the `microservice` structure, we learn to call other `service` by APIs, and it should not be different in terms of structure in monolith application. For example `service B` needs to call `service B`:

![service package dependencies](/docs/images/service_package_dependencies.png)

As we can see in above picture, we can inject `service C` as `service B` dependency and accessing `service C` by using its public API without getting care about `service C` dependencies. 

And if an import cyclic happens:

![service package cyclic](/docs/images/service_package_cyclic.png)

The compiler will tell us that we have an import cyclic, and need to fix it. To make our program running, we simply need to revisit the design decission we made and fix the problem. 

Microservice is good, but we can't take it for granted, and all microservices are coming from monolith application like this.

## Microservice Application

Will add soon