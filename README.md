This was a fun little team project to see how we can filter 
S3 inventory .csv.gz files fastest!

## Other implementations:

* [mythmon's Rust implementation](https://github.com/mythmon/rust-gz-csv-test)
* [peterbe's Python implementation](https://gist.github.com/peterbe/f147fd093aef43304a5c7e0a89c1ea0a) + [blog](https://www.peterbe.com/plog/fastest-python-datetime-parser)

## Usage

````
# get some, downloads 1GB from S3. So we can use the same input
# to benchmark implementations
> ./download.sh

# use a one file at a time strategy
> go run ./filter.go

# use a parallel (one worker / CPU) strategy
# add a GOPAR=1 env variable
> GOPAR=1 go run ./filter.go
````

## My results (on my late 2017 13" MBP)

````
Strategy: One file at a time ...
Total: 31521045, Matched: 710093, Ratio: 2.25%
Time: 52.740166887s
````

````
Strategy: Parallel, 4 Workers ...
Total: 31521045, Matched: 710093, Ratio: 2.25%
Time: 27.207802611s
````
