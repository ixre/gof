package util

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// 监听文件变化,并进行处理
func FsWatch(h func(fsnotify.Event), directory ...string) {
	//创建一个监控对象
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()
	//添加要监控的对象，文件或文件夹
	for _, path := range directory {
		fi, err := os.Stat(path)
		if err == nil {
			err = watch.Add(path)
			if fi.IsDir() {
				//通过Walk来遍历目录下的所有子目录,并添加到监听中
				err = filepath.Walk(path, func(path string, info os.FileInfo, _ error) error {
					if info.IsDir() {
						return watch.Add(path)
					}
					return nil
				})
			}
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	//log.Printf("watch directory: %s\n", directory)
	// 监控文件更改,如果更改则生成代码
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create ||
						ev.Op&fsnotify.Rename == fsnotify.Rename ||
						ev.Op&fsnotify.Write == fsnotify.Write ||
						ev.Op&fsnotify.Remove == fsnotify.Remove {
						h(ev)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()
	//循环
	select {}
}
