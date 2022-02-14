package shell

import (
	"fmt"
	"testing"
)

func TestBash(t *testing.T) {
	//SetDebug(true)
	handleOutput(Run("mkdir /home/testdir", false))
	handleOutput(Run("touch /home/testdir/1", false))
	handleOutput(Run("touch /home/testdir/2", false))
	handleOutput(Run("ls /home/testdir", false))
	handleOutput(Run("rm -rf /home/testdir", false))
}

func handleOutput(code int, output string, err error) {
	fmt.Println("[Code]:", code, "\n", output)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestRunScrapyJob(t *testing.T) {
	Run("cd /data/git/github/go2o-scrapy/sogou_goods && ls .", true)
	//StdRun("cd /data/git/github/go2o-scrapy/sogou_goods && ls . && ~/.local/bin/scrapy crawl sogou_goods")
}
