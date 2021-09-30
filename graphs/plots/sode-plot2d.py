import pandas as pd
import matplotlib.pyplot as plt
import argparse
import numpy as np
import time

filePath = "/home/nick/.go-de/sode/allPoints/rastrigin/dim-2/rand1.csv"

f = pd.read_csv(
    filePath,
    sep='\t|\n', engine='python'
)

number_rows, number_cols = f.shape
print(number_rows, number_cols)

X = f.iloc[number_rows-2][3:]
Y = f.iloc[number_rows-1][3:]


# plt.title("ZDT-1")
plt.suptitle("Frente de Pareto")
plt.xlabel("obj-1")
plt.ylabel("obj-2")
plt.scatter(X, Y, color='b', s=2, alpha=0.6)
plt.show()
