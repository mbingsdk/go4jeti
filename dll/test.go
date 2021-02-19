package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"time"
)

var data = map[string][]string{}

func pend(data map[string][]string){
	time.Sleep(10 * time.Second)
	fmt.Println(data)
}

func main(){
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("[------------MBING SHELL--------]")

	for{
		fmt.Print("$ ")
		text,_ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		s := strings.Split(text, ":")
		ss := strings.Split(s[1], ",")
		for i := range ss{
			data[s[0]] = append(data[s[0]], ss[i])
		}
		go pend(data)
	}
}