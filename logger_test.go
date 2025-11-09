package kit

import (
	"context"
	"fmt"
	"testing"
)

func Test_Clogger_WithCtx(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel, Context: true})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())
	l := L(logger).C(ctx)
	l.Inf("I'm logger")
}

func Test_Clogger_WithCaller(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel, Context: true, Caller: true})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())
	l := L(logger).C(ctx)
	l.Inf("I'm logger")
}

func Test_SetLevel(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel, Context: true})
	L(logger).Trc("I'm logger")
	logger.SetLevel(DebugLevel)
	L(logger).Trc("I'm logger")
}

func Test_Clogger_Plain(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel, Format: "plain", Context: true})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())
	l := L(logger).C(ctx).F(KV{"field": "value"})
	l.Inf("I'm logger")
}

func Test_Clogger_Plain_Err(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel, Format: "plain", Context: true})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())
	err := NewAppErrBuilder("ERR-123", "%s happened", "shit").C(ctx).Business().F(KV{"f": "v"}).Err()
	l := L(logger).E(err).St()
	l.Err("my bad")
}

func Test_Clogger_WithComponentAndMethod(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})
	l := L(logger).Cmp("service").Mth("do")
	l.Inf("I'm logger")
}

func Test_Clogger_WithComponentMethodAndCtx(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel, Context: true})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())
	l := L(logger).Cmp("service").Mth("do").C(ctx)
	l.Inf("I'm logger")
}

func Test_Clogger_All(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel, Context: true})
	ctx := NewRequestCtx().WithNewRequestId().WithSessionId(NewId()).WithUser("1", "john").ToContext(context.Background())
	l := L(logger).Cmp("service").Mth("do").C(ctx).F(KV{"field": "value"})
	l.Inf("I'm logger")
}

func Test_Clogger_WithFields(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})
	l := L(logger).F(KV{"field": "value"})
	l.Inf("I'm logger")
}

func Test_Clogger_WithMultipleFields(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})
	l := L(logger).F(KV{"field": "value"}).Inf("I'm logger")
	l.F(KV{"field2": "value2"}).Inf("I'm logger")
}

func Test_Clogger_WithErr(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})
	l := L(logger).E(fmt.Errorf("error"))
	l.Err("my bad")
}

type hookImpl struct{}

func (h *hookImpl) Error(err error) {
	fmt.Println("hook error:", err)
}

func Test_Clogger_WithErrHook(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})
	logger.SetErrorHook(&hookImpl{})
	L(logger).E(fmt.Errorf("error")).Err("my bad")
}

func Test_Clogger_WithErrStack(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})
	l := L(logger).E(fmt.Errorf("error")).St()
	l.Err("my bad")
}

func Test_Clogger_WithAppErr(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())
	err := NewAppErrBuilder("ERR-123", "%s happened", "shit").C(ctx).Business().F(KV{"f": "v"}).Err()
	l := L(logger).E(err)
	l.Err("my bad")
}

func Test_Clogger_WithAppErrAndStack(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())
	err := NewAppErrBuilder("ERR-123", "%s happened", "shit").C(ctx).Business().F(KV{"f": "v"}).Err()
	l := L(logger).E(err).St()
	l.Err("my bad")
}

func Test_Clogger_WithTraceObj(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})

	obj1 := struct {
		A string
		B int
	}{
		A: "test",
		B: 5,
	}

	type n struct {
		A string
	}
	type s struct {
		Nested *n
	}

	obj2 := &s{
		Nested: &n{
			A: "str",
		},
	}

	L(logger).TrcObj("objects: %v, %v", obj1, obj2)
	L(logger).F(KV{"obj1": obj1, "obj2": obj2}).Trc("objects")
}

func Test_Clogger_WithAppErrAndFields(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel})
	e := NewAppErrBuilder("ERR-123", "%s happened", "shit").F(KV{"f": "v"}).Err()
	l := L(logger).E(e)
	l.Err("my bad")
}

func Test_Clogger_WithAppErrAndAppContext(t *testing.T) {
	logger := InitLogger(&LogConfig{Level: TraceLevel, Context: true})
	ctx := NewRequestCtx().WithRequestId("123").ToContext(context.Background())
	e := NewAppErrBuilder("ERR-123", "%s happened", "shit").C(ctx).Err()
	l := L(logger).E(e).St()
	l.Err("my bad")
}

func Test_Clogger_All2(t *testing.T) {

	logger := InitLogger(&LogConfig{Level: TraceLevel, Context: true})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())

	bigObj := struct {
		A string
		B int
		C []string
		D map[string]string
		E []map[string]string
	}{
		A: "test",
		B: 5,
		C: []string{"a", "b", "c"},
		D: map[string]string{"a": "b", "c": "d"},
		E: []map[string]string{
			{"a": "b", "c": "d"},
		},
	}

	L(logger).Cmp("service").Mth("do").C(ctx).F(KV{"field1": "value1", "field2": "value2"}).DbgF("log iteration").F(KV{"obj": bigObj}).Dbg("big object")

}

func Benchmark_Clogger_All(b *testing.B) {

	logger := InitLogger(&LogConfig{Level: TraceLevel, Context: true, Format: "json"})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())

	bigObj := struct {
		A string
		B int
		C []string
		D map[string]string
		E []map[string]string
	}{
		A: "test",
		B: 5,
		C: []string{"a", "b", "c"},
		D: map[string]string{"a": "b", "c": "d"},
		E: []map[string]string{
			{"a": "b", "c": "d"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := L(logger).Cmp("service").Mth("do").C(ctx).F(KV{"field1": "value1", "field2": "value2"})
		for i := 0; i < 10; i++ {
			l.DbgF("log iteration %d", i).F(KV{"obj": bigObj}).Dbg("big object")
		}
	}

}

func Benchmark_Clogger_All2(b *testing.B) {

	logger := InitLogger(&LogConfig{Level: TraceLevel, Context: true, Format: "json"})
	ctx := NewRequestCtx().WithNewRequestId().WithUser("1", "john").ToContext(context.Background())

	l := L(logger).Cmp("service").Mth("do").C(ctx).F(KV{"field1": "value1", "field2": "value2"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.F(KV{"id": "123"}).Dbg("big object")
	}

}
