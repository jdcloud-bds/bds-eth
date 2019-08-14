# Upgrade manual

To apply the patch set we made for earlier version, for example we want to
upgrade from v1.8.20 to v1.8.21.

First, create a new branch based on the upstream release.

```
git clone git@git.jd.com:jbri/go-ethereum.git
cd go-ethereum
git remote add upstream https://github.com/ethereum/go-ethereum
git fetch upstream

# create our own branch v1.8.21-jbri based on the release tag
git checkout -b v1.8.21-jbri refs/tags/v1.8.21

git push origin refs/tags/v1.8.21  # make sure tags available on our server
git push origin v1.8.21-jbri       # make sure new branch available on our server
```
> If some file exceed 10M, please clean source tree by
> [bfg-repo-cleaner](https://rtyley.github.io/bfg-repo-cleaner/).

Then get the patch commit list we made for v1.8.20:

```
git rev-list v1.8.20..v1.8.20-jbri --no-merges --reverse > patch.list
```

Now cherry-pick the commits one by one with following script, you may need to
resolve conflicts manually and safely rerun the script to continue cherry-pick
the remaining commits:

```bash
PATCH_LIST=patch.list

for x in $(cat $PATCH_LIST); do
    if git cherry-pick $x; then
        sed -i "/$x/d" $PATCH_LIST
        echo "[*] Successfully cherry-pick'ed $x (removed from $PATCH_LIST)"
    else
        echo "[*] Failed to cherry-pick $x"
        echo "[*] Please manual resolve conflict and remove first line in $PATCH_LIST"
        break
    fi
done
```

Now we can push branch to origin repo.

> Do not commit `patch.list` file.