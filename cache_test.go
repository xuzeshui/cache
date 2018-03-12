package cache

import (
	//	"os"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	mc, err := NewCache("memory", &MemoryConf{CheckPeriodMs: 500, TimeoutFn: func(key string, value interface{}) {
		t.Logf("[%s] TimeoutFn: %s ==> %v", time.Now().Format("2016-01-02 15:04:05"), key, value)
	}})
	if err != nil {
		t.Error("init err")
	}
	t.Logf("start at %s", time.Now().Format("2016-01-02 15:04:05"))
	timeoutDuration := 10 * time.Second
	if err = mc.Put("xuzeshui", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !mc.IsExist("xuzeshui") {
		t.Error("check err")
	}

	if v := mc.Get("xuzeshui"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if mc.IsExist("xuzeshui") {
		t.Error("check err")
	}

	if err = mc.Put("xuzeshui", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if err = mc.Incr("xuzeshui"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := mc.Get("xuzeshui"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = mc.Decr("xuzeshui"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := mc.Get("xuzeshui"); v.(int) != 1 {
		t.Error("get err")
	}
	mc.Delete("xuzeshui")
	if mc.IsExist("xuzeshui") {
		t.Error("delete err")
	}

	//test GetMulti
	if err = mc.Put("xuzeshui", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !mc.IsExist("xuzeshui") {
		t.Error("check err")
	}
	if v := mc.Get("xuzeshui"); v.(string) != "author" {
		t.Error("get err")
	}

	if err = mc.Put("xuzeshui1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !mc.IsExist("xuzeshui1") {
		t.Error("check err")
	}

	vv := mc.GetMulti([]string{"xuzeshui", "xuzeshui1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}
}

//func TestFileCache(t *testing.T) {
//	mc, err := NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":2,"EmbedExpiry":0}`)
//	if err != nil {
//		t.Error("init err")
//	}
//	timeoutDuration := 10 * time.Second
//	if err = mc.Put("astaxie", 1, timeoutDuration); err != nil {
//		t.Error("set Error", err)
//	}
//	if !mc.IsExist("astaxie") {
//		t.Error("check err")
//	}
//
//	if v := mc.Get("astaxie"); v.(int) != 1 {
//		t.Error("get err")
//	}
//
//	if err = mc.Incr("astaxie"); err != nil {
//		t.Error("Incr Error", err)
//	}
//
//	if v := mc.Get("astaxie"); v.(int) != 2 {
//		t.Error("get err")
//	}
//
//	if err = mc.Decr("astaxie"); err != nil {
//		t.Error("Decr Error", err)
//	}
//
//	if v := mc.Get("astaxie"); v.(int) != 1 {
//		t.Error("get err")
//	}
//	mc.Delete("astaxie")
//	if mc.IsExist("astaxie") {
//		t.Error("delete err")
//	}
//
//	//test string
//	if err = mc.Put("astaxie", "author", timeoutDuration); err != nil {
//		t.Error("set Error", err)
//	}
//	if !mc.IsExist("astaxie") {
//		t.Error("check err")
//	}
//	if v := mc.Get("astaxie"); v.(string) != "author" {
//		t.Error("get err")
//	}
//
//	//test GetMulti
//	if err = mc.Put("astaxie1", "author1", timeoutDuration); err != nil {
//		t.Error("set Error", err)
//	}
//	if !mc.IsExist("astaxie1") {
//		t.Error("check err")
//	}
//
//	vv := mc.GetMulti([]string{"astaxie", "astaxie1"})
//	if len(vv) != 2 {
//		t.Error("GetMulti ERROR")
//	}
//	if vv[0].(string) != "author" {
//		t.Error("GetMulti ERROR")
//	}
//	if vv[1].(string) != "author1" {
//		t.Error("GetMulti ERROR")
//	}
//
//	os.RemoveAll("cache")
//}
