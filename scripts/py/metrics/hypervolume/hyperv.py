import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from pymoo.factory import get_problem, get_performance_indicator, get_reference_directions, get_visualization
from pymoo.performance_indicator.hv import Hypervolume
import os
import sys

# important methods -> ref_point
ref_point = np.array([ 2.9699152058577947 , 4.985423795813952,  6.988314095918014])
metric = Hypervolume(ref_point=ref_point, normalize=False)
problem = "wfg1"

# 1 -> [488.4435867485893 ,471.03860764532106 , 501.59485141330845 ]
# 2 -> [ 2.7640053099955044, 2.5431893060817545 ,  2.676556552035201]
# 3 -> [ 1692.099007341335 , 1678.2523471910183 , 1812.180711990459 ]
# 4 -> [ 2.80704182512513 , 2.6906135709301386 , 2.534695795120166]
# 5 -> [ 2.352652061652649 , 2.504239438737515 , 2.7298272947925213]
# 6 -> [ 10.548625453871612 , 10.612209406392823 , 10.634665459827534 ]
# 7 -> [ 1, 1, 27.9038366074846 ]
# WFGS
# 1 -> [ 2.9699152058577947 , 4.985423795813952,  6.988314095918014]

# general constants
NUM_EXECS = 30
variants = ["rand1", "rand2", "best1", "best2", "currtobest1", "pbest/P-0.05", "pbest/P-0.1", "pbest/P-0.15", "pbest/P-0.2"] 
# 
base_path = "/home/nick/.gode/mode/paretoFront/" + problem + "/"
# data for the plot
HVData = []

# file related constants      
GEN = 251
NP = 100

for varIndex in range(len(variants)):
    execFiles = []


    # reads from each execution of the variant
    for i in range(NUM_EXECS):
        filePath = base_path + variants[varIndex] + "/exec-" + str(i+1) + '.csv'
        execFiles.append(
            pd.read_csv(
                filePath,
                sep=';|\n', engine='python'
            )
        )

    # data of variant
    variantHV =  np.array([])

    for i in range(GEN):
        genHV = 0.0

        for file in execFiles:
            # A = np.reshape(file["A"], (NP,GEN))
            A = file["A"]
            B = file["B"]
            C = file["C"]
            genData = []
            for j in range(NP*i, NP*(i+1)):
                genData.append(np.array([
                    A[j],
                    B[j],
                    C[j]
                ]))
            #genData = np.array([
            #        A[i*NP:(i+1)*NP],
            #        B[i*NP:(i+1)*NP],
            #        C[i*NP:(i+1)*NP]
            #])
            genHV = genHV + metric.calc(np.array(genData))

        genHV = genHV / float(len(execFiles))
        variantHV = np.append(variantHV, genHV)

    HVData.append(variantHV)

    #for gen in range(0, GEN):
    #    # sum of the avarage HV of each element
    #    genHV = 0.0
    #    for file in execFiles:
    #        populationArray = []
    #        for i in range(NP):
    #            populationArray.append(np.array([
    #                file.iloc[(gen*3)][i],
    #                file.iloc[(gen*3)+1][i],
    #                file.iloc[(gen*3)+2][i],
    #            ]))

    #        variantArray = np.array(populationArray)
    #        genHV = genHV + metric.calc(variantArray)
    #    genHV = genHV / float(len(execFiles))
    #    variantHV = np.append(variantHV, genHV)

    #HVData.append(variantHV)

for i in range(len(HVData)):
    x = np.linspace(0, len(HVData[i]), len(HVData[i]))
    plt.scatter(x, HVData[i], s=2, alpha=0.6, label=variants[i])

plt.title(problem)
plt.suptitle("Average Hypervolume per Generations")
plt.xlabel("generations")
plt.ylabel("HV")
plt.legend()
plt.show()
