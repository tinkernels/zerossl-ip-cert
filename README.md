# zerossl-ip-cert &middot; [![License](https://img.shields.io/hexpm/l/plug?logo=Github&style=plastic)](https://github.com/tinkernels/zerossl-ip-cert/blob/master/LICENSE)

zerossl-ip-cert is a automation tool for issuing ZeroSSL IP certificates.

* Use ZeroSSL [REST API](https://zerossl.com/documentation/api/)  to implement certificate issuing.
* Mainly made for **IP** certificates.
* Call external program for automatically verification.
* Painless certificate renewal.

## Installation

* Package zerossl-ip-cert contains ZeroSSL [REST API](https://zerossl.com/documentation/api/) client, one can
  just `go get github.com/tinkernels/zerossl-ip-cert` and import it to use the client.
* To build static executables, clone this repository and `make release` , or you can make your desire target
  binary, just take a look at the [Makefile](https://github.com/tinkernels/zerossl-ip-cert/blob/master/Makefile).

## Usage

zerossl-ip-cert rely on configuration file to run. To archive the goal of issuing certificate automatically, you need do
some additional work, saying the external hook.

### Usage Info

```
Version: 0.2.1

Usage: zerossl-ip-cert [ -renew ] -config CONFIG_FILE

  -config string
        Config file
  -renew
        Renew existing certs only
```

### Configuration File

You can find a sample configuration
file [here](https://github.com/tinkernels/zerossl-ip-cert/blob/master/exec/sample-config.yaml), with enough comments
in it. And also a state record sample
file [here](https://github.com/tinkernels/zerossl-ip-cert/blob/master/exec/sample-current.yaml), just for troubleshooting.

### External Hook

zerossl-ip-cert use `HTTP_CSR_HASH` validation method to verify domains (including ip address surely), get more
information from the ZeroSSL official [documentation](https://zerossl.com/documentation/api/verify-domains/).

So you should have a http server running and prepare hook programs to finish the domain verification.

* **verify-hook** will be called before domain verification, some environment variables will be passed to it.

  `ZEROSSL_HTTP_FV_HOST` stands for listening host, here will be ip address.

  `ZEROSSL_HTTP_FV_PATH` stands for url path, where verification content will locate.

  `ZEROSSL_HTTP_FV_PORT` stands for listening port, as far as I know, ZeroSSL only reach 80 port of your http server.

  `ZEROSSL_HTTP_FV_CONTENT` stands for validation content, ZeroSSL will check it when domain verification started.

  And a sample script for nginx can be
  found [here](https://github.com/tinkernels/zerossl-ip-cert/blob/master/exec/sample-nginx-verify-hook.sh).

* **post-hook** will be called after certification downloading, and some other environment variables will be pass to it.

  `ZEROSSL_CERT_FPATH` stands for the store path of certificate.

  `ZEROSSL_KEY_FPATH` stands for the store path of private key.

  And a sample script for nginx can be
  found [here](https://github.com/tinkernels/zerossl-ip-cert/blob/master/exec/sample-nginx-post-hook.sh).

## License

[Apache-2.0](https://github.com/tinkernels/zerossl-ip-cert/blob/master/LICENSE)
