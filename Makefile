#  make gt a=v0.0.0 m=标签
a = test0.0.1
m = 标签
NID = `git log --pretty=%H -1`
gt:
	git tag -a $(a) $(NID) -m '$(m)'; git push -u origin $(a)
