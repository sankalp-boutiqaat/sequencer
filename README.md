## Sequencer .. What It Does?
A utility to provide thread safe counters. 

- Counters can be Incrementing.  
```
1 2 3 4 5 6 7 ... ..10
```
- Counters can be Decrementing.
```
10 9 8 7 6 5 .... ..0
```  
- Counters can be Rolling.
```
1 2 3 4 5 1 2 3 4 ...
```  
- Counters can be Decrementing and Rolling.
```
5 4 3 2 1 5 4 3 2 ...
```

## How to Use?
Detailed illustrations are kept in ![examples](https://github.com/sanksons/sequencer/tree/master/examples) directory.


## Supported Adapters
As of Now, Sequencer supports following adapters.
- redis
- redis cluster

But, It can be easily extended to use any persistent data source.
