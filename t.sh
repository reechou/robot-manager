# maybe more powerful
# for mac (sed for linux is different)
dir=`echo ${PWD##*/}`
grep "robot-manager" * -R | grep -v Godeps | awk -F: '{print $1}' | sort | uniq | xargs sed -i '' "s#robot-manager#$dir#g"
mv robot-manager.ini $dir.ini

