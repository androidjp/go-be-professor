package structure

import "fmt"

func StaticStrategyCheck(a, b *Vertex) error {
	aSSMap, ok1 := a.Value.(map[string]string)
	bSSMap, ok2 := b.Value.(map[string]string)
	if !ok1 || !ok2 {
		return fmt.Errorf("%s或%s的静态过滤器丢失", a.Key, b.Key)
	}
	info := fmt.Sprintf("%s强依赖%s", a.Key, b.Key)

	if aSSMap["cli_ver"] != bSSMap["cli_ver"] {
		return fmt.Errorf(info + "，客户端版本不匹配")
	}

	if aSSMap["cli_chan"] != bSSMap["cli_chan"] {
		return fmt.Errorf(info + "，客户端渠道不匹配")
	}

	return nil
}
