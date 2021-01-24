package main

import (
	"fmt"
	"interpreter-go/repl"
	"os"
	"os/user"
)

func main(){
	user,err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programing lanuguage!\n",user.Username)
	fmt.Printf("Fell free to type in commnads\n")
	repl.Start(os.Stdin,os.Stdout)
}
