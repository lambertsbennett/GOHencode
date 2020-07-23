# GOHencode
One-hot encoding of biological sequence data in Go.

## Premise
Many common recurrent neural network architectures can accept one-hot encoded data of variable length as input. GOHencode is a convenient tool to one-hot encode entire fasta files.

The input to GOHencode is a (gzipped) fasta file of sequences. These files have a structure as follows:

```
>Sequence-identifier
ATCGAA
```

GOHencode processes a fasta file and returns a parquet file with the following structure:
| Sequence ID | Encoded Sequence Array |
| :---        | :---                   |
| Sequence identifier 1 | [4*len(sequence)] |

The encoded sequence array has the following form for the sample sequence above:
```
A | 1 | 0 | 0 | 0 | 1 | 1  
C | 0 | 0 | 1 | 0 | 0 | 0  
T | 0 | 1 | 0 | 0 | 0 | 0  
G | 0 | 0 | 0 | 1 | 0 | 0  
```
Parquet has some limitations with storing 2D arrays, so the output ecoded array is of the form:

| Base1 | Base2 | ... | Basen|
| ---   | ---   | --- | ---  |
|[1 0 0 0]|[0 0 1 0] | ... | [...]|
