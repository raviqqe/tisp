package parse

import (
	"../ast"
	"./comb"
)

const (
	commentChar  = ';'
	invalidChars = "\x00"
	quoteString  = "quote"
	spaceChars   = " \t\n\r"
	specialChars = "()[]{}\"'`$"
)

func Parse(source string) []interface{} {
	m, err := newState(source).module()()

	if err != nil {
		panic(err.Error())
	}

	return m.([]interface{})
}

func (s *state) module() comb.Parser {
	return s.Exhaust(s.Wrap(s.blank(), s.expressions(), s.None()))
}

func (s *state) letConst() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewLetConst(xs[1].(string), xs[2])
	}, s.list(s.strippedString("let"), s.identifier(), s.expression()))
}

func (s *state) output() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		expanded := false

		if xs[0] != nil {
			expanded = true
		}

		return ast.NewOutput(xs[1], expanded)
	}, s.And(s.Maybe(s.String("..")), s.expression()))
}

func (s *state) expressions() comb.Parser {
	return s.Lazy(s.strictExpressions)
}

func (s *state) strictExpressions() comb.Parser {
	return s.Many(s.expression())
}

func (s *state) expression() comb.Parser {
	return s.strip(s.Or(
		s.firstOrderExpression(),
		s.Lazy(func() comb.Parser { return s.quote(s.expression()) })))
}

func (s *state) firstOrderExpression() comb.Parser {
	return s.Or(
		s.identifier(),
		s.String(".."),
		s.String("."),
		s.stringLiteral(),
		s.sequence("(", ")"),
		s.prepend("list", s.sequence("[", "]")),
		s.prepend("dict", s.sequence("{", "}")),
		s.prepend("set", s.sequence("'{", "}")),
		s.prepend("lambda", s.sequence("'(", ")")))
}

func (s *state) quote(p comb.Parser) comb.Parser {
	return s.And(s.Replace(quoteString, s.Char('`')), p)
}

func (s *state) identifier() comb.Parser {
	cs := string(commentChar) + invalidChars + spaceChars + specialChars
	return s.strip(s.Stringify(s.And(s.NotInString(cs+"."), s.Stringify(s.Many(s.NotInString(cs))))))
}

func (s *state) stringLiteral() comb.Parser {
	c := s.Char('"')
	f := func(x interface{}) interface{} {
		return []interface{}{quoteString, x}
	}

	return s.App(f, s.Stringify(s.Wrap(
		c,
		s.Many(s.Or(s.NotInString("\"\\"), s.String("\\\""), s.String("\\\\"))),
		s.strip(c))))
}

func (s *state) list(ps ...comb.Parser) comb.Parser {
	return s.Wrap(s.strippedString("("), s.And(ps...), s.strippedString(")"))
}

func (s *state) sequence(l, r string) comb.Parser {
	return s.Wrap(s.strippedString(l), s.expressions(), s.strippedString(r))
}

func (s *state) prepend(x interface{}, p comb.Parser) comb.Parser {
	return s.App(func(any interface{}) interface{} {
		return append([]interface{}{x}, any.([]interface{})...)
	}, p)
}

func (s *state) strip(p comb.Parser) comb.Parser {
	b := s.blank()
	return s.Wrap(s.None(), p, b)
}

func (s *state) blank() comb.Parser {
	return s.Void(s.Many(s.Or(s.InString(spaceChars), s.comment())))
}

func (s *state) comment() comb.Parser {
	return s.Void(s.And(
		s.Char(commentChar),
		s.Many(s.NotChar('\n')), s.Char('\n')))
}

func (s *state) strippedString(str string) comb.Parser {
	return s.strip(s.String(str))
}
