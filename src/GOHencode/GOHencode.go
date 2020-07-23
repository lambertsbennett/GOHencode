package main


// GOHencode.go -n PROCS -file CONTIG_FILE -o OUTPUT_FILE


import (
	"flag"
	"fmt"
	"github.com/lambertsbennett/GOHencode/src/SeqOps"
	"runtime"
	"sync"
	"time"
)

func main() {

	var proc int
	flag.IntVar(&proc,"n",2,"Number of processors or threads to leverage.")

	var contigfile string
	flag.StringVar(&contigfile,"file","","Contig file in fasta format.")

	var out string
	flag.StringVar(&out,"o","./gomrlbpout.parquet","Output file.")

	flag.Parse()

	runtime.GOMAXPROCS(proc)

	ls := SeqOps.ReadFasta(contigfile)
	lsp := SeqOps.SequenceCollection{}

	fmt.Println("Processing sequences")
	start := time.Now()
	var wg sync.WaitGroup
	in := make(chan SeqOps.Sequence, len(ls))

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go SeqOps.OHEncode(&wg, in, &lsp)
	}

	for _, s := range ls {
		in <- s
	}

	close(in)

	wg.Wait()


	t := time.Since(start)
	fmt.Printf("%v sequences analysed in %s \n",len(ls),t)

	lsp.ToParquet(out)
}