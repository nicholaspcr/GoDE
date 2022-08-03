package de

import "context"

type generationKeyType struct{}

var generationKey = &generationKeyType{}

func InjectGenerations(ctx context.Context, gen int) context.Context {
	return context.WithValue(ctx, generationKey, gen)
}

func FetchGenerations(ctx context.Context) int {
	v, ok := ctx.Value(generationKey).(int)
	if !ok {
		return 0
	}
	return v
}

type constFKeyType struct{}

var constFKey = &constFKeyType{}

func InjectFConst(ctx context.Context, F float64) context.Context {
	return context.WithValue(ctx, generationKey, F)
}

func FetchFConst(ctx context.Context) float64 {
	v, ok := ctx.Value(generationKey).(float64)
	if !ok {
		return 0.0
	}
	return v
}

type constPKeyType struct{}

var constPKey = &constPKeyType{}

func InjectPConst(ctx context.Context, P float64) context.Context {
	return context.WithValue(ctx, generationKey, P)
}

func FetchPConst(ctx context.Context) float64 {
	v, ok := ctx.Value(generationKey).(float64)
	if !ok {
		return 0.0
	}
	return v
}

type constCRKeyType struct{}

var constCRKey = &constCRKeyType{}

func InjectCRConst(ctx context.Context, CR float64) context.Context {
	return context.WithValue(ctx, generationKey, CR)
}

func FetchCRConst(ctx context.Context) float64 {
	v, ok := ctx.Value(generationKey).(float64)
	if !ok {
		return 0.0
	}
	return v
}
