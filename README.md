[![Go Report Card](https://goreportcard.com/badge/github.com/nicholaspcr/GDE3)](https://goreportcard.com/report/github.com/nicholaspcr/GDE3)
[![codecov](https://codecov.io/gh/nicholaspcr/GDE3/branch/master/graph/badge.svg?token=X96TDQSMFI)](https://codecov.io/gh/nicholaspcr/GDE3)

# GDE3 - Third Generalized Differentianl Evolution

gde3 is the golang cli for running the algorithm GDE3 with a set of well
defined problems of the literature, comporting the addition of external
problems that implement the interface methods.

As an addition there are a few python files reponsible for generating graphs of
the plots and performance indicators in related to the output of the gde3 cli.

For more details on how to execute each of these please refer to the `make
help` information that is provided by the Makefile. Observation, it requires
python3 to generate the text.


### Requirements



### TODO
- Create Makefile with commands:
    - [ ] Install gde3
    - [ ] Run
        - [ ] installs
        - [ ] runs the gde3 with default parameters
    - [ ] Python Graphs
        - [ ] Plots
        - [ ] Hypervolume
        - [ ] IGDPlus
    - [ ] Tests
        - [ ] General command for running tests
        - [ ] Generate tests via `pymoo`
    - [ ] Help command
    - [ ] Setup command
        - [ ] Install python libraries needed

- Add requirements for using the makefile
    - [ ] Golang minimal version
    - [ ] Python3 minimal version

- refactor the golang sections
    - [ ] check if its using proper interfaces
    - [ ] is it functional as it can be?
    - [ ] proper testing, use framework or just write base tests

