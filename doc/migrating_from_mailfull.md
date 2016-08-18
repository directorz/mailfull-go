Migrating from mailfull
=======================

Migrating from [directorz/mailfull](https://github.com/directorz/mailfull)

Change directory to Mailfull directory.

```
# su - mailfull
$ cd /home/mailfull
```

Initialize a directory as a Mailfull repository.

```
$ mailfull init
```

Delete unnecessary files.

```
$ rm -rf .git .gitignore bin docs lib README.md README.ja.md
$ find domains -maxdepth 2 -name '.vforward' | xargs rm -f
```
