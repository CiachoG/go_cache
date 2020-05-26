/**
* @Author: CiachoG
* @Date: 2020/5/25 18:10
* @Description：
 */
package cache

import (
	"fmt"
	"log"
	"sync"
)

//定义一个函数类型 F，并且实现接口 A 的方法，
//然后在这个方法中调用自己。这是 Go 语言中将
//其他函数（参数返回值定义与 F 一致）转换为接
//口 A 的常用技巧。

type Getter interface {
	Get(key string) ([]byte, error)
}
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//Group 缓存分组,
type Group struct {
	name      string //不同缓存有不同的名字
	getter    Getter //缓存未命中时的回调
	mainCache cache  //并发缓存
}

var (
	mu     sync.RWMutex //对groups读写进行保护
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

//Get 从缓存中获取数据，如果不存在调用用户自定义回调函数
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok { //如果缓存中存在
		log.Printf("the key of %v is hit", key)
		return v, nil
	}
	return g.load(key) //缓存不存在
}
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key) //分布式场景下会调用 getFromPeer 从其他节点获取
}
func (g *Group) getLocally(key string) (value ByteView, err error) {
	bytes, err := g.getter.Get(key) //调用用户自定义数据获取的方法
	if err != nil {
		return ByteView{}, err
	}
	value = ByteView{
		b: cloneBytes(bytes),
	}
	g.populateCache(key, value)
	return value, nil
}
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
