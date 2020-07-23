package SeqOps

import(
	"bufio"
	"compress/gzip"
	"os"
	"sync"
	"strings"
	"fmt"
)

type Sequence struct {
	Header string //contig id from fasta file
	Seq string // contig sequence
	OH [][]uint8 // one-hot encoded sequence
}

type SequenceCollection struct {
	sync.RWMutex
	Items []Sequence
}

func (cs *SequenceCollection) Append(item Sequence){
	cs.Lock()
	defer cs.Unlock()
	cs.Items = append(cs.Items, item)
}

// Create a new empty sequence struct.
func NewSequence() *Sequence {
	s := Sequence{"", "",nil}
	return &s
}

func ReadFasta(fname string) []Sequence {
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	bReader := bufio.NewReader(file)
	testBytes, err := bReader.Peek(2)
	if err != nil {
		panic(err)
	}

	file.Close()

	if testBytes[0] == 31 && testBytes[1] == 139 {
		file, err := os.Open(fname)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		gzipReader, err := gzip.NewReader(file)
		if err != nil{
			panic(err)
		}

		defer gzipReader.Close()
		scanner := bufio.NewScanner(gzipReader)
		sequenceList := make([]Sequence, 0)
		seq := new(Sequence)
		for scanner.Scan() {
			lstring := string(scanner.Text())
			if strings.Contains(lstring, ">") {
				if lstring == "" {
					fmt.Print("Empty line encountered!")
					break
				}
				seq.Header = lstring
			} else {
				if lstring == "" {
					fmt.Print("Empty line encountered!")
					break
				}
				seq.Seq = lstring
				sequenceList = append(sequenceList, *seq)
			}

		}

		if err := scanner.Err(); err != nil{
			panic(err)
		}
		return sequenceList


	}else {
		file, err := os.Open(fname)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		scanner := bufio.NewScanner(reader)

		sequenceList := make([]Sequence, 0)
		seq := new(Sequence)
		for scanner.Scan() {

			lstring := string(scanner.Text())
			if strings.Contains(lstring, ">") {
				if lstring == "" {
					fmt.Print("Empty line encountered!")
					break
				}
				seq.Header = lstring
			} else {
				if lstring == "" {
					fmt.Print("Empty line encountered!")
					break
				}
				seq.Seq = lstring
				sequenceList = append(sequenceList, *seq)
			}

		}
		if err := scanner.Err(); err != nil{
			panic(err)
		}
		return sequenceList
	}
	return nil
}