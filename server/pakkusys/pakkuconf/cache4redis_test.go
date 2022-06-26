package pakkuconf

import (
	"fmt"
	"testing"

	"github.com/wup364/pakku"
	"github.com/wup364/pakku/ipakku"
)

func TestCache(t *testing.T) {
	app := pakku.NewApplication("T").EnableCoreModule().BootStart()
	var conf ipakku.AppConfig
	if err := app.GetModuleByName("AppConfig", &conf); nil != err {
		panic(err)
	}
	rdc := NewRedisCache()
	rdc.Init(conf, "T")
	rdc.RegLib("T", 60)
	// Set
	if err := rdc.Set("T", "key-001", "value-001"); nil != err {
		panic(err)
	}
	// Get
	var str string
	if err := rdc.Get("T", "key-001", &str); nil != err {
		panic(err)
	} else {
		fmt.Println(str)
	}
	var str1 string
	if err := rdc.Get("T", "key-002", &str1); nil != err && err != ipakku.ErrNoCacheHit {
		panic(err)
	} else {
		fmt.Println(str1)
	}
}
