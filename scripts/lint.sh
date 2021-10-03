
#!/bin/bash

for f in $(find ./** -iname '*.go' -type f); do
    golines --chain-split-dots --ignore-generated --reformat-tags --shorten-comments -m 80 -w $f
done
