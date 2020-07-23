package SeqOps

import(
	"github.com/xitongsys/parquet-go/parquet"
	"log"
	"sync"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
)

func OHEncode(wg *sync.WaitGroup, in chan Sequence,lsp *SequenceCollection) {
	for s := range in {
		// Go's clunky way of making 2d slices.
		oh := make([][]uint8, 4)
		for i := range oh {
			oh[i] = make([]uint8, len(s.Seq))
		}
		for i, rne := range s.Seq {

			switch rne {
			case 'A':
				oh[0][i] = 1
			case 'C':
				oh[1][i] = 1
			case 'T':
				oh[2][i] = 1
			case 'G':
				oh[3][i] = 1
			}
		}

			rs := NewSequence()
			rs.Header = s.Header
			rs.OH = oh
			lsp.Append(*rs)
	}
	wg.Done()
}


func (sc *SequenceCollection) ToParquet(fname string){
	type tmpseq struct {
		Header    string  `parquet:"name=name, type=UTF8, encoding=PLAIN_DICTIONARY"`
		OH     [][]uint8   `parquet:"name=OHE, type=INT8, repetitiontype=REPEATED"`
	}

	fw, err := local.NewLocalFileWriter(fname)
	if err != nil {
		log.Println("Can't open file", err)
		return
	}
	pw, err := writer.NewParquetWriter(fw, new(tmpseq),4)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}
	pw.RowGroupSize = 5 * 1024 * 1024 //5M
	pw.CompressionType = parquet.CompressionCodec_GZIP

	for _,s := range sc.Items {
		seq := tmpseq{
			Header: s.Header,
			OH:  s.OH,
		}

		if err = pw.Write(seq); err != nil {
			log.Println("Write error", err)
		}
	}

	if err = pw.WriteStop(); err != nil {
		log.Println("WriteStop error", err)
		return
	}
	log.Println("Write Finished")
	fw.Close()
}