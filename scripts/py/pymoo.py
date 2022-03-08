import numpy as np
import pandas as pd
import os
import sys
import json
import matplotlib.pyplot as plt

# d = {}
# d["A"] = np.array([1,2,3,4,5])
# d["B"] = np.array([2,3,4,5,6])
# d["C"] = np.array([3,4,5,6,7])

# df = pd.DataFrame(d, columns= ['A', 'B', 'C'])
# dir_path = os.path.dirname(os.path.realpath(__file__))
# print(dir_path)
# df.to_csv(dir_path + "/dtlz1.csv")
# print(df)

df = pd.read_csv('dtlz1.csv')
x = np.linspace(1,5,5)
arr = ['A', 'B', 'C']
for a in arr:
    plt.scatter(x, df[a], s=2, alpha=0.5, label=a)

plt.title("test")
plt.suptitle("Average Hypervolume per Generations")
plt.xlabel("generations")
plt.ylabel("HV")
plt.legend()    
plt.show()