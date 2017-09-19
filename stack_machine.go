package main
import(
    "fmt"
    "strings"
)


func (s *StackMachine) Push(v interface{}) {
    s.Stack = append(s.Stack, v)
}

func (s *StackMachine) Pop() interface{} {
    l := len(s.Stack)
    val := s.Stack[l-1]
    s.Stack = s.Stack[:l-1]
    return val
}

type StackMachine struct {
    Stack []interface{}
    Ip int
    SymTable []*SymTableEntry
    Instructions []*Instruction
}

func (s *StackMachine) Run() {
    for {
        //fmt.Println(s.Instructions[s.Ip].Name)
        //fmt.Printf("%+v\n", s.Stack)
        //for  _, v  :=  range(s.SymTable) {
            //fmt.Printf("%+v ", v)
        //}
        //fmt.Printf("\n")
        switch s.Instructions[s.Ip].Name {
        case OP_PUSH:
            value, ok := s.Instructions[s.Ip].Value.(string)
            if ok {
                //fmt.Println("Value: " + value)
                if !strings.Contains(value,"'") {
                    s.OpPushIdentifier(value)
                } else {
                    s.OpPush(s.Instructions[s.Ip].Value)
                }
            } else {
                s.OpPush(s.Instructions[s.Ip].Value)
            }
        case OP_POP:
            s.OpPop(s.Instructions[s.Ip].Value)
        case OP_HALT:
            fmt.Println("Program End.")
            return
        case OP_OR:
            s.OpOr()
        case OP_AND:
            s.OpAnd()
        case OP_XOR:
            s.OpXor()
        case OP_NOT:
            s.OpNot()
        case OP_MOD:
            s.OpMod()
        case OP_ADD:
            s.OpAdd()
        case OP_MINUS:
            s.OpSub()
        case OP_DIV_FLOAT:
            s.OpDiv()
        case OP_MULT:
            s.OpMult()
        case OP_GREATER:
            s.OpGreater()
        case OP_LESS:
            s.OpLess()
        case OP_GE:
            s.OpGreaterEq()
        case OP_LE:
            s.OpLessEq()
        case OP_EQUALS:
            s.OpEq()
        case OP_NOT_EQUALS:
            s.OpNotEq()
        case OP_JFALSE:
            s.OpJFalse(s.Instructions[s.Ip].Value.(int))
        case OP_JTRUE:
            s.OpJTrue(s.Instructions[s.Ip].Value.(int))
        case OP_JMP:
            s.OpJump(s.Instructions[s.Ip].Value.(int))
        case OP_WRITELN:
            s.WriteLine()
        default:
            fmt.Println("Error instruction doesn't exist")
        }
        s.Ip += 1
    }
}

func (s *StackMachine) WriteLine() {
    value := s.Pop()
    fmt.Println(value)
}

func (s *StackMachine) OpJTrue(ip int) {
    value := s.Pop()
    if value == true {
        s.Ip = ip - 1
    }
}

func (s *StackMachine) OpJFalse(ip int) {
    value := s.Pop()
    if value == false {
        s.Ip = ip - 1
    }
}

func (s *StackMachine) OpJump(ip int) {
    s.Ip = ip - 1
}

func (s *StackMachine) OpLessEq() {
    value := s.Pop()
    value2 := s.Pop()
    if value2.(int) <= value.(int) {
        s.Push(true)
    } else {
        s.Push(false)
    }
}

func (s *StackMachine) OpGreaterEq() {
    value := s.Pop()
    value2 := s.Pop()
    if value2.(int) >= value.(int) {
        s.Push(true)
    } else {
        s.Push(false)
    }
}

func (s *StackMachine) OpNotEq() {
    value := s.Pop()
    value2 := s.Pop()
    if value2.(int) != value.(int) {
        s.Push(true)
    } else {
        s.Push(false)
    }
}
func (s *StackMachine) OpEq() {
    value := s.Pop()
    value2 := s.Pop()
    if value2.(int) == value.(int) {
        s.Push(true)
    } else {
        s.Push(false)
    }
}

func (s *StackMachine) OpGreater() {
    value := s.Pop()
    value2 := s.Pop()
    if value2.(int) > value.(int) {
        s.Push(true)
    } else {
        s.Push(false)
    }
}

func (s *StackMachine) OpLess() {
    value := s.Pop()
    value2 := s.Pop()
    if value2.(int) < value.(int) {
        s.Push(true)
    } else {
        s.Push(false)
    }
}

func (s *StackMachine) OpNot() {
    value := s.Pop()
    if !value.(bool) == true {
        s.Push(true)
    } else {
        s.Push(false)
    }
}

func (s *StackMachine) OpAnd() {
    value := s.Pop()
    value2 := s.Pop()
    if value.(bool) && value2.(bool) {
        s.Push(true)
    } else {
        s.Push(false)
    }
}

func (s *StackMachine) OpOr() {
    value := s.Pop()
    value2 := s.Pop()
    if value.(bool) || value2.(bool) {
        s.Push(true)
    } else {
        s.Push(false)
    }
}

func (s *StackMachine) OpXor() {
    value := s.Pop()
    value2 := s.Pop()
    if value.(bool) != value2.(bool) {
        s.Push(true)
    } else {
        s.Push(false)
    }
}

func (s *StackMachine) OpMod() {
    value := s.Pop()
    value2 := s.Pop()
    s.Push(value.(int) % value2.(int))
}

func (s *StackMachine) OpDiv() {
    value := s.Pop()
    value2 := s.Pop()
    s.Push(value2.(int) / value.(int))
}

func (s *StackMachine) OpMult() {
    value := s.Pop()
    value2 := s.Pop()
    s.Push(value2.(int) * value.(int))
}

func (s *StackMachine) OpAdd() {
    value := s.Pop()
    value2 := s.Pop()
    s.Push(value2.(int) + value.(int))
}

func (s *StackMachine) OpSub() {
    value := s.Pop()
    value2 := s.Pop()
    s.Push(value2.(int) - value.(int))
}

func (s *StackMachine) OpPush(value interface{}) {
    s.Push(value)
}

func (s *StackMachine) OpPushIdentifier(value interface{}) {
    for _, entry := range s.SymTable {
        if entry.Name == value {
            s.Push(entry.Value)
        }
    }
}

func (s *StackMachine) OpPop(value interface{}) {
    v := s.Pop()
    for _, entry := range s.SymTable {
        if entry.Name == value {
            entry.Value = v
        }
    }
}
