package main

// The store stands in for the primary database, and it is slow ON PURPOSE.
// Every query pays queryCost (DB_DELAY_MS, default 30ms) — the simulated
// round-trip to a database that is not on your machine — and then scans the
// whole table, because nothing here is indexed and nothing is cached.
//
// Do not "fix" this file by deleting the delay. The briefs on this rung are
// about not paying it on every request, and your load-test baseline is only
// honest if the slow path stays slow.

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const seedRowsDefault = 50000

var categories = []string{"beauty", "electronics", "fashion", "grocery", "home", "sports", "toys"}

type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	PricePaise int     `json:"price_paise"`
	Rating     float64 `json:"rating"`
	InStock    bool    `json:"in_stock"`
}

type Store struct {
	products  []Product
	queryCost time.Duration
}

// NewStore seeds the dataset deterministically: every clone of this repo gets
// exactly the same rows, so numbers quoted in a pull request are reproducible.
func NewStore() *Store {
	rows := envInt("SEED_ROWS", seedRowsDefault)
	delay := envInt("DB_DELAY_MS", 30)
	rng := rand.New(rand.NewSource(42))
	products := make([]Product, rows)
	for i := range products {
		products[i] = Product{
			ID:         i + 1,
			Name:       fmt.Sprintf("Product %d", i+1),
			Category:   categories[rng.Intn(len(categories))],
			PricePaise: 999 + rng.Intn(500000),
			Rating:     1 + rng.Float64()*4,
			InStock:    rng.Intn(10) > 1,
		}
	}
	return &Store{products: products, queryCost: time.Duration(delay) * time.Millisecond}
}

// Get pays the round-trip cost, then scans the whole table for one row.
func (s *Store) Get(id int) (Product, bool) {
	time.Sleep(s.queryCost)
	for _, p := range s.products { // full scan — no index
		if p.ID == id {
			return p, true
		}
	}
	return Product{}, false
}

// ByCategory returns every matching row. No pagination, no cache.
func (s *Store) ByCategory(category string) []Product {
	time.Sleep(s.queryCost)
	out := []Product{}
	for _, p := range s.products {
		if category == "" || p.Category == category {
			out = append(out, p)
		}
	}
	return out
}

func envInt(name string, fallback int) int {
	if v := os.Getenv(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return n
		}
	}
	return fallback
}
