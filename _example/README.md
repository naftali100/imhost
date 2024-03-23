# Usage

```shell
# Spawn 5 containers and counts the number of host names that
# are different.
./up.sh 5
```

```shellsessionc
$ ./up.sh 5
- docker compose up: starting 5 instances of imhost service.
[+] Running 6/7
 ⠼ Network example_default           Created                                                                                                                                                                              2.4s
 ✔ Container example-imhost-1        Started                                                                                                                                                                              0.8s
 ✔ Container example-imhost-5        Started                                                                                                                                                                              2.0s
 ✔ Container example-imhost-3        Started                                                                                                                                                                              1.4s
 ✔ Container example-imhost-2        Started                                                                                                                                                                              1.7s
 ✔ Container example-imhost-4        Started                                                                                                                                                                              1.1s
 ✔ Container example-loadbalancer-1  Started                                                                                                                                                                              2.1s
- stress test: starting stress test with 5 instances.
        Searching for 5 hosts.
        Found 5 hosts .....
        OK:found all hosts
        List responses:
          #01: "Hello from host: 9630feeaa63c"
          #02: "Hello from host: bc03ae785b34"
          #03: "Hello from host: f26a9bdb6e4d"
          #04: "Hello from host: d87ac3a9d07f"
          #05: "Hello from host: 2b63d018a7ee"
- docker compose down: done.
```
