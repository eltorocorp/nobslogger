#!/bin/bash

sed -E '{
	/Benchmark/!d
	s/BenchmarkAccumulatedContext\///g
	s/ (ns|B|allocs)\/op//g
}' | \
sort -bnk 3
