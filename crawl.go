package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

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
	flag.IntVar(&config.NumWorker, "n", 1, "usage: -n [NumWorker].\n")
	flag.IntVar(&config.LenQueue, "l", 1, "usage: -l [LenQueue].\n")
	flag.Parse()

	c = make(chan Worker, config.LenQueue)

	for i := 0; i < config.NumWorker; i++ {
		go run(i)
		fmt.Printf("init %d complete.\n", i)
	}
}

func main() {
	_ = goquery.Document{}

	testCase := []Model{
		{Name: "taobao", URL: "http://www.taobao.com", Labels: []Label{Label{Route: "div #J_SiteNav a", Attrs: []string{"href"}}}},
	}
	for _, m := range testCase {
		c <- m
	}

	var signalChan = make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
}

type Model struct {
	Name   string
	URL    string // 链接
	Labels []Label
}

type Label struct {
	Route string   // 标签路径
	Attrs []string // 属性
	Loop  bool     // 遍历
}

func (m Model) Work(id int) error {
	doc, err := goquery.NewDocument(m.URL)
	if err != nil {
		return err
	}

	for _, label := range m.Labels {
		sel := doc.Find(label.Route)
		for _, attr := range label.Attrs {
			attrValue, ok := sel.Attr(attr)
			if !ok {
				htmlContent, err := sel.Html()
				if err != nil {
					return err
				}
				log.Printf("count find %s, %s, %s\n", label.Route, attr, htmlContent)
				continue
			}
			fmt.Println(attrValue)
		}
	}

	return nil
}

func run(i int) {
	for w := range c {
		err := w.Work(i)
		if err != nil {
			log.Println(i, err)
		}
	}
}
