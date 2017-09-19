package main

import (
    "strconv"
)

type OpCode string

const (
    OP_PUSH = "op_push"
    OP_POP = "op_pop"
    OP_OR = "op_or"
    OP_AND = "op_and"
    OP_XOR = "op_xor"
    OP_NOT = "op_not"
    OP_MOD = "op_mod"
    OP_LESS = "op_less"
    OP_GREATER = "op_greater"
    OP_LE = "op_less_equals"
    OP_GE =  "op_greater_equals"
    OP_EQUALS = "op_equals"
    OP_NOT_EQUALS = "op_not_equals"
    OP_DIV_FLOAT = "op_div_float"
    OP_ADD = "op_add"
    OP_MINUS = "op_minus"
    OP_MULT = "op_mult"
    OP_END = "op_end"
    OP_JFALSE = "op_jfalse"
    OP_JMP = "op_jmp"
    OP_JTRUE = "op_jtrue"
    OP_WRITELN = "op_writeln"
    OP_HALT = "op_halt"
)

type Parser struct {
    Tokens []*Token
    CurrentToken *Token
    Index int
    StoredValue interface{}
    SymTable []*SymTableEntry
    Instructions []*Instruction
    Address int
    Ip int
    Op bool
}

type Instruction struct {
    Name OpCode
    Ip int
    Value interface{}
}

type SymType string

const (
    INT = "INTEGER"
    REALTYPE = "REAL"
    NIL = "NIL"
    BOOL = "BOOL"
    CHAR = "CHAR"
)

type SymTableEntry struct {
    Name string
    Address int
    Type SymType
    Value interface{}
}

func (p *Parser) Parse() {
    p.Program()
}

func (p *Parser) Program() {
    p.MustMatchAdvance(TK_PROGRAM)
    p.MustMatchAdvance(TK_IDENTIFIER)
    p.MustMatchAdvance(TK_SEMICOLON)
    p.DeclareVars()
}

func (p *Parser) DeclareVars() {
    if !p.Match(TK_VAR) {
        p.Begin()
        return
    } else{
        p.Advance()
    }

    for {
        if p.Match(TK_IDENTIFIER) {
            e := &SymTableEntry{
                Name: p.CurrentToken.Value,
                Address: p.Address,
                Type: NIL,
            }
            p.SymTable = append(p.SymTable, e)
            p.Address += 4
            p.Advance()
        } else if p.MatchAdvance(TK_COMMA) {
            continue
        } else {
            p.MustMatchAdvance(TK_COLON)
            break
        }
    }

    if p.MatchAdvance(TK_ID_INT) {
        for _, e := range p.SymTable {
            if e.Type == NIL {
                e.Type = INT
            }
        }
    } else if p.MatchAdvance(TK_ID_REAL) {
        for _, e := range p.SymTable {
            if e.Type == NIL {
               e.Type = REALTYPE
            }
        }
    } else if p.MatchAdvance(TK_ID_BOOLEAN) {
        for _, e := range p.SymTable {
            if e.Type == NIL {
               e.Type = BOOL
            }
        }
    } else if p.MatchAdvance(TK_ID_CHAR) {
        for _, e := range p.SymTable {
            if e.Type == NIL {
               e.Type = CHAR
            }
        }
    }
    p.MustMatchAdvance(TK_SEMICOLON)

    p.DeclareVars()
}

func (p *Parser) Begin() {
    p.MustMatchAdvance(TK_BEGIN)
    for !p.Match(TK_END_CODE) {
        p.ExecuteStatements()
    }
    p.AddInstruction(OP_HALT, "HALT", p.Ip)
}

func (p *Parser) ExecuteStatements() {
    for {
        switch p.CurrentToken.Type {
        case TK_IF:
            p.IfStatement()
        case TK_FOR:
            p.ForStatement()
        case TK_REPEAT:
            p.RepeatStatement()
        case TK_WHILE:
            p.WhileLoop()
        case TK_WRITELN:
            p.WriteLine()
        case TK_IDENTIFIER:
            p.StoredValue = p.CurrentToken.Value
            p.MustMatchAdvance(TK_IDENTIFIER)
        }

        if p.Match(TK_OB) {
            p.Advance()
            p.Factor()
            p.MustMatchAdvance(TK_CB)
        }

        if p.Match(TK_ASSIGNMENT) {
            p.MustMatchAdvance(TK_ASSIGNMENT)
            p.Op = true
        }

        p.DoLogic()
        if p.Match(TK_SEMICOLON) {
            p.Advance()
            if p.Op {
                p.AddInstruction(OP_POP, p.StoredValue ,p.Ip)
                p.Ip += 1
                p.Op = false
            }
        }

        if p.Match(TK_END_CODE) {
            return
        }

        if p.Match(TK_UNTIL) {
            return
        }

        if p.Match(TK_TO) {
            return
        }

        if p.Match(TK_ELSE) {
            return
        }
    }
}

func(p *Parser) IfStatement() {
    p.MustMatchAdvance(TK_IF)
    p.DoLogic()
    p.MustMatchAdvance(TK_THEN)
    p.Instructions = append(p.Instructions, &Instruction{Name:OP_JFALSE, Ip:p.Ip, Value: 0})
    ip := p.Ip
    p.Ip += 1
    p.ExecuteStatements()

    if p.Match(TK_ELSE) {
        ip2 := p.Ip
        p.Instructions = append(p.Instructions, &Instruction{Name:OP_JMP, Ip:p.Ip, Value: NIL})
        p.Ip += 1
        p.Advance()
        p.UpdateInstruction(ip)
        p.ExecuteStatements()
        p.UpdateInstruction(ip2)
    }
}

func (p *Parser) ForStatement(){
    p.MustMatchAdvance(TK_FOR)
    p.ExecuteStatements()
    loopVar := p.SymTable[0].Name
    p.MustMatchAdvance(TK_TO)
    loopStartIp := p.Ip
    p.AddInstruction(OP_PUSH,loopVar, p.Ip)
    p.Ip += 1
    p.Factor()
    p.MustMatchAdvance(TK_DO)
    p.AddInstruction(OP_GREATER, ">", p.Ip)
    p.Ip += 1
    ip := p.Ip
    p.AddInstruction(OP_JTRUE, ip, p.Ip )
    p.Ip += 1
    p.ExecuteStatements()

    p.AddInstruction(OP_PUSH, loopVar, p.Ip)
    p.Ip += 1
    p.AddInstruction(OP_PUSH, 1, p.Ip)
    p.Ip += 1
    p.AddInstruction(OP_ADD, "+", p.Ip)
    p.Ip += 1
    p.AddInstruction(OP_POP,  loopVar, p.Ip)
    p.Ip += 1
    p.AddInstruction(OP_JMP, loopStartIp, p.Ip)
    p.Ip += 1
    p.UpdateInstruction(ip)
}

func (p *Parser) WriteLine() {
    p.MustMatchAdvance(TK_WRITELN)
    p.MustMatchAdvance(TK_OPEN_PAREN)
    p.DoLogic()
    p.MustMatchAdvance(TK_CLOSE_PAREN)
    p.MustMatchAdvance(TK_SEMICOLON)
    p.AddInstruction(OP_WRITELN, "", p.Ip)
    p.Ip += 1
}

func (p *Parser) UpdateInstruction(ip int) {
    p.Instructions[ip].Value = p.Ip
}

func (p *Parser) DoLogic() {
    switch p.CurrentToken.Type {
    case TK_LESS:
        p.MustMatchAdvance(TK_LESS)
        p.Expr()
        p.BuildInstruction(TK_LESS)
    case TK_GREATER:
        p.MustMatchAdvance(TK_GREATER)
        p.Expr()
        p.BuildInstruction(TK_GREATER)
    case TK_LE:
        p.MustMatchAdvance(TK_LE)
        p.Expr()
        p.BuildInstruction(TK_LE)
    case TK_GE:
        p.MustMatchAdvance(TK_GE)
        p.Expr()
        p.BuildInstruction(TK_GE)
    case TK_EQUALS:
        p.MustMatchAdvance(TK_EQUALS)
        p.Expr()
        p.BuildInstruction(TK_EQUALS)
    case TK_NOT_EQUALS:
        p.MustMatchAdvance(TK_NOT_EQUALS)
        p.Expr()
        p.BuildInstruction(TK_NOT_EQUALS)
    default:
        p.Expr()
    }
}

func (p *Parser) RepeatStatement() {
    p.MustMatchAdvance(TK_REPEAT)
    ip := p.Ip
    p.ExecuteStatements()
    p.MustMatchAdvance(TK_UNTIL)
    p.DoLogic()
    p.AddInstruction(OP_JFALSE, ip, p.Ip)
    p.Ip += 1
}

func (p *Parser) WhileLoop() {
    p.MustMatchAdvance(TK_WHILE)
    ip := p.Ip
    p.DoLogic()
    p.MustMatchAdvance(TK_DO)
    ip2 := p.Ip
    p.AddInstruction(OP_JFALSE, p.Ip, p.Ip)
    p.Ip += 1
    p.ExecuteStatements()
    p.AddInstruction(OP_JMP, ip, p.Ip)
    p.Ip += 1
    p.UpdateInstruction(ip2)
    return
}

func (p *Parser) Expr() {
    p.Term()
    p.ExprRecursive()
}

func (p *Parser) ExprRecursive() {
    switch p.CurrentToken.Type {
    case TK_PLUS:
        p.MustMatchAdvance(TK_PLUS)
        p.Term()
        p.BuildInstruction(TK_PLUS)
        p.ExprRecursive()
    case TK_MINUS:
        p.MustMatchAdvance(TK_MINUS)
        p.Term()
        p.BuildInstruction(TK_MINUS)
        p.ExprRecursive()
    case TK_OR:
        p.MustMatchAdvance(TK_OR)
        p.Term()
        p.BuildInstruction(TK_OR)
        p.ExprRecursive()
    case TK_XOR:
        p.MustMatchAdvance(TK_XOR)
        p.Term()
        p.BuildInstruction(TK_XOR)
        p.ExprRecursive()
    default:
        return
    }
}

func (p *Parser) Term() {
    p.Factor()
    p.TermRecursive()
}

func (p *Parser) TermRecursive() {
    switch p.CurrentToken.Type {
    case TK_LESS:
        p.MustMatchAdvance(TK_LESS)
        p.Expr()
        p.BuildInstruction(TK_LESS)
    case TK_GREATER:
        p.MustMatchAdvance(TK_GREATER)
        p.Expr()
        p.BuildInstruction(TK_GREATER)
    case TK_LE:
        p.MustMatchAdvance(TK_LE)
        p.Expr()
        p.BuildInstruction(TK_LE)
    case TK_GE:
        p.MustMatchAdvance(TK_GE)
        p.Expr()
        p.BuildInstruction(TK_GE)
    case TK_EQUALS:
        p.MustMatchAdvance(TK_EQUALS)
        p.Expr()
        p.BuildInstruction(TK_EQUALS)
    case TK_NOT_EQUALS:
        p.MustMatchAdvance(TK_NOT_EQUALS)
        p.Expr()
        p.BuildInstruction(TK_NOT_EQUALS)
    case TK_MULT:
        p.MustMatchAdvance(TK_MULT)
        p.Factor()
        p.BuildInstruction(TK_MULT)
        p.TermRecursive()
    case TK_DIV_FLOAT:
        p.MustMatchAdvance(TK_DIV_FLOAT)
        p.Factor()
        p.BuildInstruction(TK_DIV_FLOAT)
        p.TermRecursive()
    case TK_MOD:
        p.MustMatchAdvance(TK_MOD)
        p.Factor()
        p.BuildInstruction(TK_MOD)
        p.TermRecursive()
    case TK_AND:
        p.MustMatchAdvance(TK_AND)
        p.Factor()
        p.BuildInstruction(TK_AND)
        p.TermRecursive()
    default:
        return
    }
}

func (p *Parser) Factor() {
    switch p.CurrentToken.Type {
    case TK_IDENTIFIER:
        p.BuildInstruction(TK_IDENTIFIER)
        p.Advance()
        return
    case TK_STRING:
        p.BuildInstruction(TK_STRING)
        p.Advance()
        return
    case TK_INT:
        p.BuildInstruction(TK_INT)
        p.Advance()
        return
    case TK_NOT:
        p.MustMatchAdvance(TK_NOT)
        p.Factor()
        p.BuildInstruction(TK_NOT)
        return
    case TK_OPEN_PAREN:
        p.MustMatchAdvance(TK_OPEN_PAREN)
        p.DoLogic()
        p.MustMatchAdvance(TK_CLOSE_PAREN)
        return
    }
}

func (p *Parser) AddInstruction(op OpCode, value interface{}, ip int) {
    i := &Instruction{
        Name: op,
        Value: value,
        Ip: ip,
    }
    p.Instructions = append(p.Instructions, i)
}

func (p *Parser) BuildInstruction(t TokenType) {
    switch t {
    case TK_IDENTIFIER:
        p.AddInstruction(OP_PUSH, p.CurrentToken.Value, p.Ip)
    case TK_STRING:
        p.AddInstruction(OP_PUSH, p.CurrentToken.Value, p.Ip)
    case TK_INT:
        i, _ := strconv.Atoi(p.CurrentToken.Value)
        p.AddInstruction(OP_PUSH, i, p.Ip)
    case TK_PLUS:
        p.AddInstruction(OP_ADD, "+", p.Ip)
    case TK_MINUS:
        p.AddInstruction(OP_MINUS, "+", p.Ip)
    case TK_MULT:
        p.AddInstruction(OP_MULT, "*", p.Ip)
    case TK_DIV_FLOAT:
        p.AddInstruction(OP_DIV_FLOAT, "/", p.Ip)
    case TK_MOD:
        p.AddInstruction(OP_MOD, "mod", p.Ip)
    case TK_OR:
        p.AddInstruction(OP_OR, "or", p.Ip)
    case TK_AND:
        p.AddInstruction(OP_AND, "and", p.Ip)
    case TK_NOT:
        p.AddInstruction(OP_NOT, "not", p.Ip)
    case TK_XOR:
        p.AddInstruction(OP_XOR, "xor", p.Ip)
    case TK_GREATER:
        p.AddInstruction(OP_GREATER, ">", p.Ip)
    case TK_LESS:
        p.AddInstruction(OP_LESS, "<", p.Ip)
    case TK_GE:
        p.AddInstruction(OP_GE, ">=", p.Ip)
    case TK_LE:
        p.AddInstruction(OP_LE, "<=", p.Ip)
    case TK_EQUALS:
        p.AddInstruction(OP_EQUALS, "=", p.Ip)
    case TK_NOT_EQUALS:
        p.AddInstruction(OP_NOT_EQUALS, "=", p.Ip)
    default:
        return
    }

    p.Ip += 1
}

func (p *Parser) Advance() {
    p.CurrentToken = p.Tokens[p.Index]
    p.Index += 1
}

func (p *Parser) MatchAdvance(t TokenType) bool {
    b := p.Match(t)
    if b {
        p.Advance()
    }
    return b
}

func (p *Parser) Match(t TokenType) bool {
    if p.CurrentToken.Type == t {
        return true
    } else {
        return false
    }
}

func (p *Parser) MustMatch(t TokenType) {
    if p.CurrentToken.Type != t {
        p.Error(t)
    }
}

func (p *Parser) MustMatchAdvance(t TokenType) {
    p.MustMatch(t)
    p.Advance()
}

func (p *Parser) Error(t TokenType) {
    panic("ERROR, expected token: " + string(t) + ". Got " + string(p.CurrentToken.Type))
}
