# golang logger

实现log的基本操作，实现按照时间周期写入

## Install 

```sh
# log基本操作，并实现按时间周期写入log
# 按时间周期写入，保证一个周期内，只会写入一次。
# 对于很多写log的情况，我们都需要控制一定的输出频率，避免log文件被写爆掉。
go get -u github.com/ibbd-dev/go-log

# 异步写log
go get -u github.com/ibbd-dev/go-log/async-log
```

## Example

