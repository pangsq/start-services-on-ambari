# start-services-on-ambari

Send a request to ambari-server for making all services to be started

## how to compile

```shell
dep ensure
go build
```

## how to run

```
./start-services-on-ambari --ambari http://192.168.127.131:30888/api/v1/clusters/delphini --user <admin_user> --password <password>
```