/**
* @Author: CiachoG
* @Date: 2020/5/26 15:13
* @Description：
 */
package main

import (
	"cache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"golang": "666",
	"java":   "777",
	"python": "888",
}

func main() {
	cache.NewGroup("lang", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[DB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "0.0.0.0:23451"
	peers := cache.NewHttpPool(addr)
	log.Println("go cache is run at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
