# map-chan
Go Parallelism with Ordered Concurrency

# Introduction
A Go utility library designed to enhance parallel execution capabilities in asynchronous processes, while ensuring an ordered sequence according to a specific degree of differentiation. This library is suitable for scenarios such as network message forwarding, consumption queue concurrent consumption and the like.

## Features:
**Asynchronous Task Execution**

Ensures that messages with the same key are executed in sequence while supporting hot updates for concurrency and buffer size.

**Asynchronous Message Pushing**

Manages a service discovery pool that ensures messages with the same key are forwarded to the same connection while supporting hot updates for connection addresses and number of connections.

