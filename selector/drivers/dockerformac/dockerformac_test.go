package dockerformac

import "testing"

func TestName(t *testing.T) {
	driver := NewDriver()
	if driver.Name() != "docker-for-mac" {
		t.Errorf("Invalid name")
	}
}

func TestList(t *testing.T) {
	driver := NewDriver()
	list := driver.List()
	
	l := 0
	for _ = range list {
		l = l + 1
	}
	if l != 1 {
		t.Errorf("Invalid list, count=%v", l)
	}
}
