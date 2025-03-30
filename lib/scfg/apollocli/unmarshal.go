package apollocli

import (
	"github.com/mitchellh/mapstructure"
	"time"
)

func GetUnmarshalFunc(f UnmarshalFunc) UnmarshalFunc {
	return func(inputBytes []byte, output interface{}) error {
		decoder, err := mapstructure.NewDecoder(defaultDecoderConfig(output))
		if err != nil {
			return err
		}

		var target interface{}
		if err := f(inputBytes, &target); err != nil {
			return err
		}
		return decoder.Decode(target)
	}
}

// defaultDecoderConfig returns default mapsstructure.DecoderConfig with suppot of time.Duration、string slices、time.Time
func defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),     // string: 1s、1ns to time.Duration
			mapstructure.StringToTimeHookFunc(time.RFC1123), // default use rfc1123 time format
			mapstructure.StringToSliceHookFunc(","),         // default string: a,b,c to slice: string[]{a,b,c}
		),
	}
	return c
}
