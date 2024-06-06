package main 

// Node - has the data about the backend server

type Node struct {
  URL *url.URL
  Active bool 
  weight float64
}

// NodePool - Has slice of the nodes and most recently used node indexx

type NodePool struct {
  nodes []*Node 
  current uint64
}

func main () {
    
}

