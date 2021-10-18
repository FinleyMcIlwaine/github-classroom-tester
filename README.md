# GitHub Classroom Tester

Haskell and Go applications that run GitHub classroom autograding test commands locally.

## Install and Use

### 1. Clone the repo:

```
git clone git@github.com:FinleyMcIlwaine/github-classroom-tester.git tester
```

### 2. Install with `cabal` or `go` (after changing directories to where you cloned the project):

Haskell:
```
cd tester/tester-haskell
cabal install --installdir <path-to-your-project>
```

Go:
```
cd tester/tester-go
go build -o <path-to-your-project>
```

Where `<path-to-your-project>` is the relative or absolute path to the directory containing the project with the autograding.

### 3. Run the tester:

```
cd <path-to-your-project>
./tester
```

## Output

The tool will output the results of all tests to standard out. A log of the test results is written to `tester-log.txt` in the directory the `tester` command was ran in.

![tester.gif](gif of tester)
