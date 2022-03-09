#!/bin/bash

ZDTS=("zdt1" "zdt2" "zdt3" "zdt4" "zdt6")
VARIANTS=("rand1" "rand2" "best1" "best2" "currtobest1" "pbest")
PBEST_VALUES=("0.05" "0.1" "0.15" "0.2")

for problem in "${ZDTS[@]}"
do
    echo $problem    
    for variant in "${VARIANTS[@]}"
    do 
        echo $variant
        if [ "$variant" = "pbest" ]; then
            for pvalue in "${PBEST_VALUES[@]}"
            do
                
            done
        else

        fi      
    done
done
