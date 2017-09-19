package main

type TokenType string

const (
    TK_BEGIN  = "TK_BEGIN"
    TK_BREAK  = "TK_BREAK"
    TK_CONST  = "TK_CONST"
    TK_DO =  "TK_DO"
    TK_DOT = "TK_DOT"
    TK_DOWNTO = "TK_DOWNTO"
    TK_ELSE = "TK_ELSE"
    TK_END = "TK_END"
    TK_END_CODE = "TK_END_CODE"
    TK_FOR = "TK_FOR"
    TK_FUNCTION = "TK_FUNCTIOn"
    TK_IDENTIFIER = "TK_IDENTIFIER"
    TK_IF = "TK_IF"
    TK_LABEL = "TK_LABEL"
    TK_PROGRAM = "TK_PROGRAM"
    TK_REPEAT = "TK_REPEAT"
    TK_STRING = "TK_STRING"
    TK_THEN = "TK_THEN"
    TK_TO = "TK_TO"
    TK_TYPE = "TK_TYPE"
    TK_VAR = "TK_VAR"
    TK_UNTIL = "TK_UNTIL"
    TK_WHILE = "TK_WHILE"
    TK_INT = "TK_INT"
    TK_REAL = "TK_REAL"
    TK_CHAR = "TK_CHAR"
    TK_BOOLEAN = "TK_BOOLEAN"
    TK_ID_INT = "TK_ID_INT"
    TK_ID_REAL = "TK_ID_REAL"
    TK_ID_CHAR = "TK_ID_CHAR"
    TK_ID_BOOLEAN = "TK_ID_BOOLEAN"
    TK_OF = "TK_OF"
    TK_PLUS = "TK_PLUS"
    TK_MINUS = "TK_MINUS"
    TK_MULT = "TK_MULT"
    TK_DIV_FLOAT = "TK_DIV_FLOAT"
    TK_DIV = "TK_DIV"
    TK_MOD = "TK_MOD"
    TK_COLON = "TK_COLON"
    TK_EQUALS = "TK_EQUALS"
    TK_ASSIGNMENT = "TK_ASSIGNMENT"
    TK_GREATER = "TK_GREATER"
    TK_LESS = "TK_LESS"
    TK_GE = "TK_GE"
    TK_LE = "TK_LE"
    TK_EXCL = "TK_EXCL"
    TK_NOT_EQUALS = "TK_NOT_EQUALS"
    TK_AND = "TK_AND"
    TK_XOR = "TK_XOR"
    TK_OR = "TK_OR"
    TK_NOT = "TK_NOT"
    TK_SEMICOLON = "TK_SEMICOLON"
    TK_OPEN_PAREN = "TK_OPEN_PAREN"
    TK_CLOSE_PAREN = "TK_CLOSE_PAREN"
    TK_QUOTE = "TK_QUOTE"
    TK_BEGIN_COMMENT = "TK_BEGIN_COMMENT"
    TK_END_COMMENT = "TK_END_COMMENT"
    TK_COMMA = "TK_COMMA"
    TK_RANGE = "TK_RANGE"
    TK_ARRAY = "TK_ARRAY"
    TK_OB = "TK_OB"
    TK_CB = "TK_CB"
    TK_WRITELN = "TK_WRITELN"
    TK_EMPTY = ""
)

var Keywords = map[string]TokenType {
        "BEGIN"     : TK_BEGIN,
        "BREAK"     : TK_BREAK,
        "BOOLEAN"   : TK_ID_BOOLEAN,
        "CHAR"      : TK_ID_CHAR,
        "CONST"     : TK_CONST,
        "DO"        : TK_DO,
        "DOWNTO"    : TK_DOWNTO,
        "ELSE"      : TK_ELSE,
        "END"       : TK_END,
        "END."      : TK_END_CODE,
        "FOR"       : TK_FOR,
        "INTEGER"   : TK_ID_INT,
        "FUNCTION"  : TK_FUNCTION,
        "IDENTIFIER": TK_IDENTIFIER,
        "IF"        : TK_IF,
        "LABEL"     : TK_LABEL,
        "PROGRAM"   : TK_PROGRAM,
        "OF"        : TK_OF,
        "REAL"      : TK_ID_REAL,
        "REPEAT"    : TK_REPEAT,
        "STRING"    : TK_STRING,
        "THEN"      : TK_THEN,
        "TO"        : TK_TO,
        "TYPE"      : TK_TYPE,
        "VAR"       : TK_VAR,
        "UNTIL"     : TK_UNTIL,
        "WHILE"     : TK_WHILE,
        "WRITELN"   : TK_WRITELN,
        "DIV"       : TK_DIV,
        "MOD"       : TK_MOD,
        "AND"       : TK_AND,
        "XOR"       : TK_XOR,
        "OR"        : TK_OR,
        "NOT"       : TK_NOT,
        "ARRAY"     : TK_ARRAY,
}

var Symbols = map[string]TokenType {
        "+"         : TK_PLUS,
        "-"         : TK_MINUS,
        "*"         : TK_MULT,
        "/"         : TK_DIV_FLOAT,
        ":"         : TK_COLON,
        "."         : TK_DOT,
        "="         : TK_EQUALS,
        ":="        : TK_ASSIGNMENT,
        ">"         : TK_GREATER,
        "<"         : TK_LESS,
        ">="        : TK_GE,
        "<="        : TK_LE,
        "!"         : TK_EXCL,
        "!="        : TK_NOT_EQUALS,
        ";"         : TK_SEMICOLON,
        "("         : TK_OPEN_PAREN,
        ")"         : TK_CLOSE_PAREN,
        "\""        : TK_QUOTE,
        "(*"        : TK_BEGIN_COMMENT,
        "*)"        : TK_END_COMMENT,
        ","         : TK_COMMA,
        ".."         : TK_RANGE,
        "["         : TK_OB,
        "]"         : TK_CB,
}