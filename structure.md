```
|--lru/ 
 	|--lru.go   // 实现lru算法，存储Value是一个接口
|--byteView.go  // ByteView实现lru的Value接口，利用copy实现只读
|--cache.go   // 实例化lru，提供加锁的add和get
|--gocache.go // 实现Getter接口，适配用户自定义回调从数据源获取数据的函数；实现group分组，对缓存进行分组，实现从分组中获取缓存：如果存在，返回缓存的值，如果不存在调用用户自定义回调函数。
|--http.go // 提供被其他节点访问的能力
```



