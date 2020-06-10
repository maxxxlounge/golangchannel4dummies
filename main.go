package main

import (
	"fmt"
	"sync"
)

type NumberGenerator struct {
	sec int
	sync.Mutex
}

type SharedMaker struct {
	Nr *NumberGenerator
}

func (nr *NumberGenerator) GetSectional() int {
	// lokko la risorsa condivisa con i mutex per impedire che 2 elemnti accedano alla risorso in contemporanea
	nr.Lock()
	defer nr.Unlock()
	nr.sec++
	return nr.sec
}

func ProcessMessage(wg *sync.WaitGroup, nr *NumberGenerator, in, out chan int) {
	defer wg.Done()
	//chiamo la mia funzione per ogni ingresso del canale in entrata
	for _ = range in {
		out <- nr.GetSectional()
	}
}

func main() {
	/*t0 := time.Now()
	defer func() {
		fmt.Println(time.Since(t0))
	}()*/
	// definisco i 2 canali con i relativi waitgroup per aspettare che siano processati tutti gli elementi
	in := make(chan int, 1000)
	out := make(chan int, 10000)
	wg := &sync.WaitGroup{}
	wgout := &sync.WaitGroup{}

	nr := NumberGenerator{}

	// aggiungo 1 al mio waitgroup e ed elaboro gli elementi in uscita che provengono dal canale di out
	wgout.Add(1)
	go func(wgout *sync.WaitGroup, out chan int) {
		defer wgout.Done()
		for o := range out {
			fmt.Printf("%v\n", o)
		}
	}(wgout, out)

	// definisco 10 worker per processare gli nput
	for i := 0; i < 10; i++ {
		//per ogni worker aggiungo un punto al waitgroup e inizio a processare il messaggio passando il canale in ingresso e ein uscita
		wg.Add(1)
		go ProcessMessage(wg, &nr, in, out)
	}

	//inizio ad inserire gli eelmenti da elaborare nel channel in input
	for i := 0; i < 100000; i++ {
		in <- i
	}
	//chiudo il channel in input una volta inseriti tutti gli elementi (pallina)
	close(in)
	//attendo che sia finito il waitgroup associato al channel 1 e quindi sono sicuro di aver inserito nel chan tutti gli elementi
	wg.Wait()
	//inserisco l'ultimo elemento nel channel di uscita (pallina) come ultimo elemento da oprocessare
	close(out)
	// attendo la chiusura del channel per chiudere l'esecuzione
	wgout.Wait()

}
