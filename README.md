# education

将`GOPATH`设置为`/root/go`,拉取项目：
```shell
cd $GOPATH/src && git clone https://github.com/Rjhnihao/collection.git
```

在`/etc/hosts`中添加：
```shell
vim /etc/hosts
```

然后在文件添加以下内容：
```
127.0.0.1  orderer.example.com
127.0.0.1  peer0.org1.example.com
127.0.0.1  peer1.org1.example.com
```

添加依赖：
```shell
cd collection && go mod tidy
```

运行项目：
```shell
make kill
make clean
make
./collection
```
后台运行：
```shell
nohup ./collection > collection.log 2>&1 &
```
