# WC-Tool - Coding Challenge #1

This is my own version of the Unix command line tool `wc`. It is a solution to the first problem in [John Crickett's Coding Challenges](https://codingchallenges.fyi/challenges/challenge-wc/).

## How to build

Execute `make build` to create the executable. It will be saved in the `bin/` directory.

## How to run

Execute the binary with a file as an argument to get the result.

The following are the possible arguments:

```bash
# Outputs the number of bytes

> ./bin/ccwc -c test.text
341836 test.txt

# Outputs the number of line breaks

> ./bin/ccwc -l test.text
7137 test.txt

# Outputs the number of words

> ./bin/ccwc -w test.text
58159 test.txt

# Outputs with -c, -l and -w flags
> ./bin/ccwc test.text
7137 58159 341836 test.txt
```

You can also read from standard input

```bash
> cat test.txt | ./bin/ccwc -l
7137
```

## License

[MIT](LICENSE) © André Brandão
