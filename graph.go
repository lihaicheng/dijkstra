package dijkstra

import (
	"errors"
	"fmt"
)

// Graph contains all the graph details
type Graph struct {
	best        int64
	visitedDest bool
	//slice of all verticies available
	Verticies       []Vertex
	visiting        dijkstraList
	mapping         map[string]int
	usingMap        bool // 使用从string类型ID到int类型ID的映射
	highestMapIndex int
	running         bool
}

// NewGraph creates a new empty graph
func NewGraph() *Graph {
	new := &Graph{}
	new.mapping = map[string]int{}
	return new
}

// AddNewVertex adds a new vertex at the next available index
func (g *Graph) AddNewVertex() *Vertex {
	for i, v := range g.Verticies {
		if i != v.ID {
			g.Verticies[i] = Vertex{ID: i}
			return &g.Verticies[i]
		}
	}
	return g.AddVertex(len(g.Verticies))
}

// AddVertex adds a single vertex
func (g *Graph) AddVertex(ID int) *Vertex {
	g.AddVerticies(Vertex{ID: ID})
	return &g.Verticies[ID]
}

// GetVertex gets the reference of the specified vertex. An error is thrown if
// there is no vertex with that index/ID.
func (g *Graph) GetVertex(ID int) (*Vertex, error) {
	if ID >= len(g.Verticies) {
		return nil, errors.New("Vertex not found")
	}
	return &g.Verticies[ID], nil
}

// 这段代码是一个用于验证图的合法性的方法。
// 主要的目标是确保图中的每一条边的目标节点是有效的，
// 即目标节点的索引在图的节点范围内，并且没有指向 ID 为 0 的节点（这里的 ID 为 0 的节点可能被认为是未初始化的节点）
func (g Graph) validate() error {
	for _, v := range g.Verticies {
		for a := range v.arcs {
			// 目标节点的索引 a 是否超出了图中节点的范围，
			// 节点a不是0号节点但是节点的ID是0，这是为了防止默认值 这是需要强调的
			if a >= len(g.Verticies) || (g.Verticies[a].ID == 0 && a != 0) {
				return errors.New(fmt.Sprint("Graph validation error;", "Vertex ", a, " referenced in arcs by Vertex ", v.ID))
			}
		}
	}
	return nil
}

// SetDefaults sets the distance and best node to that specified
func (g *Graph) setDefaults(Distance int64, BestNode int) {
	for i := range g.Verticies {
		g.Verticies[i].bestVerticies = []int{BestNode}
		g.Verticies[i].distance = Distance
	}
}
