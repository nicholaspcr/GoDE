import pandas as pd
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
import numpy as np

problems=['dtlz1','dtlz2','dtlz3','dtlz4','dtlz5','dtlz6','dtlz7','wfg1','wfg2','wfg3','wfg4','wfg5','wfg6','wfg7','wfg8','wfg9']

for problem in problems:
    variant = "rand1"
    
    filePath = "/home/nick/.gode/mode/multiExecutions/" + \
       problem + "/" + variant + "/rankedPareto.csv"
    
    
    f = pd.read_csv(
        filePath,
        sep=';|\n', engine='python'
    )
    
    number_rows, number_cols = f.shape
    
    X = f["A"]
    Y = f["B"]
    Z = f["C"]
    
    fig = plt.figure(figsize=(19.20, 10.80))
    ax = Axes3D(fig)
    # ax.xticks(fontsize=14)
    
    ax.set_title(
            problem.upper(),
            fontdict={'fontsize': 20, 'fontweight': 'medium'}
    )
    
    ax.tick_params(axis='both', labelsize=14, rotation=45)
    
    ax.set_xlabel("obj-1", fontsize=20)
    ax.set_ylabel("obj-2", fontsize=20)
    ax.set_zlabel("obj-3", fontsize=20)
    ax.scatter(X, Y, Z, s=20)
    
    # for rotate the axes and update.
    ax.view_init(30,60)
    
    
    # ax.set_yticklabels(ax.get_yticks(), rotation = 45)
    plt.legend()
    # plt.show()
    plt.savefig('/home/nick/Documents/ic/imgs/plots/generated/plot_' + problem + '.png',
    format='svg', dpi=300)
    plt.close()    # close the figure window
    
