package file_test

import (
	"fmt"
	"io"
	"mylib/scfg/file"
	"mylib/scfg/kcfg"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSourceLoad(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		wantErr   error
		wantKeyVs []*kcfg.KeyValue
	}{
		{
			name:     "normal test case 1",
			filePath: "../../mock/cfg/file_test.yaml",
			wantErr:  nil,
			wantKeyVs: []*kcfg.KeyValue{
				{
					Key:    "file_test.yaml",
					Format: "yaml",
					Value:  readFile("../../mock/cfg/file_test.yaml"),
				},
			},
		},
		{
			name:     "normal test case 2",
			filePath: "../../mock/cfg/file_test.json",
			wantErr:  nil,
			wantKeyVs: []*kcfg.KeyValue{
				{
					Key:    "file_test.json",
					Format: "json",
					Value:  readFile("../../mock/cfg/file_test.json"),
				},
			},
		},
		{
			name:     "bad test case 3",
			filePath: "../../mock/cfg",
			wantErr:  fmt.Errorf("file source can not load file path is a dir"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := file.NewSource(tt.filePath)
			gotKeyVs, gotErr := fs.Load()
			assert.Equal(t, tt.wantErr, gotErr)

			for i, wantKeyV := range tt.wantKeyVs {
				assert.Equal(t, wantKeyV.Key, gotKeyVs[i].Key)
				assert.Equal(t, wantKeyV.Format, gotKeyVs[i].Format)
				assert.Equal(t, wantKeyV.Value, gotKeyVs[i].Value)
			}
		})
	}

}

func readFile(path string) []byte {
	f, _ := os.Open(path)
	defer f.Close()
	data, _ := io.ReadAll(f)
	return data
}
