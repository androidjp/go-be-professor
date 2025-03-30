package kcfg

import (
	hash "github.com/mitchellh/hashstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"sync/atomic"
	"time"
)

type Value interface {
	Bool() (bool, error)
	Int64() (int64, error)
	Float64() (float64, error)
	String() (string, error)
	Duration() (time.Duration, error)
	Map() (map[string]interface{}, error)
	Scan(interface{}) error
	Load() interface{}
	Store(interface{})
	ValueFrom() string
	HashCode() uint64
	Error() error
}

type value struct {
	atomic.Value
	from string
	v    *viper.Viper
	key  string
}

func NewValue(v *viper.Viper, key string, from string) Value {
	return &value{
		v:    v,
		key:  key,
		from: from,
	}
}
func (v *value) Bool() (bool, error) {
	return cast.ToBoolE(v.Load())
}

func (v *value) Int64() (int64, error) {
	return cast.ToInt64E(v.Load())
}

func (v *value) Float64() (float64, error) {
	return cast.ToFloat64E(v.Load())
}

func (v *value) String() (string, error) {
	return cast.ToStringE(v.Load())
}

func (v *value) Duration() (time.Duration, error) {
	return cast.ToDurationE(v.Load())
}

func (v *value) Map() (map[string]interface{}, error) {
	return cast.ToStringMapE(v.Load())
}

func (v *value) Scan(val interface{}) error {
	if val == nil {
		return nil
	}

	if err := v.v.UnmarshalKey(v.key, val); err != nil {
		return err
	}

	return nil
}

func (v *value) Load() interface{} {
	return v.Value.Load()
}

func (v *value) Store(val interface{}) {
	v.Value.Store(val)
}

func (v *value) ValueFrom() string {
	return v.from
}

func (v *value) Error() error {
	return nil
}

func (v *value) HashCode() uint64 {
	code, err := hash.Hash(v.Load(), nil)
	if err != nil {
		panic(err)
	}
	return code
}

type errValue struct {
	err error
}

func (v errValue) Bool() (bool, error)                  { return false, v.err }
func (v errValue) Int64() (int64, error)                { return 0, v.err }
func (v errValue) Float64() (float64, error)            { return 0.0, v.err }
func (v errValue) Duration() (time.Duration, error)     { return 0, v.err }
func (v errValue) String() (string, error)              { return "", v.err }
func (v errValue) Scan(interface{}) error               { return v.err }
func (v errValue) Load() interface{}                    { return nil }
func (v errValue) Store(interface{})                    {}
func (v errValue) Map() (map[string]interface{}, error) { return nil, v.err }
func (v errValue) ValueFrom() string                    { return "err value" }
func (v errValue) HashCode() uint64                     { return 0 }
func (v errValue) Error() error                         { return v.err }
