package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type InOperator struct {
	*baseOperator
	not bool
}

func NewInOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	not bool,
) (*InOperator, error) {
	if core.IsNil(left) {
		return nil, core.Error(core.ErrMissedArgument, "left expression")
	}

	if core.IsNil(right) {
		return nil, core.Error(core.ErrMissedArgument, "right expression")
	}

	return &InOperator{&baseOperator{src, left, right}, not}, nil
}

func (operator *InOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return values.False, core.SourceError(operator.src, err)
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return values.False, core.SourceError(operator.src, err)
	}

	err = core.ValidateType(right, core.ArrayType)

	if err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	arr := right.(*values.Array)
	found := arr.IndexOf(left) > -1

	if operator.not {
		return values.NewBoolean(found == false), nil
	}

	return values.NewBoolean(found == true), nil
}
