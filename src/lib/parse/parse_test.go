package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainModule(t *testing.T) {
	for _, str := range []string{"", "x = 42\nf x = + x 123\nwrite 123"} {
		result, err := newStateWithoutFile(str).mainModule()()

		t.Log(result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestMainModuleFail(t *testing.T) {
	for _, str := range []string{"(", "(()"} {
		result, err := newStateWithoutFile(str).mainModule()()

		t.Log(err.Error())

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestSubModule(t *testing.T) {
	for _, str := range []string{"", "x = 123\nf x = + x 123"} {
		result, err := newStateWithoutFile(str).subModule()()

		t.Log(result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestXFailSubModule(t *testing.T) {
	for _, str := range []string{"(", "write 123"} {
		result, err := newStateWithoutFile(str).subModule()()

		t.Log(err.Error())

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestImportModule(t *testing.T) {
	for _, str := range []string{`import "foo"`, `import "foo/bar"`} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.importModule())()
		assert.Equal(t, nil, err)
	}
}

func TestLetFunctionOneLine(t *testing.T) {
	for _, str := range []string{
		"foo = 123",
		"foo x = x",
		"foo x = f x y",
		"foo x y (z 123) (v 456) ..args . a b (c 123) (d 456) ..kwargs = 123",
	} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.letFunctionOneLine())()
		assert.Equal(t, nil, err)
	}
}

func TestLetFunctionBlock(t *testing.T) {
	for _, str := range []string{
		"foo\n  = 123",
		"foo x\n  = x",
		"foo x\n  = f x y",
		"foo x y (z 123) (v 456) ..args . a b (c 123) (d 456) ..kwargs\n  = 123",
	} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.letFunctionBlock())()
		assert.Equal(t, nil, err)
	}
}

func TestEqualExpression(t *testing.T) {
	for _, str := range []string{
		"= 123",
		"= x",
		"= f x y",
		"= match 123 123 123",
	} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.equalExpression())()
		assert.Equal(t, nil, err)
	}
}

func TestSignature(t *testing.T) {
	for _, str := range []string{"", "x", "x y", "(x 123)", "..args", ". x", ". (x 123)", ". ..kwargs", "..args . ..kwargs"} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.signature())()
		assert.Equal(t, nil, err)
	}
}

func TestSignatureFail(t *testing.T) {
	for _, str := range []string{
		"=",
	} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.signature())()
		assert.NotEqual(t, nil, err)
	}
}

func TestHalfSignatureFail(t *testing.T) {
	for _, str := range []string{
		"=",
	} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.halfSignature())()
		assert.NotEqual(t, nil, err)
	}
}

func TestExpandedArgumentFail(t *testing.T) {
	for _, str := range []string{
		"=",
	} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.Many(s.identifier()))()
		assert.NotEqual(t, nil, err)
	}
}

func TestOutput(t *testing.T) {
	for _, str := range []string{"output", "..outputs", "foo bar", "..(foo bar)"} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.output())()
		assert.Equal(t, nil, err)
	}
}

func TestStringLiteral(t *testing.T) {
	for _, str := range []string{`""`, `"sl"`, "\"   string literal  \"", `"\""`, `"\\"`} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.stringLiteral())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestStrip(t *testing.T) {
	s := newStateWithoutFile("ident   ")
	result, err := s.Exhaust(s.strip(s.identifier()))()

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

func TestList(t *testing.T) {
	for _, str := range []string{"[]", "[123 456]", "[(f x) 123]"} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestExpression(t *testing.T) {
	for _, str := range []string{
		"ident",
		"ident  ",
		"foo bar",
		"foo ; (this is) comment \n     bar   ; lsdfj\n ",
	} {
		t.Logf("source: %#v", str)

		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestExpressionFail(t *testing.T) {
	for _, str := range []string{
		"=",
	} {
		t.Logf("source: %#v", str)

		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", err)

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestMatchExpression(t *testing.T) {
	for _, str := range []string{
		"match 123 123 true",
		"match (foo bar) [123 ..elems] (process elems) xs (write xs)",
		"match (foo bar) [\"foo\" 123 ..rest] (process rest) xs (write xs)",
	} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.match())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestApp(t *testing.T) {
	for _, str := range []string{
		"f x", "f x y", "f ..x", "f . x 123", "f . x 123 y 456",
		"func . ..kwargs", "f ..x (func x y) 123 456 ..foo . a 123 b 456 ..c ..(d 123 456 789)"} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.app())()
		t.Logf("%#v", result)
		assert.Equal(t, err, nil)
	}
}

func TestArguments(t *testing.T) {
	for _, str := range []string{
		"",
		"x",
		"x y",
		"..x",
		". x 123",
		". x 123 y 456",
		". ..kwargs",
		"..x (func x y) 123 456 ..foo . a 123 b 456 ..c ..(d 123 456 789)"} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.arguments())()
		t.Logf("%#v", result)
		assert.Equal(t, err, nil)
	}
}

func TestIdentifier(t *testing.T) {
	for _, str := range []string{
		"ident",
		"123",
		"nil",
		"true",
	} {
		result, err := newStateWithoutFile(str).identifier()()

		t.Log(err)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestIdentifierFail(t *testing.T) {
	for _, str := range []string{"", ".", "..", ".foo", "="} {
		s := newStateWithoutFile(str)
		result, err := s.identifier()()
		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestBlank(t *testing.T) {
	for _, str := range []string{"", "   ", "\n\n", " ; laskdjf \n  "} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.blank())()

		t.Log(result, err)

		assert.Equal(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func newStateWithoutFile(source string) *state {
	return newState("", source)
}
