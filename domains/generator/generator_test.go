package generator

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()
	code = m.Run()
}

func TestGenerateShortURL(t *testing.T) {
	ctx := context.Background()
	svc := New()
	res := svc.Generate(ctx, time.Now().Nanosecond())
	require.NotEmpty(t, res)
	require.Len(t, res, 10) // make sure that the length is not more than 10
}

func TestGenerate_SameIDSameResult(t *testing.T) {
	ctx := context.Background()
	svc := New()
	id := time.Now().Nanosecond()
	res := svc.Generate(ctx, id)
	require.NotEmpty(t, res)
	require.Len(t, res, 10)

	res1 := svc.Generate(ctx, id)
	require.NotEmpty(t, res1)
	require.Len(t, res1, 10)

	// make sure same ID gives same response
	require.Equal(t, res, res1)
}

func TestGenerate_DiffID_DiffResult(t *testing.T) {
	ctx := context.Background()
	svc := New()

	res := svc.Generate(ctx, time.Now().Nanosecond())
	require.NotEmpty(t, res)
	require.Len(t, res, 10)

	res1 := svc.Generate(ctx, time.Now().Nanosecond())
	require.NotEmpty(t, res1)
	require.Len(t, res1, 10)

	// make sure same ID gives same response
	require.NotEqual(t, res, res1)
}

// func TestGenerate_AddMachineID(t *testing.T) {
// 	// 376673000
// 	ctx := context.Background()
// 	svc := New()

// 	res := svc.Generate(ctx, 50376673000)
// 	require.NotEmpty(t, res)
// 	require.Len(t, res, 10)
// }
