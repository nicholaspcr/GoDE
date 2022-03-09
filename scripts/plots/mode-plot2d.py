import pandas as pd
import matplotlib.pyplot as plt

problems = ['zdt1','zdt2','zdt3','zdt4','zdt6']


for problem in problems:

    variant = 'rand1'
    
    filePath = '/home/nick/.gode/mode/multiExecutions/{}/{}/rankedPareto.csv'.format(problem, variant)
    
    f = pd.read_csv(
        filePath,
        sep=';|\n', engine='python'
    )
    
    number_rows, number_cols = f.shape
    
    # number_cols = number_cols - 1
    X = f["A"]
    Y = f["B"]
    
    fig = plt.figure(figsize=(19.20, 10.80))
    ax = fig.add_subplot()
    
    plt.title(problem, fontsize=20)
    
    ax.tick_params(axis='both', labelsize=16, rotation=45)
    ax.set_xlabel("obj-1", fontsize=20)
    ax.set_ylabel("obj-2", fontsize=20)
    ax.scatter(X, Y, s=20)
    
    plt.legend()
    # plt.show()
    plt.savefig('/home/nick/Documents/ic/imgs/plots/generated/plot_' + problem + '.svg', format='svg', dpi=1200)
    plt.close()    # close the figure window
    
