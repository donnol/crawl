package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var c chan Worker

type Config struct {
	NumWorker int
	LenQueue  int
}

type Worker interface {
	Work(id int) error
}

func init() {
	var config Config
	flag.IntVar(&config.NumWorker, "n", 10, "usage: -n [NumWorker].\n")
	flag.IntVar(&config.LenQueue, "l", 100, "usage: -l [LenQueue].\n")
	flag.Parse()

	c = make(chan Worker, config.LenQueue)

	for i := 0; i < config.NumWorker; i++ {
		go run(i)
		fmt.Printf("init %d complete.\n", i)
	}

	var signalChan = make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
}

func main() {
	_ = goquery.Document{}

	testCase := []Model{
		{Name: "jd"},
		{Name: "je"},
		{Name: "jf"},
		{Name: "jg"},
		{Name: "jh"},
	}
	for i := 0; i < 1000; i++ {
		for _, m := range testCase {
			c <- m
		}
	}
	time.Sleep(time.Second * 2)
}

type Model struct {
	Name   string
	URL    string   // 链接
	Labels []string // 标签路径
	Attrs  []string // 属性
}

func (m Model) Work(id int) error {
	fmt.Printf("%d: %s\n", id, m.Name)
	return nil
}

func run(i int) {
	for w := range c {
		err := w.Work(i)
		if err != nil {
			log.Fatal(i, err)
		}
	}
}
