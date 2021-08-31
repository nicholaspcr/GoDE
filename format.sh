PASS=""

DIRS=( )



for file in $(ls -la | grep "\.go\$")
do
    # formats the file
    golines $files
    if [[ $? -ne 0 ]]; then
		PASS="golines error"
	fi
done

if [ ! -z "$PASS" ]; then
    echo "COMMIT FAILED - $PASS"
    exit 1
fi





