import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from pymoo.factory import get_problem, get_reference_directions
from pymoo.performance_indicator.igd_plus import IGDPlus
from pymoo.factory import get_performance_indicator


# general constants
problem = "wfg1"
NUM_EXECS = 30
variants = ["rand1", "rand2", "best1", "best2", "currtobest1", "pbest/P-0.05", "pbest/P-0.1", "pbest/P-0.15", "pbest/P-0.2"] 
#  file path
base_path = "/home/nick/.gode/mode/paretoFront/" + problem + "/"
# IGD data
IGD_DATA = []

ref_dirs = get_reference_directions("das-dennis", 3, n_partitions=12)

# pf = get_problem(problem)
# metric = IGDPlus(pf=ref_point)

pf = get_problem("wfg1").pareto_front()
metric = get_performance_indicator("igd+", pf)


for varIndex in range(len(variants)):
    variantFiles = []

    for i in range(NUM_EXECS):
        filePath = base_path + variants[varIndex] + "/exec-" + str(i+1) + '.csv'
        
        variantFiles.append(
            pd.read_csv(
                filePath,
                sep='\t|\n', engine='python'
            )
        )


    # file related constants      
    GEN = int(len(variantFiles[0])/3) # GEN = QTD_LINES / QTD_OBJS
    NP = len(variantFiles[0].iloc[0]) # NP = QTD_COLS
    # NP = 100

    # data of variant
    variantIGD =  np.array([])
    for gen in range(0, GEN):              
        # sum of the avarage IGD of each element
        genIGD = 0.0
    
        for file in variantFiles:
            populationArray = []

            for i in range(NP):
                populationArray.append(np.array([
                    file.iloc[(gen*3)][i],
                    file.iloc[(gen*3)+1][i],
                    file.iloc[(gen*3)+2][i],
                ]))

            variantArray = np.array(populationArray)            
            
            
            genIGD = genIGD + metric.calc(variantArray)
                
        genIGD = genIGD / float(len(variantFiles))
        variantIGD = np.append(variantIGD, genIGD)     

    IGD_DATA.append(variantIGD)


x = np.linspace(0, int(len(variantFiles[0])/3), int(len(variantFiles[0])/3))
for i in range(len(IGD_DATA)):
    if len(x) != len(IGD_DATA[i]):
        print("variant -> " + str(i))
        print(len(IGD_DATA[i]))
    plt.scatter(x, IGD_DATA[i], s=2, alpha=0.6, label=variants[i])

plt.title(problem)
plt.suptitle("Average IGD per Generations")
plt.xlabel("generations")
plt.ylabel("IGD")
plt.legend()    
plt.show()
