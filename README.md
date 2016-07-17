# dummi - dummy microservice

A tiny golang application simulating a microservice. It

- listens for inbound connections
- connects to other named microservices specified on the command line
- sends messages containing a sequence number
- logs sent and received messages
- attempts to reconnect on failure

Included are two Docker Compose files:

- docker-compose.yml creates a tiny app consisting of three
  interconnected containers
- docker-compose-weavesock.yml creates a simulation of the
[weavesock demo](https://github.com/weaveworks/weaveDemo)

## WeaveSock Demo

Unlike the real thing, this doesn't have an HTML endpoint you can
point your browser at. Other than that, it is architecturally
identical and a neat little example to play with in order to
familiarize yourself with
[Weave Scope](https://www.weave.works/products/weave-scope/).

Simply run it with

    docker-compose -f docker-compose-weavesock.yml up


Here are a few things to try, all using Weave Scope

1. understand the app architecture by looking at the toplogy shown
2. get an overview of memory and cpu usage
3. explore the 'front_end' container in more detail:
  - see details of the inbound and outbound connections
  - check the command line arguments
  - look at the logs to see the inbound and outbound messages
4. stop the 'orders' container
  - see it disappear from the topology view
5. check the logs of the 'front_end' container
  - observe errors
6. try pinging 'orders' from the 'front_end' container via a shell
  - this will fail
7. restart the 'orders' container
  - see it re-appear in the topology and connections getting
    restablished
8. check the logs of the 'front_end' container again
  - there should be no more errors
