package tests

import (
	genv1 "GRPC_Calc/proto/gen"
	"GRPC_Calc/tests/suite"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestCalculations_HappyPath(t *testing.T) {
	const (
		mn = 0
		mx = 100
	)

	ctx, st := suite.New(t)

	number := gofakeit.Number(mn, mx)
	uid := strconv.Itoa(gofakeit.Number(0, 1000))
	expr := fmt.Sprintf("%d + %d - %d * (%d + %d) / %d", number, number, number, number, number, number)

	respCalc, err := st.CalcClient.Calculate(ctx, &genv1.ExprRequest{
		Expr: expr,
		Uid:  uid,
	})
	require.NoError(t, err)

	ans := respCalc.GetAnswer()
	require.NotEmpty(t, ans)
}
