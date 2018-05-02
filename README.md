This was a fun little team project to see how we can filter 
S3 inventory .csv.gz files fastest!

## Other implementations:

* [mythmon's Rust implementation](https://github.com/mythmon/rust-gz-csv-test)
* [peterbe's Python implementation](https://gist.github.com/peterbe/f147fd093aef43304a5c7e0a89c1ea0a) + [blog](https://www.peterbe.com/plog/fastest-python-datetime-parser)

## Usage

```
# get some working data, downloads 1GB from S3 into testdata/ subdirectory
> ./download.sh



# Processing using a one file at a time
> go run ./filter.go


# Processing in parallel (workers = num cpus)
> GOPAR=1 go run ./filter.go
```

## My results (on my late 2017 13" MBP)

```
Strategy: One file at a time ...
Total: 31521045, Matched: 710093, Ratio: 2.25%
Time: 52.740166887s
```

```
Strategy: Parallel, 4 Workers ...
Total: 31521045, Matched: 710093, Ratio: 2.25%
Time: 27.207802611s
```
