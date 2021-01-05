package evolution

import (
	"math/rand"
)

// Population is a list of all competing genomes.
type Population []Lifeform

// Lifeform is a genome scored by the fitness function.
type Lifeform struct {
	genes Genome
	score float32
}

// GenPool is the current generation of all genes.
type GenPool []Genome

// Genome is the building block to generate a life form.
type Genome []byte

// FitnessFunc calculates the fitness score of a genome. Higher values indicate a more desiered behavior.
type FitnessFunc func(g Genome) float32

// ParentSelector choses a parent from a polulation for reproduction.
type ParentSelector func(overallScrore float32, p Population) Lifeform

// GenExchange creates a new genome from a set of parents.
type GenExchange func(mother, father Lifeform) Genome

// Mutation changes a genome slightly.
type Mutation func(g Genome) Genome

// NewPopulation creates a initial population of genomes to start the evolution.
func NewPopulation(r rand.Rand, genSize, populationSize int, fit FitnessFunc) Population {
	population := make([]Lifeform, populationSize, populationSize)

	for i := 0; i < populationSize; i++ {
		genome := make([]byte, genSize, genSize)
		r.Read(genome)

		lifeform := Lifeform{
			genes: genome,
			score: fit(genome),
		}

		population = append(population, lifeform)
	}

	return population
}

// PreferHigherScore selects a parent based on its fitness score. Higher values increase the
// propability of selection. This creates the pressure required for evolution.
func PreferHigherScore(r rand.Rand, overallScore float32, p Population) Lifeform {
	offset := r.Float32() * overallScore

	for _, lifeform := range p {
		offset -= lifeform.score
		if offset < 0 {
			return lifeform
		}
	}

	return p[len(p)-1]
}

// SimpleGenExchange creates a new genome based on mothers and fathers genome.
// The beginning is chosen from the mother, the end from the father.
// The pivot Byte also starts with the mother and end with the fathers bits.
func SimpleGenExchange(r *rand.Rand, mother, father Lifeform) Genome {
	genSize := len(mother.genes)
	genome := make([]byte, genSize, genSize)

	pivotBytes := r.Int() % genSize
	pivotBits := r.Int() % 9 // 9 so allow shifting to generate all values from 0 to 255

	copy(genome, father.genes)
	copy(genome, mother.genes[0:pivotBytes])

	var pivotBitsMask byte = 255 << pivotBits
	pivotByte := (mother.genes[pivotBytes] & pivotBitsMask) | (father.genes[pivotBytes] & ^pivotBitsMask)
	genome[pivotBytes] = pivotByte

	return genome
}

// Evolve generates a new generation of lifeforms.
func (p Population) Evolve(fit FitnessFunc, ps ParentSelector, exch GenExchange, mut Mutation) Population {
	populationSize := len(p)
	population := make([]Lifeform, populationSize, populationSize)
	overallScore := p.overallScore()

	for i := 0; i < populationSize; i++ {
		mother := ps(overallScore, p)
		father := ps(overallScore, p)
		genome := exch(mother, father)
		genome = mut(genome)

		lifeform := Lifeform{
			genes: genome,
			score: fit(genome),
		}

		population = append(population, lifeform)
	}

	return population
}

func (p Population) overallScore() float32 {
	var score float32
	for _, lifeform := range p {
		score += lifeform.score
	}
	return score
}

// BestScored returns the lifeform with the highest fitness score value.
func (p Population) BestScored() Lifeform {
	lifeform := p[0]

	for _, competitor := range p {
		if competitor.score > lifeform.score {
			lifeform = competitor
		}
	}

	return lifeform
}
