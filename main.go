package main

import (
    "os"
    "fmt"
    "text/tabwriter"
)
func main () {
    var tokens []*Token
    s := Scanner{Col: 1, Row: 1, SavedToken: TK_EMPTY, Value: "", State: NULL, Tokens: tokens}
    s.Scan(os.Args[1])

    fmt.Println("TOKENS:")
    w := tabwriter.NewWriter(os.Stdout,0,0,1,' ', tabwriter.Debug)
    fmt.Fprintln(w, "Name\tValue\tCol\tRow")
    for _, token := range s.Tokens {
        fmt.Fprintln(w, token.Print())
    }
    w.Flush()

    fmt.Print("\n\n")
    fmt.Println("INSTRUCTIONS:")

    var symTable []*SymTableEntry
    p := Parser{Tokens:s.Tokens, SymTable:symTable, CurrentToken: s.Tokens[0], Index: 1, Address:0}
    p.Parse()

    fmt.Fprintln(w, "Pointer\tOpCode\tValue")
    for _, v := range p.Instructions {
        fmt.Fprintf(w, "IP: %d\t%v\t%v\n", v.Ip,  v.Name, v.Value )
    }
    w.Flush()

    fmt.Print("\n\n")
    fmt.Println("SYMBOL TABLE:")
    fmt.Fprintln(w, "Identifier\tAddress\tValue")
    for  _, v  :=  range(p.SymTable) {
        fmt.Fprintf(w, "%s\t0x%d\t%v\n", v.Name, v.Address, v.Value)
    }
    w.Flush()

    fmt.Print("\n\n")
    fmt.Println("OUTPUT:")
    var stack []interface{}
    sm := StackMachine{Stack: stack, Ip: 0, SymTable: p.SymTable, Instructions: p.Instructions }
    sm.Run()
    fmt.Print("\n\n")

    fmt.Println("SYMBOL TABLE:")
    fmt.Fprintln(w, "Identifier\tAddress\tValue")
    for  _, v  :=  range(sm.SymTable) {
        fmt.Fprintf(w, "%s\t0x%d\t%v\n", v.Name, v.Address, v.Value)
    }
    w.Flush()
}
