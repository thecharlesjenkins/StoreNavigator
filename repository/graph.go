package repository

//StoreGraph contains all data pertaining to the graph of stores
//Both the graph and the store have an Items structure. The Items in the Store is a
//simple list of Items, while the Items in this struct is a map that allows quick access to a specific item by ID.
type StoreGraph struct {
	GraphedStore *Store
	Items        map[int]Item `json:"items"`
	AdjList      []Vertex     `json:"adjList"`
	StartVertex  *Vertex      `json:"start"`
	EndVertex    *Vertex      `json:"end"`
}

//point is a simple structure to hold a row and column
type point struct {
	row    int
	column int
}

//Vertex contains the data at the vertex as well as distances to all neighbor vertices
type Vertex struct {
	StoredItem Item             `json:"item"`
	Neighbors  []VertexDistance `json:"neighbors"`
}

//VertexDistance contains the distance to a specific vertex
type VertexDistance struct {
	DestinationVertex Vertex `json:"item"`
	Distance          int    `json:"distance"`
}

//GraphStore will only be successful if path is connected
//Uses breadth-first-search to combine path and items into graph
//TODO: Doesn't need to use breadth-first-search, can just iteratate over all entries
func (s *Store) GraphStore() StoreGraph {
	graph := StoreGraph{}
	if len(s.Path) > 0 && len(s.Path[0]) > 0 {
		graph.GraphedStore = s
		height := len(s.Path)
		width := len(s.Path[0])

		graph.AdjList = make([]Vertex, height*width)
		var j int //Row
		i := 0    //Column

		memo := make(map[point]*Item, len(s.Items))

		index := 0

		for j = 0; j < len(s.Path); j++ {
			for i < len(s.Path[j]) {
				neighborVertex := graph.AdjList[getIndex(&s.Path, i, j)]
				graph.pushNeighbors(neighborVertex, i, j-1)
				graph.pushNeighbors(neighborVertex, i, j+1)
				graph.pushNeighbors(neighborVertex, i-1, j)
				graph.pushNeighbors(neighborVertex, i+1, j)
				graph.AdjList[getIndex(&s.Path, i, j)] = Vertex{StoredItem: graph.getItem(memo, &index, i, j)}
			}
		}
	}

	return graph
}

//getIndex converts a row and column to a single number to be accessed easily from list
func getIndex(path *PathArray, row int, column int) int {
	return row*len((*path)[row]) + column
}

//getItem expects an reference to an index that will be able to consistently store a value.
//Caches all of the incremented values to lazily build the graph and retrieve the item simultaneously.
//Also uses memoization to quickly access already cached Items that have not yet been used
func (graph *StoreGraph) getItem(memo map[point]*Item, index *int, row int, column int) Item {
	var itemAtIndex *Item = memo[point{row: row, column: column}]
	for itemAtIndex == nil && *index < len(graph.GraphedStore.Items) {
		if graph.GraphedStore.Items[*index].Row == row && graph.GraphedStore.Items[*index].Column == column {
			itemAtIndex = &graph.GraphedStore.Items[*index]
		}
		graph.Items[graph.GraphedStore.Items[*index].ID] = graph.GraphedStore.Items[*index]
		memo[point{row: row, column: column}] = &graph.GraphedStore.Items[*index]
		*index++
	}
	if itemAtIndex == nil {
		itemAtIndex = &Item{Type: graph.GraphedStore.Path[row][column].getItemType(), Row: row, Column: column}
	}
	return *itemAtIndex
}

//pushNeighbors adds all neighbors of a vertex to the graph
func (graph *StoreGraph) pushNeighbors(neighborVertex Vertex, row int, column int) {
	if row >= 0 && row < len(graph.GraphedStore.Path) && column >= 0 &&
		column < len(graph.GraphedStore.Path[row]) && graph.GraphedStore.Path[row][column].getItemType() != Wall {
		vertex := graph.AdjList[getIndex(&graph.GraphedStore.Path, row, column)]
		distance := abs(neighborVertex.StoredItem.Row-row) + abs(neighborVertex.StoredItem.Column-column)
		neighborVertex.Neighbors = append(neighborVertex.Neighbors, VertexDistance{DestinationVertex: vertex, Distance: distance})
	}
}

//abs value of an integer
func abs(num int) int {
	if num < 0 {
		return num * -1
	}
	return num
}
