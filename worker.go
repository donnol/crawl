package main

import "log"

type Worker interface {
	Work(id int) error
}

func run(numWorker int) {
	for i := 0; i < numWorker; i++ {
		go func(i int) {
			for w := range c {
				err := w.Work(i)
				if err != nil {
					log.Println(i, err)
				}
			}
		}(i)
		log.Printf("init No.%d worker.\n", i)
	}
}
