package env_test

import (
	"fmt"
	"mylib/scfg/env"
	"mylib/scfg/kcfg"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv_Load(t *testing.T) {
	prefix1, prefix2 := "ENV_", "KLION_"

	testData := map[string]string{
		prefix1 + "xxx":    "aaaaa",
		prefix2 + "ererer": "rrrrr",
		"cccc":             "bbbbb",
	}

	for k, v := range testData {
		os.Setenv(k, v)
	}

	envSource := env.NewSource(prefix1, prefix2)
	sr := kcfg.NewSourceReader(envSource)
	if err := sr.Init(); err != nil {
		t.Fatal(err)
	}

	t.Run("case1 test get from env", func(t *testing.T) {
		for k, v := range testData {
			avtualV, err := sr.GetValue(k)
			if err != nil {
				expectedErr := fmt.Errorf("key [%s] not found from [%s] config source", k, sr.Source().SourceName())
				assert.Equal(t, expectedErr, err)
				continue
			}
			actualStr, _ := avtualV.String()
			assert.Equal(t, v, actualStr)
		}
	})

}
