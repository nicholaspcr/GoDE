import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import os

problem = "zdt1"
variants = ["rand1", "rand2", "best1", "best2", "currtobest1" , "pbest/P-0.05", "pbest/P-0.1", "pbest/P-0.15", "pbest/P-0.2"]

GEN = 100

dir_path = os.path.dirname(os.path.realpath(__file__))
file_path = dir_path + '/files/' + problem + '.csv'
HVData = pd.read_csv(file_path)

x = np.linspace(0,GEN,GEN)

for var in variants:
    plt.scatter(x, HVData[var], s=2, alpha=0.6, label=var)

plt.title(problem)
plt.suptitle("Average Hypervolume per Generations")
plt.xlabel("generations")
plt.ylabel("HV")
plt.legend()
plt.show()

