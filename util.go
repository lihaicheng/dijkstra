package dijkstra

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Import imports a graph from the specified file returns the Graph, a map for
// if the nodes are not integers and an error if needed.
func Import(filename string) (g Graph, err error) {
	g.usingMap = false
	var lowestIndex int // lowestIndex 用于跟踪节点映射的最低整数索引
	var i int           // i 是当前正在处理的节点的索引
	var arc int         // arc 是边的目标节点的索引
	var dist int64      // dist 是边的权重
	var ok bool         //ok 是一个布尔值，表示在映射中查找是否成功
	got, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	g.mapping = map[string]int{}

	input := strings.TrimSpace(string(got))
	for _, line := range strings.Split(input, "\n") {
		f := strings.Fields(strings.TrimSpace(line))
		if len(f) == 0 || (len(f) == 1 && f[0] == "") {
			continue
		}
		//no need to check for size cause there must be something as the string is trimmed and split
		if g.usingMap {
			if i, ok = g.mapping[f[0]]; !ok {
				g.mapping[f[0]] = lowestIndex
				i = lowestIndex
				lowestIndex++
			}
		} else {
			// 两个边界
			// 如果少了一行怎么办 见*1行，如果节点的数量小于等于当前节点索引 i，则扩展节点切片。
			// 如果转化失败怎么办 自动转成使用string-int映射
			i, err = strconv.Atoi(f[0])
			if err != nil {
				g.usingMap = true
				g.mapping[f[0]] = lowestIndex
				i = lowestIndex
				lowestIndex++
			}
		}
		if temp := len(g.Verticies); temp <= i { //Extend if we have to // *1
			// 创建1+i-len(g.Verticies)个空节点追加到g.Verticies中
			// 对于普通情况，就是为第i个节点创建一个节点实例
			// 但是对于漏掉一行，如M.txt中没有下标为3的行，此时temp == len(g.Verticies) == 3， i == 4 , 1+i-len(g.Verticies) == 2
			g.Verticies = append(g.Verticies, make([]Vertex, 1+i-len(g.Verticies))...) // 此时 len(g.Verticies) == 5
			for ; temp < len(g.Verticies); temp++ {                                    // 3 < 5
				g.Verticies[temp].ID = temp
				g.Verticies[temp].arcs = map[int]int64{}
				g.Verticies[temp].bestVerticies = []int{-1}
			}
		}
		if len(f) == 1 {
			//if there is no FROM here
			continue
		}
		for _, set := range f[1:] {
			got := strings.Split(set, ",")
			if len(got) != 2 {
				err = ErrWrongFormat
				return
			}
			dist, err = strconv.ParseInt(got[1], 10, 64)
			if err != nil {
				err = ErrWrongFormat
				return
			}
			if g.usingMap {
				arc, ok = g.mapping[got[0]]
				if !ok {
					arc = lowestIndex
					g.mapping[got[0]] = arc
					lowestIndex++
				}
			} else {
				arc, err = strconv.Atoi(got[0])
				if err != nil {
					err = ErrMixMapping
					return
				}
			}
			g.Verticies[i].arcs[arc] = dist
		}
	}
	err = g.validate()
	return
}

// ExportToFile exports the verticies to file currently does not take into account
// mappings (from string to int)
func (g Graph) ExportToFile(filename string) error {
	var i string
	var err error
	os.MkdirAll(filename, 0777)
	if _, err = os.Stat(filename); err == nil {
		os.Remove(filename)
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, v := range g.Verticies {
		if g.usingMap {
			if i, err = g.GetMapped(v.ID); err != nil {
				return errors.New("mapping fail when exporting; " + err.Error())
			}
			fmt.Fprint(f, i)
		} else {
			fmt.Fprint(f, v.ID)
		}
		for key, val := range v.arcs {
			if g.usingMap {
				if i, err = g.GetMapped(key); err != nil {
					return errors.New("mapping fail when exporting; " + err.Error())
				}
				fmt.Fprint(f, " ", i, ",", val)
			} else {
				fmt.Fprint(f, " ", key, ",", val)
			}
		}
		fmt.Fprint(f, "\n")
	}
	return nil
}
