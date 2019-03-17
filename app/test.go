package main

import (
	"fmt"
	"gearbox"
	"gearbox/cache"
	"gearbox/host"
)

func main() {
	gb := gearbox.NewApp(&gearbox.Args{
		HostConnector: host.GetConnector(),
	})

	store := cache.NewCache(gb.GetHostConnector().GetCacheDir())

	_, err := store.Get("config")
	if err != nil {
		fmt.Println(err)
	}

	err = store.Set("config", gb.GetConfig().Bytes(), "1s")

	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := store.Get("config")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(gearbox.UnmarshalConfig(c).About))

}
