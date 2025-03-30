package slogger

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"reflect"
	"unsafe"
)

//----------------------------------------------------------------------------------------
// AllKeysReadContextHandler
//----------------------------------------------------------------------------------------

type AllKeysReadContextHandler struct {
	slog.Handler
}

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h *AllKeysReadContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// v2版本：读取所有的key-value (存在问题，待修复)
	m := GetKeyValues(ctx)
	for k, v := range m {
		r.AddAttrs(slog.Any(convertToString(k), v))
	}
	return h.Handler.Handle(ctx, r)
}

func convertToString(v interface{}) string {
	switch value := v.(type) {
	case string:
		return value
	case int:
		return fmt.Sprintf("%d", value)
	case int8:
		return fmt.Sprintf("%d", value)
	case int16:
		return fmt.Sprintf("%d", value)
	case int32:
		return fmt.Sprintf("%d", value)
	case int64:
		return fmt.Sprintf("%d", value)
	case float32:
		return fmt.Sprintf("%f", value)
	case float64:
		return fmt.Sprintf("%f", value)
	case bool:
		return fmt.Sprintf("%t", value)
	case nil:
		return "nil"
	default:
		t := reflect.TypeOf(value)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Struct {
			method := reflect.ValueOf(value).MethodByName("String")
			if method.IsValid() && method.Type().NumIn() == 0 && method.Type().NumOut() == 1 && method.Type().Out(0).Kind() == reflect.String {
				result := method.Call(nil)
				return result[0].String()
			}
			return fmt.Sprintf("%+v", value)
		}
		return fmt.Sprintf("%+v", value)
	}
}

type iface struct {
	itab, data uintptr
}
type emptyCtx int

type valueCtx struct {
	context.Context
	key, val interface{}
}

func GetKeyValues(ctx context.Context) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	getKeyValue(ctx, m)
	return m
}

func getKeyValue(ctx context.Context, m map[interface{}]interface{}) {
	ictx := *(*iface)(unsafe.Pointer(&ctx))
	if ictx.data == 0 || int(*(*emptyCtx)(unsafe.Pointer(ictx.data))) == 0 {
		return
	}

	valCtx := (*valueCtx)(unsafe.Pointer(ictx.data))
	if valCtx != nil && valCtx.key != nil {
		m[valCtx.key] = valCtx.val
	}
	getKeyValue(valCtx.Context, m)
}

//----------------------------------------------------------------------------------------
// DefaultContextHandler
//----------------------------------------------------------------------------------------

const (
	ctxKeyRequestID string = "request_id"
)

type DefaultContextHandler struct {
	ctxReadKeys     []string               // 配置层读取的context key（仅支持string）
	extraReadKeyMap map[string]interface{} // 动态为logger设置的context key（其中：k-displayKey，v-key）
	slog.Handler
}

func (h *DefaultContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// 默认：根据配置的key读取value并添加到Record中
	for _, key := range h.ctxReadKeys {
		if v, ok := ctx.Value(key).(string); ok {
			r.AddAttrs(slog.Any(key, v))
		}
	}

	// 额外：根据配置的key读取value并添加到Record中
	if h.extraReadKeyMap != nil {
		for displayKey, key := range h.extraReadKeyMap {
			if v, ok := ctx.Value(key).(string); ok {
				r.AddAttrs(slog.Any(displayKey, v))
			}
		}
	}
	return h.Handler.Handle(ctx, r)
}

func (h *DefaultContextHandler) addExtraCtxReadKey(displayKey string, key interface{}) {
	var mp map[string]interface{}
	if h.extraReadKeyMap == nil {
		mp = make(map[string]interface{}, 1)
	} else {
		mp = make(map[string]interface{}, len(h.extraReadKeyMap)+1)
		mp[displayKey] = key
		maps.Copy(mp, h.extraReadKeyMap)
	}
	h.extraReadKeyMap = mp
}

func NewContextHandler(cfg *Config, h slog.Handler) slog.Handler {
	switch cfg.CtxReadType {
	// case ContextReadTypeAllKeys:
	// 	return &AllKeysReadContextHandler{
	// 		Handler: h,
	// 	}
	default:
		keys := make([]string, 0, len(cfg.CtxReadKeys)+1)
		keys = append(keys, ctxKeyRequestID)
		keys = append(keys, cfg.CtxReadKeys...)
		extraReadKeyMap := make(map[string]interface{}, 0)
		return &DefaultContextHandler{
			ctxReadKeys:     keys,
			extraReadKeyMap: extraReadKeyMap,
			Handler:         h,
		}
	}
}
