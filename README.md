# Simple go Networking Utilities

This repository contains a handful of HTTP servers and clients, userful for testing connection between machines. The utilities are distributed as executables, docker images (available at [hub.docker.com/u/lucasegp](https://hub.docker.com/u/lucasegp)) and docker-compose yaml versions for easy use.

## Projects
### Simple Server and Client
The server/ and client/ directories contain the simplest utilities possible: an HTTP server that responds some text on "/", and a client that gets the "/" page of a server continuously.
### HTTP Proxy
The proxy/ directory contains a "man in the middle" like server that logs the full coommunication between a client and a server (can be highly userful to debug APIs and be certain about security).
### HTTP File Transfer
The httptransfer/ directory contains a server similar to server/, however it communicates files, and not text messages. It also can receive files via a POST request (but the files written to and read from are different, this is intentional for testing purposes).

## How to Use the Projects
All projects can be compiled or executed locally going to the project's directory and running `go build .` or `go run .`. It is also possible to use `docker-compose up name_of_project` or `docker pull lucasegp/name_of_project` and `docker run lucasegp/name_of_project` to use containerized linux versions. 

### Documentation on httptransfer's Files
The HTTP file transfer utility has a key difference from others because it deals with volumes (to access the files to send or receive data). Because of that, if the host has SELinux enabled, some flags need to be set to solve "permission denied" errors on the data directory (by default ./httptranfer/data).
To do that, after creating the data directory and adding a "get_file", just one command is needed to enable read/write access from docker to it:
```
sudo chcon -R -t svirt_sandbox_file_t ./httptranfer/data/
```
This command modifies the 'SELinux type' of the directory (and any files inside it) to one above, enabling it's access from containers, as described in [this Red Hat's article](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux_atomic_host/7/html/container_security_guide/docker_selinux_security_policy) about docker's SELinux security policy and [this other one](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/security-enhanced_linux/sect-security-enhanced_linux-working_with_selinux-selinux_contexts_labeling_files) about SELinux context manipulation.
