package main

type Broadcaster struct {
	input       chan string
	subscribe   chan chan string
	unsubscribe chan chan string
}

func NewBroadcaster() *Broadcaster {
	input := make(chan string)
	subscribe := make(chan chan string)
	unsubscribe := make(chan chan string)

	b := Broadcaster{input, subscribe, unsubscribe}
	go b.run()

	return &b
}

func (b *Broadcaster) run() {
	collection := make(map[chan string]struct{})
	for {
		select {
		case msg := <-b.input:
			for ch := range collection {
				ch <- msg
			}
		case subCh := <-b.subscribe:
			collection[subCh] = struct{}{}
		case unsubCh := <-b.unsubscribe:
			delete(collection, unsubCh)
			close(unsubCh)
		}
	}
}

func (b *Broadcaster) Publish(msg string) {
	b.input <- msg
}

func (b *Broadcaster) Subscribe() chan string {
	ch := make(chan string)
	b.subscribe <- ch
	return ch
}

func (b *Broadcaster) Unsubscribe(ch chan string) {
	b.unsubscribe <- ch
}
