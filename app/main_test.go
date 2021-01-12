package main

import (
	"go.uber.org/goleak"
	"testing"
)

func TestMain(m *testing.M)  {
	goleak.VerifyTestMain(m)
}

func Test_Main(t *testing.T) {
}