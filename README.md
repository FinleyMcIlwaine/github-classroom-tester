# GitHub Classroom Tester

Haskell application that runs GitHub classroom autograding test commands locally.

## Install

### 1. Clone the repo:

```
git clone git@github.com:FinleyMcIlwaine/github-classroom-tester.git tester
```

### 2. Install with `cabal` (after changing directories to where you cloned the project):

```
cd tester
cabal install --installdir <path-to-your-project>
```

Where `<path-to-your-project>` is the relative or absolute path to the directory containing the project with the autograding.

### 3. Run the tester:

```
cd <path-to-your-project>
./tester
```

## Output

The tool will output the results of all tests to standard out. A log of the test results is written to `testlog.txt` in the directory the `tester` command was ran in.

<img src="https://s9.gifyu.com/images/tester.gif"/>
