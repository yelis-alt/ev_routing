package dto

// GenerationDTO is one candidate path in the genetic search.
type GenerationDTO struct {
	Dna     map[int]int   // vertex id -> x (1/0 included)
	Parents []map[int]int // prior generations, for crossover
}
