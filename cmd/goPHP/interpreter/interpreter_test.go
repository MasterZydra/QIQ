package interpreter

import "testing"

func TestText(t *testing.T) {
	actual, err := NewInterpreter().Process("<html>...</html>")
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
	}
	expected := "<html>...</html>"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestEchoShortTag(t *testing.T) {
	actual, err := NewInterpreter().Process(`<html><?= "abc" ?><?= 42; ?></html>`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
	}
	expected := "<html>abc42</html>"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestEchoExpression(t *testing.T) {
	actual, err := NewInterpreter().Process(`<html><?php echo "abc", 42 ?><?php echo "def", 24; ?></html>`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
	}
	expected := "<html>abc42def24</html>"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestVariableDeclaration(t *testing.T) {
	actual, err := NewInterpreter().Process(`<?php $var = "hi"; $var = "hello"; echo $var, " world";`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
	}
	expected := "hello world"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}
