xtunnel is a 'fork' from [qtunnel](https://github.com/getqujing/qtunnel)

[![Build Status](https://travis-ci.org/elvizlai/xtunnel.svg?branch=master)](https://travis-ci.org/elvizlai/xtunnel)

Install from docker hub:
```
docker pull sdrzlyz/xtunnel
```

```
Usage of ./xtunnel:
  -crypto string
       	encryption method: blank, rc4, rc4-md5, aes256cfb, chacha20, salsa20 (default "blank")
  -listen string
       	xtunnel local listen (default "127.0.0.1:9000")
  -logto string
       	stdout or syslog (default "stdout")
  -mode string
       	run mode: proxy_server, proxy_client, tunnel_server, tunnel_client
  -remote string
       	xtunnel remote backend (default "127.0.0.1:9001")
  -secret string
       	password used to encrypt data (default "")
```