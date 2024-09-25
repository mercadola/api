package test

import "github.com/mercadola/api/pkg/exceptions"

type CaseFunction struct {
	Name                string
	MockName            string
	InputParams         any
	Expected            any
	ExpectedErr         any
	AssertNumberOfCalls int
}

type Case string

const (
	Success Case = "success"
	Failure Case = "failure"
)

type TestCase struct {
	// METODO A SER TESTADO
	Name        string
	InputParams any
	Case        Case

	// MOCK DAS FUNÇÃO
	FunctionsCalledMock []CaseFunction

	// ASSERT
	Expected    any
	ExpectedErr *exceptions.AppException
}
