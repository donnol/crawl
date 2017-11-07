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
		log.Printf("init No.%d worker.\n", i)
	}

	log.Printf("init completed.\n\n")
}

func main() {
	_ = goquery.Document{}

	testCase := []Model{
		{
			Name: "taobao", URL: "http://www.taobao.com", Labels: []Label{
				Label{
					Route: "div #J_SiteNav a",
					Attrs: []string{
						"href",
					},
				},
				Label{
					Route: "img",
					Attrs: []string{
						"src",
					},
					Flag: "all",
				},
				Label{
					Route: ".screen-outer .service-bd li a",
					Attrs: []string{
						"href",
					},
					Flag: "list",
				},
			},
		},
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
	Flag  string   // all: 所有标签；list: 标签列表
}

func (m Model) Work(id int) error {
	doc, err := goquery.NewDocument(m.URL)
	if err != nil {
		return err
	}

	findAttr := func(s *goquery.Selection, attrs []string) error {
		for _, attr := range attrs {
			attrValue, ok := s.Attr(attr)
			if !ok {
				htmlContent, err := s.Html()
				if err != nil {
					return err
				}
				log.Printf("count find %s, %s\n", attr, htmlContent)
				continue
			}
			fmt.Println(attrValue)
		}
		return nil
	}
	for _, label := range m.Labels {
		sel := doc.Find(label.Route)
		// selHtml, err := sel.Html()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println(selHtml)

		// 列表查找
		if label.Flag == "list" {
			sel.Each(func(i int, s *goquery.Selection) {
				err = findAttr(s, label.Attrs)
				if err != nil {
					log.Println(err)
				}
			})
		} else if label.Flag == "all" {
			// TODO
		} else {
			err = findAttr(sel, label.Attrs)
			if err != nil {
				log.Println(err)
			}
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
