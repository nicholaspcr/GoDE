import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from matplotlib.patches import Polygon

import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import os
import sys

# problem = "dtlz2"
# variants = ["rand1", "rand2", "best1", "best2", "currtobest1" , "pbest/P-0.05", "pbest/P-0.1", "pbest/P-0.15", "pbest/P-0.2"]  

# GEN = 151

# original_array = np.loadtxt("files/" + problem + ".txt").reshape(len(variants), GEN)
# print(original_array)


def set_box_color(bp, color):
    plt.setp(bp['boxes'], color=color)
    plt.setp(bp['whiskers'], color=color)
    plt.setp(bp['caps'], color=color)
    plt.setp(bp['medians'], color=color)

variants = ["rand1", "rand2", "best1", "best2", "currtobest1", "pbest/P-0.05", "pbest/P-0.1", "pbest/P-0.15", "pbest/P-0.2"]

problem = "dtlz1"
filePath = "/home/nick/.gode/mode/multiExecutions/" + \
   problem + "/rand1/rankedPareto.csv"

f = pd.read_csv(
    filePath,
    sep='\t|\n', engine='python'
)

number_rows, number_cols = f.shape
# print(number_rows, number_cols)

data = [
    f["A"],
    f["B"],
    f["C"]
]

fig, subs = plt.subplots(figsize=(10, 6))
fig.subplots_adjust(left=0.075, right=0.95, top=0.9, bottom=0.25)

bp = subs.boxplot(data, notch=0, sym='+', vert=1, whis=1.5)
plt.setp(bp['boxes'], color='black')
plt.setp(bp['whiskers'], color='black')
plt.setp(bp['fliers'], color='red', marker='+')

# Add a horizontal grid to the plot, but make it very light in color
# so we can use it for reading data values but not be distracting
subs.yaxis.grid(True, linestyle='-', which='major', color='lightgrey',
               alpha=0.5)

subs.set(
    axisbelow=True,  # Hide the grid behind plot objects
    title='Comparison of IID Bootstrap Resampling Across Five Distributions',
    xlabel='Distribution',
    ylabel='Value',
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

plt.show()
