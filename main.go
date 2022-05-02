package main

import "github.com/hasindulanka/MethaneToolkitGo/methane"

func main() {

	methane.MakeDir(methane.WSRoot)
	methane.MakeDir(methane.WSCache)

	methane.Print("-------------------------------------")
	methane.Print("          Methane Toolkit            ")
	methane.Print("       github.com/HasinduLanka       ")
	methane.Print("-------------------------------------")
	methane.Print("")

}
