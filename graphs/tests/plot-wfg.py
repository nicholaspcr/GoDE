import pandas as pd
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
import argparse
import numpy as np
import time
from pymoo.factory import get_problem, get_reference_directions, get_visualization
from pymoo.util.plotting import plot


n_var = 12
n_obj = 3

# problems = ["wfg1", "wfg2", "wfg3", "wfg4", "wfg5", "wfg6", "wfg7", "wfg8", "wfg9" ]
problem = "wfg9"

p = get_problem(problem, n_var, n_obj)
ref_dirs = get_reference_directions("das-dennis", 3, n_partitions=12)
pf = p.pareto_front(ref_dirs)

print(pf)
X = pf[:, 0]
Y = pf[:, 1]
Z = pf[:, 2]

fig = plt.figure()
ax = Axes3D(fig)

ax.set_title("DTLZ")
plt.title("DTLZ")
ax.set_xlabel("obj-1")
ax.set_ylabel("obj-2")
ax.set_zlabel("obj-3")
ax.scatter(X, Y, Z, s=5)
plt.show()

