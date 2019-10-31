# extract-git-objects

A tool for extracting files from the Git object store onto the filesystem.

## Usage

```sh
go install github.com/Artmann/extract-git-objects

extract-git-objects
```

This will create a files structure similr to this

```
- objects/24/315ebe8c3cfdefa6542ebe8cfea6d227795d44.go
- objects/35/baf22e2d0f8adccf77bd1166bd9c4e85874b40.yml
- objects/64/58a1db269593c0b300ddf46597d7dffd08cbf1.go
- objects/6b/e02374db118b9fb99cd98c6b403e5a558d0d57.js
- objects/8f/0b8aebd862ae6b2787cde7bd3ea36c5a441c10.go
- objects/a3/b3956fef0c927f670e68c98e5e676e536235ec.go
- objects/c1/54cc97f7a1765c7772e589658467351c559d05.sh
- objects/ca/2bc394db8a3448881c854fc92da510a178db81.go
- objects/d6/d1aa461633fd97e4ef94ecd52ee506d7ee9e97
- objects/d8/e1482034747db5f3e021bca4c92ce6c22e750e.sh
- objects/da/78ea3a37914c2bfe781861536eace852e32394.go
- objects/df/01c9f51dc4185f05a4a51d3dea9850deadfb7a.gitignore
- objects/f1/9b564ece400dd588541199b0c6b8477376801c.go
```
