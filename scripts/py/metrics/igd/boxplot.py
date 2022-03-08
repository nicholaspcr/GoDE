import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from matplotlib.patches import Polygon

import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import os

variants = ["rand1", "rand2", "best1", "best2", "currtobest1" , "pbest/P-0.05", "pbest/P-0.1", "pbest/P-0.15", "pbest/P-0.2"]  

problems = ['zdt1', 'zdt2', 'zdt3','zdt4','zdt6', 'dtlz1', 'dtlz2', 'dtlz3', 'dtlz4', 'dtlz5', 'dtlz6','dtlz7','wfg1','wfg2','wfg3','wfg4', 'wfg5','wfg6','wfg7','wfg8','wfg9']

for problem in problems:
    dir_path = os.path.dirname(os.path.realpath(__file__))
    file_path = dir_path + '/files/' + problem + '.csv'
    IGD_DATA = pd.read_csv(file_path)
    
    data = []
    for var in variants:
        data.append(IGD_DATA[var])
    
    def set_box_color(bp, color):
        plt.setp(bp['boxes'], color=color)
        plt.setp(bp['whiskers'], color=color)
        plt.setp(bp['caps'], color=color)
        plt.setp(bp['medians'], color=color)
    
    
    fig, subs = plt.subplots(figsize=(19.20, 10.80))
    fig.subplots_adjust(left=0.075, right=0.95, top=0.9, bottom=0.25)
    
    bp = subs.boxplot(data, notch=0, sym='+', vert=1, whis=1.5, showfliers=False)
    plt.setp(bp['boxes'], color='black')
    plt.setp(bp['whiskers'], color='black')
    plt.setp(bp['fliers'], color='red', marker='+')
    
    # Add a horizontal grid to the plot, but make it very light in color
    # so we can use it for reading data values but not be distracting
    subs.yaxis.grid(True, linestyle='-', which='major', color='lightgrey',
                   alpha=0.5)
            
    subs.set(
        axisbelow=True,  # Hide the grid behind plot objects
        title='IgdPlus variants performance in ' + problem,
        xlabel='Variants',
        ylabel='IgdPlus',
    )
    
    # Now fill the boxes with desired colors
    box_colors = ['darkkhaki', 'royalblue']
    num_boxes = len(data)
    medians = np.empty(num_boxes)
    for i in range(num_boxes):
        box = bp['boxes'][i]
        box_x = []
        box_y = []
        for j in range(5):
            box_x.append(box.get_xdata()[j])
            box_y.append(box.get_ydata()[j])
        box_coords = np.column_stack([box_x, box_y])
        
        # Alternate between Dark Khaki and Royal Blue
        subs.add_patch(Polygon(box_coords, facecolor=box_colors[i % 2]))
    
    
        # Now draw the median lines back over what we just filled in
        med = bp['medians'][i]
        median_x = []
        median_y = []
        for j in range(2):
            median_x.append(med.get_xdata()[j])
            median_y.append(med.get_ydata()[j])
            subs.plot(median_x, median_y, 'k')
        medians[i] = median_y[0]
        # Finally, overplot the sample averages, with horizontal alignment
        # in the center of each box
        # subs.plot(np.average(med.get_xdata()), np.average(data[i]),
        #          color='w', marker='*', markeredgecolor='k')
    
    
    x_indices = []
    for i in range(len(variants)):
        x_indices.append(i+1)
    plt.xticks(x_indices,variants)
    
    plt.tick_params(axis='x', labelsize=20)
    plt.tick_params(axis='y', labelsize=20)
    
    plt.title(
            'IGD+ variants performance in ' + problem.upper(),
            fontdict={'fontsize': 20, 'fontweight': 'medium'}
            )
    plt.xlabel('Variantes', fontsize=20)
    plt.ylabel('IGD+', fontsize=20)
    
    
    plt.legend()
    # plt.show()
    plt.savefig('/home/nick/Documents/ic/imgs/igd+/generated/igdplus_' + problem + '.svg', format='svg',
            dpi=1200)
    plt.close()    # close the figure window

