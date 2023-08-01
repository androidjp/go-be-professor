package structure

import (
	"fmt"
	"github.com/eapache/queue"
)

// 邻接表
type Vertex struct {
	Key        string       // 节点的唯一标识
	Parents    []*Vertex    // 节点的父节点
	Children   []*Vertex    // 节点的子节点
	Value      interface{}  // 节点的信息数据
	VerifyFunc []VerifyFunc // 遍历过程中的校验函数
}

func (v *Vertex) Verify(target *Vertex) error {
	if target == v {
		return nil
	}
	var err error
	for _, verifyFunc := range v.VerifyFunc {
		if err = verifyFunc(v, target); err != nil {
			return err
		}
	}
	return nil
}

type VerifyFunc func(a, b *Vertex) error

func (v *Vertex) AddVerifyFunc(f ...VerifyFunc) {
	v.VerifyFunc = append(v.VerifyFunc, f...)
}

type DAG struct {
	Vertexes []*Vertex
}

func (dag *DAG) AddVertex(v *Vertex) {
	dag.Vertexes = append(dag.Vertexes, v)
}

func (dag *DAG) AddEdge(from, to *Vertex) {
	from.Children = append(from.Children, to)

	to.Parents = append(from.Parents, from)
}

func (dag *DAG) BFS(root *Vertex) (map[string]bool, error) {
	q := queue.New()

	visitMap := make(map[string]bool)
	visitMap[root.Key] = true

	q.Add(root)

	// 依赖链路
	dependChain := make([]string, 0)

	for {
		if q.Length() == 0 {
			fmt.Println("done")
			break
		}
		current := q.Remove().(*Vertex)
		// 自定义检测逻辑
		if err := root.Verify(current); err != nil {
			return nil, err
		}
		// 加入依赖链路中，主要是作记录
		dependChain = append(dependChain, current.Key)

		//fmt.Println("bfs key", current.Key)

		for _, v := range current.Children {
			//fmt.Printf("from:%v to:%s\n", current.Key, v.Key)
			if v.Key == root.Key {
				// 报错：循环依赖
				return nil, ErrorCycleDependent(root.Key, dependChain)
			}
			if _, ok := visitMap[v.Key]; !ok {
				visitMap[v.Key] = true
				q.Add(v)
			}
		}
	}

	return visitMap, nil
}

func (dag *DAG) DFS(root *Vertex) (map[string]bool, error) {
	stack := []*Vertex{root}

	visitMap := make(map[string]bool)
	visitMap[root.Key] = true

	dependChain := make([]string, 0)

	for {
		if len(stack) == 0 {
			fmt.Println("done")
			break
		}
		if len(stack)-1 < 0 {
			panic("unexpected")
		}
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// 自定义检测逻辑
		if err := root.Verify(current); err != nil {
			return nil, err
		}
		// 加入依赖链路中，主要是作记录
		dependChain = append(dependChain, current.Key)

		//fmt.Println("dfs key", current.Key)

		for _, v := range current.Children {
			fmt.Printf("from:%v to:%s\n", current.Key, v.Key)
			if v.Key == root.Key {
				// 报错：循环依赖
				return nil, ErrorCycleDependent(root.Key, dependChain)
			}
			if _, ok := visitMap[v.Key]; !ok {
				visitMap[v.Key] = true
				//fmt.Println("add visit", v.Key)
				if v.Key == root.Key {
					//panic("back root")
					// 报错：循环依赖
					return nil, ErrorCycleDependent(root.Key, dependChain)
				}
				stack = append(stack, v)
			}
		}
	}
	return visitMap, nil
}
