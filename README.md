## PoaGo

PoaGo (pōgō) is a Golang implementation of the Partial Order Alignment algorithm described by [Lee 2002](http://bioinformatics.oxfordjournals.org/content/18/3/452.short). This implementation was also inspired by a [blog post](http://simpsonlab.github.io/2015/05/01/understanding-poa/) by the SimpsonLab. 

### Install
```
go build
```

### Example
```
./PoaGo -f ./examples/example4.fa
```

Output is default to CLUSTAL format. Input can be fastQ or fastA. 

TODOs:
1. Profile
2. Concurrent DP
... I'm sure there's more
