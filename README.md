[![Go Report Card](https://goreportcard.com/badge/github.com/nicholaspcr/GDE3)](https://goreportcard.com/report/github.com/nicholaspcr/GDE3)
[![codecov](https://codecov.io/gh/nicholaspcr/GDE3/branch/master/graph/badge.svg?token=X96TDQSMFI)](https://codecov.io/gh/nicholaspcr/GDE3)

# Deprecated

This repository is reprecated in favor of the
[GoDE](https://github.com/nicholaspcr/GoDE). The reason is that this repository
is refered in papers and therefore can't be modified in order to comport the
execution of different algorithms.


# GDE3 - Third Generalized Differential Evolution

gde3 is the golang cli for running the algorithm GDE3 with a set of well
defined problems of the literature, comporting the addition of external
problems that implement the interface methods.

As an addition there are a few python files reponsible for generating graphs of
the plots and performance indicators in related to the output of the gde3 cli.

For more details on how to execute each of these please refer to the `make
help` information that is provided by the Makefile. Observation, it requires
python3 to generate the text.


