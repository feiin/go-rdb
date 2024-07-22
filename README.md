
# go-rdb

尝试用Go写个类似redis server的KV内存数据库,空余时间学习写着玩...


## test

server:

```
$ go run main.go

```

client:
```
$ redis-cli -h 127.0.0.1 -p 6378

127.0.0.1:6378> set aaa 44
127.0.0.1:6378> get aaa
"44"

```

## 实现

- [x] 使用ART做内存索引
- 实现数据结构
    - [x] string
    - [ ] list
    - [ ] sortset
    - [ ] set
    - [ ] hash

