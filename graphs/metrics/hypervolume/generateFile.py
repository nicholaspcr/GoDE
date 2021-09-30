import numpy as np
import pandas as pd
from pymoo.factory import get_problem, get_reference_directions 
from pymoo.performance_indicator.hv import Hypervolume
import os

variants = ["rand1", "rand2", "best1", "best2", "currtobest1", "pbest/P-0.05", "pbest/P-0.1", "pbest/P-0.15", "pbest/P-0.2" ]

# general constants
NUM_EXECS = 30
NUM_OBJ = 2

N_VAR = 30
N_OBJ = 2

# file related constants      
GEN = 100
NP = 100

# important methods -> ref_point
problem = "zdt6"
ref_point = np.array([1, 9.735527117321219])
   
if problem.startswith("zdt"):
    pf = get_problem(problem).pareto_front()
# Dtlz-5 to dtlz-7 have their own files related to their pareto fronts
elif problem == "dtlz5" or problem == "dtlz6" or problem == "dtlz7":
   pf = get_problem(problem, n_var=N_VAR,n_obj=N_OBJ).pareto_front()
else:
   ref_dirs = get_reference_directions("das-dennis", 3, n_partitions=12)
   pf = get_problem(problem, n_var=N_VAR,n_obj=N_OBJ).pareto_front(ref_dirs=ref_dirs)


metric = Hypervolume(pf=pf, ref_point=ref_point, normalize=True)

# dtlz-1 -> [482.1632866685836,479.911835933397,471.9731868733398]
# dtlz-2 -> [2.651414902010465,2.5206368624614965,2.656093434231162]
# dtlz-3 -> [1784.9822112456513,1683.7871520696372,1679.1459524987113]
# dtlz-4 -> [2.7493608245409247,2.665459302333755,2.691506519652278]
# dtlz-5 -> [2.6184046195044153,2.3154562025982375,2.490037232873547]
# dtlz-6 -> [10.460414515081052,10.523716498291654,10.571261523682367]
# dtlz-7 -> [1,1,24.464595045398383]

# zdt-1 -> [1, 5.687041127771669]
# zdt-2 -> [1, 6.71194298397789]
# zdt-3 -> [1, 6.020951819554247]
# zdt-4 -> [1, 129.8511197453462]
# zdt-6 -> [1, 9.735527117321219]

# WFGS
# WFG-1 -> [2.9699152058577947,4.985423795813952,6.988314095918014]
# WFG-2 -> [2.733039177657597,4.647396034258947,6.597212425291411]
# WFG-3 -> [2.9788962857725294,4.14440451610059,6.66644969309052]
# WFG-4 -> [2.571914973474083,4.552928861224317,6.476940203979565]
# WFG-5 -> [2.614787579010435,4.621198824818527,6.660416531655121]
# WFG-6 -> [2.823975776872104,4.811698190000338,6.839291818048697]
# WFG-7 -> [2.646804889830734,2.646804889830734,6.624496532156823]
# WFG-8 -> [2.9828456406010377,4.8457906153983545,6.74641836458321]
# WFG-9 -> [2.879043535029743,4.887602113615263,6.908800119495598]

base_path = "/home/nick/.gode/mode/paretoFront/" + problem + "/"
# data for the plot
HVData = {}

for variant in variants:
    execFiles = []


    # reads from each execution of the variant
    for i in range(NUM_EXECS):
        filePath = base_path + variant + "/exec-" + str(i+1) + '.csv'
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
            #C = file["C"]
            genData = []
            for j in range(NP*i, NP*(i+1)):
                genData.append(np.array([
                    A[j],
                    B[j],
                    # C[j]
                ]))
            genHV = genHV + metric.calc(np.array(genData))

        genHV = genHV / float(len(execFiles))
        variantHV = np.append(variantHV, genHV)

    HVData[variant] = variantHV

# write HVData in a separate file
df = pd.DataFrame(HVData, columns=variants)
dir_path = os.path.dirname(os.path.realpath(__file__))
df.to_csv(dir_path + "/files/" + problem + ".csv")
