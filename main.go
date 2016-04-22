package main

import (
	"os"
)

func main() {
	ps := PostSlack{}
	ps.Run(os.Args)
}
