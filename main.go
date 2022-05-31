package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"owl/biz/android"
	"owl/biz/ws"
	"owl/service"
	"owl/service/method"
	"syscall"
	"time"

	"github.com/pkg/browser"
	adb "github.com/zach-klippenstein/goadb"

	"github.com/gorilla/mux"
)

const (
	VERSION = "1.0.0"
)

func raise(err error) {
	if err != nil {
		panic(err)
	}
}

func healthHandler(_ http.ResponseWriter, _ *http.Request) {
}

var args struct {
	Port int
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(os.Stdout)

	flag.IntVar(&args.Port, "port", 8848, "server listening port")
	rand.Seed(time.Now().Unix())
}

//go:embed html
var f_html embed.FS

func appendPATH() {
	cur, _ := os.Getwd()
	path := os.Getenv("PATH")
	os.Setenv("PATH", path+":"+cur)
}

func main() {
	flag.Parse()
	log.Printf("Owl手机侦察兵 %s\n", VERSION)

	appendPATH()

	log.Println("正在连接ADB服务器..")
	err := android.AdbStart()
	if err != nil {
		log.Println("请先安装adb:", err)
		log.Println("https://dl.google.com/android/repository/platform-tools-latest-windows.zip")
		time.Sleep(time.Second * 10)
		return
	}
	log.Println("连接ADB服务器成功!")

	listenAddr := fmt.Sprintf("127.0.0.1:%d", args.Port)
	router := mux.NewRouter()
	// 健康检查
	router.HandleFunc("/health", healthHandler)
	// 手机截屏
	router.PathPrefix("/screen").Handler(http.FileServer(http.FS(android.GetFS())))
	// websocket
	router.HandleFunc("/ws", ws.Handler)
	go ws.Run()
	// 接口
	service.RegisterMethod(router, "/v1/test", method.Test)
	service.RegisterMethod(router, "/v1/device/info", method.Device_info)
	service.RegisterMethod(router, "/v1/device/plist", method.Device_plist)
	service.RegisterMethod(router, "/v1/device/labelicon", method.Device_labelicon)
	service.RegisterMethodRaw(router, "/v1/export", method.Export)

	// 静态页面
	{
		f, err := fs.Sub(f_html, "html")
		raise(err)
		router.PathPrefix("/").Handler(http.FileServer(http.FS(f)))
		//http.Handle("/static/", http.FileServer(http.FS(f_static)))
	}

	// Interrupt handler.
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport.
	go func() {
		//log.Println("listening:", listenAddr)
		err := http.ListenAndServe(listenAddr, router)
		raise(err)
	}()

	log.Printf("浏览器打开界面: http://localhost:%d\n", args.Port)
	err = browser.OpenURL(fmt.Sprintf("http://localhost:%d", args.Port))
	if err != nil {
		log.Println("open url warning:", err)
	}

	// android
	go func() {
		watcher := android.NewWatcher()
		for event := range watcher.C() {
			log.Println("event:", event)
			if event.NewState == adb.StateOnline {
				ws.Broadcast("online")
			} else if event.NewState == adb.StateDisconnected {
				ws.Broadcast("disconnected")
			}
		}
	}()

	// Run!
	log.Println("exit:", <-errc)

}
