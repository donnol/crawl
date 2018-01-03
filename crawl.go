package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var c chan Worker

type Config struct {
	NumWorker int
	LenQueue  int
}

func init() {
	var config Config
	flag.IntVar(&config.NumWorker, "n", 1, "usage: -n [NumWorker].\n")
	flag.IntVar(&config.LenQueue, "l", 1, "usage: -l [LenQueue].\n")
	flag.Parse()

	c = make(chan Worker, config.LenQueue)

	run(config.NumWorker)

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Printf("init completed.\n\n")
}

func main() {
	testCase := []Model{
		Model{
			Name: "taobao", URL: "https://www.taobao.com", Dynamic: true, Labels: []Label{
				Label{
					Route: "div #J_SiteNav a",
					Attrs: []string{
						"href",
					},
				},
				Label{
					Route: "img",
					Attrs: []string{},
					Flag:  "all",
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
		Model{
			Name: "wenku", URL: "https://wenku.baidu.com/view/197534bd011ca300a7c39019", Dynamic: true, Labels: []Label{
				Label{
					Route: "#WkDialogDownDoc.dialog-container.dialog-org.border-none.dialog-top.doc-title",
				},
				Label{
					Route: ".reader-word-layer", //.reader-word-s1-0.reader-word-s1-2
					Attrs: []string{},
					Flag:  "all",
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
	Name    string
	URL     string // 链接
	Dynamic bool   // 动态网站
	Labels  []Label
}

type Label struct {
	Route string   // 标签路径
	Attrs []string // 属性
	Flag  string   // all: 所有标签；list: 标签列表
}

func (m Model) Work(id int) error {
	var doc *goquery.Document
	var err error
	if m.Dynamic {
		log.Println("=== begin phantom.")
		content, err := phantom(m.URL)
		if err != nil {
			return err
		}
		fmt.Printf("content is: %s\n", content)
		if len(content) == 0 {
			return errors.New("empty content.")
		}
		doc, err = goquery.NewDocumentFromReader(bytes.NewReader(content))
		if err != nil {
			return err
		}
	} else {
		doc, err = goquery.NewDocument(m.URL)
		if err != nil {
			return err
		}
	}

	findAttr := func(s *goquery.Selection, attrs []string) error {
		for _, attr := range attrs {
			attrValue, ok := s.Attr(attr)
			if !ok {
				_, err := s.Html()
				if err != nil {
					return err
				}
				// log.Printf("count find %s, %s\n", attr, htmlContent)
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
			sel.Each(func(i int, s *goquery.Selection) {
				// fmt.Printf("%#v\n", s)
				for _, node := range s.Nodes {
					// fmt.Printf("%#v\n", node)
					if node.Type == html.ElementNode && node.Data == label.Route {
						for _, element := range node.Attr {
							if element.Key == "src" {
								fmt.Println(element.Val)
							}
							if element.Key == "data-original" {
								fmt.Println(element.Val)
							}
						}
					}
				}
				return
			})
		} else {
			err = findAttr(sel, label.Attrs)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func regexpFind(docHtml string) []string {
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindAllStringSubmatch(docHtml, -1)
	out := make([]string, len(imgs))
	for i := range out {
		out[i] = imgs[i][1]
	}
	// fmt.Println(out)
	return out
}
