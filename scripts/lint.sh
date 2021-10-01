
#!/bin/bash

for f in $(find ./** -iname '*.go' -type f); do
    golines --chain-split-dots --ignore-generated -m 80 -w $f
done
