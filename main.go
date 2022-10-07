package main

import (
	"embed"
	"fmt"
	"github.com/hktalent/scan4all/lib/api"
	"github.com/hktalent/scan4all/lib/util"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
)

//go:embed config/*
var config embed.FS

// Version can be set at link time to override debug.BuildInfo.Main.Version,
// which is "(devel)" when building from within the module. See
// golang.org/issue/29814 and golang.org/issue/29228.
var Version string

func main() {
	//os.Args = []string{"", "-host", "http://192.168.0.109", "-v"}
	//os.Args = []string{"", "-host", "http://127.0.0.1", "-v"}
	//os.Args = []string{"", "-list", "7b8fa7a85f9f6ae6f9178504d2202666fb8dc772.xml", "-v"}

	runtime.GOMAXPROCS(runtime.NumCPU())
	util.DoInit(&config)
	// set version
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		util.Version = buildInfo.Main.Version
	} else {
		Version = util.Version
	}

	szTip := ""
	if util.GetValAsBool("enableDevDebug") {
		// debug 优化时启用///////////////////////
		go func() {
			szTip = "Since you started http://127.0.0.1:6060/debug/pprof/ with -debug, close the program with: control + C"
			fmt.Println("debug info: \nopen http://127.0.0.1:6060/debug/pprof/\n\ngo tool pprof -seconds=10 -http=:9999 http://localhost:6060/debug/pprof/heap")
			http.ListenAndServe(":6060", nil)
		}()
		//////////////////////////////////////////*/
	}
	api.StartScan(nil)
	log.Printf("wait for all threads to end\n%s", szTip)
	util.Wg.Wait()
	util.CloseAll()
}
