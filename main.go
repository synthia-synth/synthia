package main

//go:generate -command yacc go tool yacc
//go:generate yacc -o lang.go -p "lang" lang.y
