import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from pymoo.factory import get_problem, get_performance_indicator, get_reference_directions, get_visualization
from pymoo.performance_indicator.igd_plus import IGDPlus
import os


# general constants
problem = "zdt6"
NUM_EXECS = 30
NUM_OBJ = 2
N_VAR = 10

# file related constants      
GEN = 100
NP = 100


#  file path
base_path = "/home/nick/.gode/mode/paretoFront/" + problem + "/"
variants = ["rand1", "rand2", "best1", "best2", "currtobest1", "pbest/P-0.05", "pbest/P-0.1", "pbest/P-0.15", "pbest/P-0.2"]

# IGD data
IGD_DATA = {}
ref_dirs = get_reference_directions("das-dennis", NUM_OBJ, n_partitions=N_VAR)

if problem.startswith("zdt"):
    pf = get_problem(problem).pareto_front()
# Dtlz-5 to dtlz-7 have their own files related to their pareto fronts
elif problem == "dtlz5" or problem == "dtlz6" or problem == "dtlz7":
   pf = get_problem(problem, n_var=N_VAR,n_obj=NUM_OBJ).pareto_front()
else:
   ref_dirs = get_reference_directions("das-dennis", 3, n_partitions=12)
   pf = get_problem(problem, n_var=N_VAR,n_obj=NUM_OBJ).pareto_front(ref_dirs=ref_dirs)

#pareto_front(n_pareto_points=ref_dirs)
metric = IGDPlus(pf=pf, normalize=True)


for variant in variants:
    execFiles = []

    for i in range(NUM_EXECS):
        filePath = base_path + variant + "/exec-" + str(i+1) + '.csv'

        execFiles.append(
            pd.read_csv(
                filePath,
                sep=';|\n', engine='python'
            )
        )


    # data of variant
    variantIGD =  np.array([])

    for i in range(GEN):
        # sum of the avarage IGD of each element
        genIGD = 0.0

        for file in execFiles:
            A = file["A"]
            B = file["B"]
            # C = file["C"]

            genData = []
            for j in range(NP*i, NP*(i+1)):
                genData.append(np.array([
                    A[j],
                    B[j],
                   # C[j]
                ]))
            genIGD = genIGD + metric.calc(np.array(genData))

        genIGD = genIGD / float(len(execFiles))
        variantIGD = np.append(variantIGD, genIGD)

    IGD_DATA[variant] = variantIGD

# write IGDdat in a separate file
df = pd.DataFrame(IGD_DATA, columns=variants)
dir_path = os.path.dirname(os.path.realpath(__file__))
df.to_csv(dir_path + "/files/" + problem + ".csv")
