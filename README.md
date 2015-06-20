json2csv
========

convert a stream of newline separated json data to csv whitch supports json array.

[![Build Status](https://travis-ci.org/SunRunAway/json2csv.png?branch=master)](https://travis-ci.org/SunRunAway/json2csv)


Installation
============

If you have a working golang install, you can use `go get`.

```bash
go get github.com/SunRunAway/json2csv
```

The binary build is here <https://github.com/SunRunAway/json2csv/releases>

Usage
=====

```
usage: json2csv
  -d=",": delimiter used for output values
  -i="": /path/to/input.json (optional; default is stdin)
  -o="": /path/to/output.json (optional; default is stdout)
  -p=true: prints header to output
```

To convert:

```json
{"user": {"name":"jehiah", "password": "root"}, "remote": [{"ip":"127.0.0.1", "port": 8080}]}
{"user": {"name":"jeroenjanssens", "password": "123"}, "remote": [{"ip":"10.0.0.1", "port": 1080}]}
{"user": {"name":"unknown", "password": ""}, "remote": [{"ip":"10.0.0.1", "port": 27017}, {"ip":"10.0.0.2", "port": 27017}]}
```

to:

```
remote[0].ip,remote[0].port,remote[1].ip,remote[1].port,user.name,user.password
127.0.0.1,8080,,,jehiah,root
10.0.0.1,1080,,,jeroenjanssens,123
10.0.0.1,27017,10.0.0.2,27017,unknown,
```
    
you would either

```bash
json2csv -i input.json -o output.csv
```

or

```bash
cat input.json | json2csv > output.csv
```

Thanks
=====

Many thanks to [jehiah](https://github.com/jehiah/json2csv) which i copied some docs and user interface.
