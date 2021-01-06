package evolution

import (
	"encoding/binary"
	"log"
	"math"
	"math/rand"
)

// Population is a list of all competing genomes.
type Population []Lifeform

// Lifeform is a genome scored by the fitness function.
type Lifeform struct {
	Genes      Genome
	Evaluation float32
	Fitness    float32
}

// GenPool is the current generation of all genes.
type GenPool []Genome

// Genome is the building block to generate a life form.
type Genome []byte

// EvaluationFunc evaluates the genome.
type EvaluationFunc func(g GenomeReader) float32

// FitnessFunc calculates the fitness score based on the evaluation score.
// Higher values indicate a more desiered behavior.
type FitnessFunc func(value, min, max float32) float32

// ParentSelector choses a parent from a polulation for reproduction.
type ParentSelector func(r *rand.Rand, overallScrore float32, p Population) int

// GenExchange creates a new genome from a set of parents.
type GenExchange func(r *rand.Rand, mother, father Lifeform) Genome

// Mutation changes a genome slightly.
type Mutation func(r *rand.Rand, g Genome) Genome

// NewPopulation creates a initial population of genomes to start the evolution.
func NewPopulation(r *rand.Rand, genSize, populationSize int, eva EvaluationFunc, fit FitnessFunc) Population {
	population := make([]Lifeform, populationSize)

	for i := 0; i < populationSize; i++ {
		log.Printf("lifeform %d", i)

		genome := make([]byte, genSize)
		r.Read(genome)

		lifeform := Lifeform{
			Genes:      genome,
			Evaluation: eva(GenomeReader{Gens: genome}),
		}

		population[i] = lifeform
	}

	var min float32 = math.MaxFloat32
	var max float32 = 0
	for _, lifeform := range population {
		if lifeform.Evaluation < min {
			min = lifeform.Evaluation
		}
		if lifeform.Evaluation > max {
			max = lifeform.Evaluation
		}
	}

	for i := range population {
		population[i].Fitness = fit(population[i].Evaluation, min, max)
	}

	return population
}

// SimpleParentSelector selects a parent based on its fitness score. Higher values increase the
// propability of selection. This creates the pressure required for evolution.
func SimpleParentSelector(r *rand.Rand, overallScore float32, p Population) int {
	offset := r.Float32() * overallScore

	for i, lifeform := range p {
		offset -= lifeform.Fitness
		if offset < 0 {
			return i
		}
	}

	return len(p) - 1
}

// SimpleGenExchange creates a new genome based on mothers and fathers genome.
// The beginning is chosen from the mother, the end from the father.
// The pivot Byte also starts with the mother and end with the fathers bits.
func SimpleGenExchange(r *rand.Rand, mother, father Lifeform) Genome {
	genSize := len(mother.Genes)
	genome := make([]byte, genSize)

	pivotBytes := r.Int() % genSize
	pivotBits := r.Int() % 9 // 9 so allow shifting to generate all values from 0 to 255

	copy(genome, father.Genes)
	copy(genome, mother.Genes[0:pivotBytes])

	var pivotBitsMask byte = 255 << pivotBits
	pivotByte := (mother.Genes[pivotBytes] & pivotBitsMask) | (father.Genes[pivotBytes] & ^pivotBitsMask)
	genome[pivotBytes] = pivotByte

	return genome
}

// SimpleMutation flipes a random bit with a 1 in 100 propability.
func SimpleMutation(r *rand.Rand, g Genome) Genome {
	if r.Int()%2 != 0 {
		return g
	}

	genSize := len(g)

	pivotBytes := r.Int() % genSize
	pivotBits := r.Int() % 9 // 9 so allow shifting to generate all values from 0 to 255

	var pivotBitsMask byte = 1 << pivotBits

	g[pivotBytes] ^= pivotBitsMask

	return SimpleMutation(r, g)
}

// Evolve generates a new generation of lifeforms.
func (p Population) Evolve(
	r *rand.Rand,
	eva EvaluationFunc,
	fit FitnessFunc,
	ps ParentSelector,
	exch GenExchange,
	mut Mutation,
) Population {
	populationSize := len(p)
	population := make([]Lifeform, populationSize)
	overallScore := p.overallFitness()

	for i := 0; i < populationSize; i++ {
		log.Printf("lifeform %d", i)

		motherIdx := ps(r, overallScore, p)
		fatherIdx := ps(r, overallScore, p)

		for motherIdx == fatherIdx {
			motherIdx = ps(r, overallScore, p)
		}

		mother := p[motherIdx]
		father := p[fatherIdx]

		genome := mother.Genes
		for testEq(genome, mother.Genes) || testEq(genome, father.Genes) {
			genome = exch(r, mother, father)
			genome = mut(r, genome)
		}

		lifeform := Lifeform{
			Genes:      genome,
			Evaluation: eva(GenomeReader{Gens: genome}),
		}

		population[i] = lifeform
	}

	var min float32 = math.MaxFloat32
	var max float32 = 0
	for _, lifeform := range population {
		if lifeform.Evaluation < min {
			min = lifeform.Evaluation
		}
		if lifeform.Evaluation > max {
			max = lifeform.Evaluation
		}
	}

	for i, lifeform := range population {
		population[i].Fitness = fit(lifeform.Evaluation, min, max)
	}

	return population
}

func testEq(a, b Genome) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func (p Population) overallFitness() float32 {
	var fitness float32
	for _, lifeform := range p {
		fitness += lifeform.Fitness
	}
	return fitness
}

// Fittest returns the lifeform with the highest fitness score.
func (p Population) Fittest() Lifeform {
	lifeform := p[0]

	for _, competitor := range p {
		if competitor.Fitness > lifeform.Fitness {
			lifeform = competitor
		}
	}

	return lifeform
}

type GenomeReader struct {
	Gens Genome
	idx  int
}

func (gr *GenomeReader) Byte() byte {
	b := gr.Gens[gr.idx]
	gr.idx += 1
	return b
}

func (gr *GenomeReader) Uint8() uint8 {
	b := gr.Gens[gr.idx]
	gr.idx += 1
	return uint8(b)
}

func (gr *GenomeReader) Uint16() uint16 {
	bytes := gr.Gens[gr.idx : gr.idx+2]
	gr.idx += 2
	return binary.LittleEndian.Uint16(bytes)
}

func (gr *GenomeReader) Float64() float64 {
	bytes := gr.Gens[gr.idx : gr.idx+8]
	gr.idx += 8
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func (gr *GenomeReader) Length() int {
	return len(gr.Gens)
}

func (gr *GenomeReader) Seek(i int) {
	gr.idx = i
}
