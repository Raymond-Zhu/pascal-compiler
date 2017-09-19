package main

import (
    "io/ioutil"
    "strconv"
    "strings"
    //"fmt"
)

type ScannerState int

const (
    COMMENT = iota
    STRING
    NUMERIC
    REAL
    SKIP
    NULL
)

type Token struct {
    Value string
    Type TokenType
    Col int
    Row int
}

func (t *Token) Print() string {
    return t.Value + "\t" + string(t.Type) + "\t" + strconv.Itoa(t.Col) + "\t" + strconv.Itoa(t.Row)
}

type Scanner struct {
    Col int
    Row int
    SavedToken TokenType
    Value string
    State ScannerState
    Tokens []*Token
}

func (s *Scanner) Scan(file string){
    b, _ := ioutil.ReadFile(file)
    source := string(b)

    for i, r := range source {
        str := string(r)
        //fmt.Println("String: " + str)
        //fmt.Println("Value: " + s.Value)
        //fmt.Println("Type: " + s.SavedToken)
        //fmt.Println(r)
        if s.State == SKIP {
            s.State = NULL
            s.Advance(str)
            continue
        } else if s.State == COMMENT {
            s.HandleComment(str, string(source[i+1]))
        } else if s.State == STRING {
            s.HandleString(str)
        } else if s.State == NUMERIC {
            s.HandleNumeric(str)
        } else {
            s.BuildOtherTokens(str)
        }
        s.Advance(str)
    }
}

func (s *Scanner) Advance(c string) {
    if c == "\n" || c == "\r"{
        s.Col = 1
        s.Row += 1
    }
    s.Col += 1
}


func (s *Scanner) HandleComment(current string, next string) {
    if current == "*" && next == ")" {
        s.State = SKIP
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"*)", Type: TK_END_COMMENT})
        s.SavedToken = TK_EMPTY
    }
}

func (s *Scanner) HandleString(c string) {
    if c == "'" {
        s.Value += c
        s.State = NULL
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:s.Value, Type: TK_STRING})
        s.Value = ""
    } else {
        s.Value += c
    }
}

func (s *Scanner) HandleNumeric(c string) {
    if _, err := strconv.Atoi(c); err == nil {
        s.Value += c
    } else if c == "." {
        s.Value += c
        s.State = REAL
    } else {
        if s.State == REAL {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:s.Value, Type: TK_REAL})
            s.Value = ""
        } else {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:s.Value, Type: TK_INT})
            s.Value = ""
        }
        s.State = NULL
    }

    if c == ";" {
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:";", Type: TK_SEMICOLON})
        s.SavedToken = TK_EMPTY
        s.Value = ""
    }
}

func IsNumeric(c string) bool {
    if _, err := strconv.Atoi(c); err == nil {
        return true
    }
    return false
}

func IsWhiteSpace(c string) bool {
    if c == "\n" || c == "\t" || c == " " || c == "\r"{
        return true
    }
    return false
}

func (s *Scanner) BuildOtherTokens(c string) {
    if IsWhiteSpace(c) {
        if s.SavedToken != TK_EMPTY {
            if Keywords[strings.ToUpper(s.Value)] != TK_EMPTY {
                s.SavedToken = Keywords[strings.ToUpper(s.Value)]
                s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:strings.ToLower(s.Value), Type:s.SavedToken })
                s.SavedToken = TK_EMPTY
                s.Value = ""
                return
            } else if Symbols[s.Value] != TK_EMPTY {
                s.SavedToken = Symbols[strings.ToUpper(s.Value)]
                s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:strings.ToLower(s.Value), Type:s.SavedToken })
                s.SavedToken = TK_EMPTY
                s.Value = ""
                return
            } else {
                if s.SavedToken == TK_OPEN_PAREN {
                    s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"(", Type:TK_OPEN_PAREN})
                    s.SavedToken = TK_EMPTY
                    s.Value = ""
                } else if s.SavedToken == TK_COLON {
                    s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:":", Type:TK_COLON})
                    s.SavedToken = TK_EMPTY
                    s.Value = ""
                } else {
                    s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:strings.ToLower(s.Value), Type:s.SavedToken})
                    s.SavedToken = TK_EMPTY
                    s.Value = ""
                }
                return
            }
        } else {
            return
        }
    }

    if c == "<" {
        if s.SavedToken == TK_EMPTY {
            s.SavedToken = TK_LESS
            s.Value += c
            return
        }
    }

    if c == ">" {
        if s.SavedToken == TK_EMPTY {
            s.SavedToken = TK_GREATER
            s.Value += c
            return
        }
    }

    if c == "(" {
        if s.SavedToken == TK_EMPTY {
            s.SavedToken = TK_OPEN_PAREN
            return
        } else {
            if token := Keywords[strings.ToUpper(s.Value)]; token != "" {
                s.Tokens = append(s.Tokens, &Token{Col:s.Col-1, Row:s.Row, Value: s.Value, Type: token})
            } else {
                s.Tokens = append(s.Tokens, &Token{Col:s.Col-1, Row:s.Row, Value: s.Value, Type: s.SavedToken})
            }
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"(", Type: TK_OPEN_PAREN})
            s.SavedToken = TK_EMPTY
            s.Value = ""
            return
        }
    }

    if c == "*" {
        if s.SavedToken == TK_OPEN_PAREN {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"(*", Type: TK_BEGIN_COMMENT})
            s.SavedToken = TK_EMPTY
            s.State = COMMENT
            return
        } else {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"*", Type: TK_MULT})
            return
        }
    }

    if c == ")" {
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:")", Type: TK_CLOSE_PAREN})
        s.SavedToken = TK_EMPTY
        s.Value = ""
    }

    if c =="[" {
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"[", Type: TK_OB})
        s.SavedToken = TK_EMPTY
        s.Value = ""
    }

    if c =="]" {
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"]", Type: TK_CB})
        s.SavedToken = TK_EMPTY
        s.Value = ""
    }

    if c == "," {
        if s.SavedToken == TK_IDENTIFIER {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:s.Value, Type: s.SavedToken})
            s.SavedToken = TK_EMPTY
            s.Value = ""
        }
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:",", Type: TK_COMMA})
        return
    }

    if c == "!" {
        if s.SavedToken == TK_EMPTY {
            s.SavedToken = TK_EXCL
            s.Value += c
            return
        }
    }

    if c == ":" {
        if s.SavedToken == TK_EMPTY {
            s.SavedToken = TK_COLON
            s.Value += c
            return
        }
    }

    if c == "=" {
        if s.SavedToken == TK_EMPTY {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"=", Type: TK_EQUALS})
        } else if s.SavedToken == TK_COLON {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:":=", Type: TK_ASSIGNMENT})
        } else if s.SavedToken == TK_GREATER {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:">=", Type:TK_GE})
        } else if s.SavedToken == TK_LESS {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"<=", Type:TK_LE})
        }
        s.SavedToken = TK_EMPTY
        s.Value = ""
        return
    }

    if c == ";" && s.State != NUMERIC {
        if s.SavedToken != TK_EMPTY {
            if Keywords[strings.ToUpper(s.Value)] != TK_EMPTY {
                s.SavedToken = Keywords[strings.ToUpper(s.Value)]
            }
            s.Tokens = append(s.Tokens, &Token{Col:s.Col-1, Row:s.Row, Value:strings.ToLower(s.Value), Type:s.SavedToken})
        }

        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:";", Type: TK_SEMICOLON})
        s.SavedToken = TK_EMPTY
        s.Value = ""
        return
    }

    if IsNumeric(c) {
        s.State = NUMERIC
        s.Value += c
        return
    }

    if c == "." {
        if s.Value == "end" {
            s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"end.", Type: TK_END_CODE})
            s.SavedToken = TK_EMPTY
            s.Value =""
            return
        }
    }

    if c == "'" {
        s.State = STRING
        s.Value += c
        return
    }

    if c == "-" {
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"-", Type: TK_MINUS})
        s.SavedToken = TK_EMPTY
        s.Value = ""
        return
    }

    if c == "/" {
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"/", Type: TK_DIV_FLOAT})
        s.SavedToken = TK_EMPTY
        s.Value = ""
        return
    }

    if c == "+" {
        s.Tokens = append(s.Tokens, &Token{Col:s.Col, Row:s.Row, Value:"+", Type: TK_PLUS})
        s.SavedToken = TK_EMPTY
        s.Value = ""
        return
    }

    s.Value += c
    if Keywords[strings.ToUpper(s.Value)] == "" && Symbols[strings.ToUpper(s.Value)] == "" {
        s.SavedToken = TK_IDENTIFIER
        return
    }
}
