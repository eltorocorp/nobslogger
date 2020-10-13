#!/bin/bash

echo "Running nobslogger benchmarks..."
rm benchmarks.out 2> /dev/null
rm benchmarks.md 2> /dev/null
go test -bench Accum -benchmem | tee benchmarks.out
cat benchmarks.out | \
sed -E '{
	/Benchmark/!d
	s/BenchmarkAccumulatedContext\///g
	s/ (ns|B|allocs)\/op//g
}' | \
sort -bnk 3 | \
awk '
BEGIN{
    FS = "\t"
    print "|Package|Time|Time %|Allocations|"
    print "|-------|----|------|-----------|"
}
NR == 1 {
    baseline = $3
}
{
    pct = -100 * (1-($3/baseline))
    printf "|%s|%i \ ns/op|%5.0f%%|%i \ allocs/op|\n", $1, $3, pct, $5
}
' > benchmarks.md
echo "Benchmark results exported to ./benchmarks.md and ./benchmarks.md."