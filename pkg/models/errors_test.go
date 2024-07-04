package models

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmallUnit_RegValidateErrAtIdentifier(t *testing.T) {

	testCases := []struct {
		Expr              string
		ExpectedIsMatched bool
	}{
		{Expr: `abc`, ExpectedIsMatched: true},
		{Expr: `abc_def`, ExpectedIsMatched: true},
		{Expr: `abc_`, ExpectedIsMatched: false},
		{Expr: `_abc`, ExpectedIsMatched: false},
		{Expr: `abc.def`, ExpectedIsMatched: true},
		{Expr: `abc.`, ExpectedIsMatched: false},
		{Expr: `abc[10]`, ExpectedIsMatched: true},
		{Expr: `abc[-10]`, ExpectedIsMatched: true},
		{Expr: `abc[-10.5]`, ExpectedIsMatched: false},
		{Expr: `abc['10']`, ExpectedIsMatched: true},
		{Expr: `abc['\\']`, ExpectedIsMatched: true},
		{Expr: `abc['\']`, ExpectedIsMatched: false},
		{Expr: `abc['_']`, ExpectedIsMatched: true},
		{Expr: `abc[-10].def`, ExpectedIsMatched: true},
		{Expr: `abc[-10].def.`, ExpectedIsMatched: false},
	}

	for _, testCase := range testCases {

		testCase := testCase

		t.Run(fmt.Sprintf("Expr=`%s`", testCase.Expr), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.ExpectedIsMatched, regValidateErrAt.MatchString(testCase.Expr))
		})
	}
}

func TestSmallUnit_ModelValidateError(t *testing.T) {

	rootErr := NewModelValidateError(errors.New("validation error message"))
	parentNodeCaseAErr := WrapModelValidateError("childA", rootErr)
	parentNodeCaseBErr := WrapModelValidateError("childB.", rootErr)
	parentNodeCaseCErr := WrapModelValidateError(".childC", rootErr)
	grandParentNodeCaseAA := WrapModelValidateError("childA", parentNodeCaseAErr)

	testCases := []struct {
		Err             error
		ExpectedMessage string
	}{
		{Err: rootErr, ExpectedMessage: "validation error: validation error message"},
		{Err: parentNodeCaseAErr, ExpectedMessage: "validation error in 'childA': validation error message"},
		{Err: parentNodeCaseBErr, ExpectedMessage: "validation error in 'childB': validation error message"},
		{Err: parentNodeCaseCErr, ExpectedMessage: "validation error in 'childC': validation error message"},
		{Err: grandParentNodeCaseAA, ExpectedMessage: "validation error in 'childA.childA': validation error message"},
	}

	for _, testCase := range testCases {

		testCase := testCase

		t.Run(fmt.Sprintf("ExpectedMessage='%s'", testCase.ExpectedMessage), func(t *testing.T) {
			t.Parallel()
			assert.ErrorContains(t, testCase.Err, testCase.ExpectedMessage)
		})
	}
}
