package lexer

import (
	"Compiler/c-monkey-v5/src/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	var input string = `let five = 5;
	let ten = 10;
	let add = fn(x, y) {
		x+y;
	};
	let result = add(five, ten);

	!-*/5;
	5 < 10 > 5;
	5 == 5;
	5 != 10;
	5 <= 10;
	10 >= 5;
	true;
	false;
	if;
	else;
	return;
	"foobar"
	"foo bar"
	[1, 3];
	{"foo": "bar"}
	`

	tests := []struct {
		expectedType token.TokenType
		expectedLit  string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.NOT_EQ, "!="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT_EQ, "<="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.GT_EQ, ">="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.SEMICOLON, ";"},
		{token.ELSE, "else"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.SEMICOLON, ";"},

		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},

		{token.EOF, ""},
	}

	stringUnderTest := New(input)

	for i, testToken := range tests {
		token := stringUnderTest.NextToken()

		if token.Type != testToken.expectedType {
			t.Fatalf("testing value [%d] - tokentype wrong. expected=%q, got=%q", i+1, testToken.expectedType, token.Type)
		}

		if token.Literal != testToken.expectedLit {
			t.Fatalf("testing value [%d] - literal wrong. expected=%q, got=%q", i+1, testToken.expectedLit, token.Literal)
		}
	}
}
