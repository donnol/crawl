package main

import "log"

type Worker interface {
	Work(id int) error
}

func run(i int) {
	for w := range c {
		err := w.Work(i)
		if err != nil {
			log.Println(i, err)
		}
	}
}
