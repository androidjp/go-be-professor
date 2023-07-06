package structure

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDAG_BFS(t *testing.T) {
	Convey("should 返回正常", t, func() {
		Convey("given A是需要被检测的节点，A->B->C", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddEdge(A, B)
			dag.AddEdge(B, C)

			m, err := dag.BFS(A)
			So(err, ShouldBeNil)
			So(m, ShouldHaveLength, 3)
		})
		Convey("given A是需要被检测的节点，A->B->C->B", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddEdge(A, B)
			dag.AddEdge(B, C)
			dag.AddEdge(C, B)

			m, err := dag.BFS(A)
			So(err, ShouldBeNil)
			So(m, ShouldHaveLength, 3)
		})
	})
	Convey("should 报错：key=B 存在循环依赖: B --> C --> B", t, func() {

		Convey("given B是需要被检测的节点，A->B->C->B", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddEdge(A, B)
			dag.AddEdge(B, C)
			dag.AddEdge(C, B)

			_, err := dag.BFS(B)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "key=B 存在循环依赖: B --> C --> B")
		})
	})
	Convey("should 报错：key=A 存在循环依赖: A --> B --> C --> A", t, func() {
		Convey("given A是需要被检测的节点，A->B->C->A", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddEdge(A, B)
			dag.AddEdge(B, C)
			dag.AddEdge(C, A)

			_, err := dag.BFS(A)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "key=A 存在循环依赖: A --> B --> C --> A")
		})
	})
	Convey("should 报错：key=1 存在循环依赖: 1 --> 5 --> 3 --> 4 --> 2 --> 1", t, func() {
		Convey("given v1是需要被检测的节点，v1->v5, v1->v3->v2->v1", func() {
			dag := &DAG{}
			v1 := &Vertex{Key: "1"}
			v2 := &Vertex{Key: "2"}
			v3 := &Vertex{Key: "3"}
			v4 := &Vertex{Key: "4"}
			v5 := &Vertex{Key: "5"}

			// 对于有环图，从root点出发，最终会回到root
			// non dag
			//    5
			//   >
			//  /
			// 1 <----2
			//  \   .>  \
			//   > /     >
			//    3----  >4
			dag.AddEdge(v1, v5)
			dag.AddEdge(v2, v1)
			dag.AddEdge(v1, v3)
			dag.AddEdge(v3, v4)
			dag.AddEdge(v3, v2)
			dag.AddEdge(v2, v4)
			dag.AddEdge(v2, v1)

			_, err := dag.BFS(v1)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "key=1 存在循环依赖: 1 --> 5 --> 3 --> 4 --> 2 --> 1")
		})
	})

	Convey("should 报错：A强依赖B，客户端渠道不匹配", t, func() {
		Convey("given A是待上线插件，需要匹配版本号A和渠道1，而A-->B-->C，B此时只匹配版本号A和渠道2", func() {
			A := &Vertex{Key: "A", Value: map[string]string{"cli_ver": "A", "cli_chan": "1"}}
			A.AddVerifyFunc(StaticStrategyCheck)
			B := &Vertex{Key: "B", Value: map[string]string{"cli_ver": "A", "cli_chan": "2"}}
			C := &Vertex{Key: "C", Value: map[string]string{"cli_ver": "A", "cli_chan": "2"}}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddEdge(A, B)
			dag.AddEdge(B, C)

			_, err := dag.BFS(A)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "A强依赖B，客户端渠道不匹配")
		})
	})
}

func TestDAG_DFS(t *testing.T) {
	Convey("should 返回正常", t, func() {
		Convey("given A是需要被检测的节点，A->B->C", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddEdge(A, B)
			dag.AddEdge(B, C)

			_, err := dag.DFS(A)
			So(err, ShouldBeNil)
		})
		Convey("given A是需要被检测的节点，A->B->C->B", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddEdge(A, B)
			dag.AddEdge(B, C)
			dag.AddEdge(C, B)

			_, err := dag.DFS(A)
			So(err, ShouldBeNil)
		})
	})

	Convey("should 报错：key=B 存在循环依赖: B --> C --> B", t, func() {

		Convey("given B是需要被检测的节点，A->B->C->B", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddEdge(A, B)
			dag.AddEdge(B, C)
			dag.AddEdge(C, B)

			_, err := dag.DFS(B)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "key=B 存在循环依赖: B --> C --> B")
		})
	})
	Convey("should 报错：key=A 存在循环依赖: A --> B --> C --> A", t, func() {
		Convey("given A是需要被检测的节点，A->B->C->A", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddEdge(A, B)
			dag.AddEdge(B, C)
			dag.AddEdge(C, A)

			_, err := dag.DFS(A)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "key=A 存在循环依赖: A --> B --> C --> A")
		})
	})
	Convey("should 报错：key=1 存在循环依赖: 1 --> 3 --> 2 --> 1", t, func() {
		Convey("given v1是需要被检测的节点，v1->v5, v1->v3->v2->v1", func() {
			dag := &DAG{}
			v1 := &Vertex{Key: "1"}
			v2 := &Vertex{Key: "2"}
			v3 := &Vertex{Key: "3"}
			v4 := &Vertex{Key: "4"}
			v5 := &Vertex{Key: "5"}

			// 对于有环图，从root点出发，最终会回到root
			// non dag
			//    5
			//   >
			//  /
			// 1 <----2
			//  \   .>  \
			//   > /     >
			//    3----  >4
			dag.AddEdge(v1, v5)
			dag.AddEdge(v2, v1)
			dag.AddEdge(v1, v3)
			dag.AddEdge(v3, v4)
			dag.AddEdge(v3, v2)
			dag.AddEdge(v2, v4)
			dag.AddEdge(v2, v1)

			_, err := dag.DFS(v1)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "key=1 存在循环依赖: 1 --> 3 --> 2 --> 1")
		})
	})
}

func TestDAG_BFS2(t *testing.T) {
	Convey("should 返回正常 并最终返回路径中走过的所有节点信息，包含A在内一共8个节点", t, func() {
		Convey("given A是需要被检测的节点，A->B->C->F->G，交叉但不循环依赖", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			D := &Vertex{Key: "D"}
			E := &Vertex{Key: "E"}
			F := &Vertex{Key: "F"}
			G := &Vertex{Key: "G"}
			H := &Vertex{Key: "H"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddVertex(D)
			dag.AddVertex(E)
			dag.AddVertex(F)
			dag.AddVertex(G)
			dag.AddVertex(H)
			dag.AddEdge(A, B)
			dag.AddEdge(A, D)
			dag.AddEdge(A, H)
			dag.AddEdge(B, C)
			dag.AddEdge(B, E)
			dag.AddEdge(D, C)
			dag.AddEdge(D, E)
			dag.AddEdge(C, F)
			dag.AddEdge(E, F)
			dag.AddEdge(F, G)
			/**
			A -> B -> C -> F -> G
			  \    X    /
			    D ->  E
			  \
			    H
			*/

			m, err := dag.BFS(A)
			So(err, ShouldBeNil)
			So(m, ShouldHaveLength, 8)
		})
	})

	Convey("should 返回正常 并最终返回路径中走过的所有节点信息，包含A在内一共6个节点", t, func() {
		Convey("given A是需要被检测的节点，交叉但不循环依赖", func() {
			A := &Vertex{Key: "A"}
			B := &Vertex{Key: "B"}
			C := &Vertex{Key: "C"}
			D := &Vertex{Key: "D"}
			E := &Vertex{Key: "E"}
			F := &Vertex{Key: "F"}
			G := &Vertex{Key: "G"}
			H := &Vertex{Key: "H"}
			dag := &DAG{}

			dag.AddVertex(A)
			dag.AddVertex(B)
			dag.AddVertex(C)
			dag.AddVertex(D)
			dag.AddVertex(E)
			dag.AddVertex(F)
			dag.AddVertex(G)
			dag.AddVertex(H)
			dag.AddEdge(A, B)
			dag.AddEdge(A, D)
			dag.AddEdge(A, H)
			dag.AddEdge(B, C)
			dag.AddEdge(B, E)
			dag.AddEdge(D, C)
			dag.AddEdge(D, E)
			dag.AddEdge(F, G)
			/**
			A -> B -> C     F -> G
			  \    X
			    D ->  E
			  \
			    H
			*/

			m, err := dag.BFS(A)
			So(err, ShouldBeNil)
			So(m, ShouldHaveLength, 6)
		})
	})

}
