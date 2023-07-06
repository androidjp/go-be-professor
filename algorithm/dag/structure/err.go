package structure

import (
	"fmt"
	"strings"
)

func ErrorCycleDependent(root string, dependChain []string) error {
	sb := strings.Builder{}
	for _, item := range dependChain {
		sb.WriteString(item)
		sb.WriteString(" --> ")
	}
	sb.WriteString(root)
	return fmt.Errorf("key=%s 存在循环依赖: %s", root, sb.String())
}
