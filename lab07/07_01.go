package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	jrpc "github.com/AdamSLevy/jsonrpc2/v14"
)

var precision int = 2

type Args struct {
	A float64 `json:"x"`
	B float64 `json:"y"`
}

func (a *Args) UnmarshalJSON(data []byte) error {
	var slice []float64
	if err := json.Unmarshal(data, &slice); err == nil && len(slice) >= 2 {
		a.A = slice[0]
		a.B = slice[1]
		return nil
	}
	type alias Args
	return json.Unmarshal(data, (*alias)(a))
}

type PrecisionArgs struct {
	N int `json:"n"`
}

func (p *PrecisionArgs) UnmarshalJSON(data []byte) error {
	var slice []int
	if err := json.Unmarshal(data, &slice); err == nil && len(slice) >= 1 {
		p.N = slice[0]
		return nil
	}
	type alias PrecisionArgs
	return json.Unmarshal(data, (*alias)(p))
}

type Result struct {
	Value float64 `json:"val"`
}

func (res *Result) Precise() {
	p := math.Pow10(precision)
	res.Value = math.Round(p*res.Value) / p
}

func parseArgs(params json.RawMessage) (*Args, error) {
	var args Args
	if err := json.Unmarshal(params, &args); err != nil {
		return nil, err
	}
	return &args, nil
}

func methodSum(_ context.Context, params json.RawMessage) interface{} {
	args, err := parseArgs(params)
	if err != nil {
		return jrpc.ErrorInvalidParams(err)
	}
	res := Result{Value: args.A + args.B}
	res.Precise()
	log.Printf("[LOG] Sum: %v + %v = %v (prec: %d)", args.A, args.B, res.Value, precision)
	return res
}

func methodSub(_ context.Context, params json.RawMessage) interface{} {
	args, err := parseArgs(params)
	if err != nil {
		return jrpc.ErrorInvalidParams(err)
	}
	res := Result{Value: args.A - args.B}
	res.Precise()
	log.Printf("[LOG] Sub: %v - %v = %v (prec: %d)", args.A, args.B, res.Value, precision)
	return res
}

func methodMul(_ context.Context, params json.RawMessage) interface{} {
	args, err := parseArgs(params)
	if err != nil {
		return jrpc.ErrorInvalidParams(err)
	}
	res := Result{Value: args.A * args.B}
	res.Precise()
	log.Printf("[LOG] Mul: %v * %v = %v (prec: %d)", args.A, args.B, res.Value, precision)
	return res
}

func methodDiv(_ context.Context, params json.RawMessage) interface{} {
	args, err := parseArgs(params)
	if err != nil {
		return jrpc.ErrorInvalidParams(err)
	}
	if args.B == 0 {
		return jrpc.NewError(-32602, "division by 0 restricted", args)
	}
	res := Result{Value: args.A / args.B}
	res.Precise()
	log.Printf("[LOG] Div: %v / %v = %v (prec: %d)", args.A, args.B, res.Value, precision)
	return res
}

func methodPre(_ context.Context, params json.RawMessage) interface{} {
	var args PrecisionArgs
	if err := json.Unmarshal(params, &args); err != nil {
		return jrpc.ErrorInvalidParams(err)
	}
	precision = args.N
	log.Printf("[LOG] Precision set to %d", precision)
	return nil
}

func main() {
	file, err := os.OpenFile("07_01.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("[FATAL] unable to open log file, aborting")
		os.Exit(1)
	}
	log.SetOutput(file)

	methods := jrpc.MethodMap{
		"Calculator.Sum": methodSum,
		"Calculator.Sub": methodSub,
		"Calculator.Mul": methodMul,
		"Calculator.Div": methodDiv,
		"Calculator.Pre": methodPre,
	}

	http.Handle("/rpc", jrpc.HTTPRequestHandler(methods, log.New(os.Stderr, "", 0)))

	fmt.Println("server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
