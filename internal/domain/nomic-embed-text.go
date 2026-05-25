package domain

type Document struct {
	ID      string
	Content string
	Vector  []float64
}

type VectorDB struct {
	documents []Document
}
