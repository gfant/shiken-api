## How to run it.

First of all, change in `executor/handler/constants.go` and `reader/handler/constants.go`  the variable `headPath` to the output of the command `pwd` while being in the folder where you cloned this repository.

Then in each folder (`executor` and `reader`) do

```bash
go run main.go
```

This will run the servers required.

### Reader

Will provide 
* The list of problems
* The information of the problem chosen

### Executor

Will
* Execute the code sent by the user