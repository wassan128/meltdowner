# MeltDowner
* convert markdown to html
* future work: blog manager for me

### How to build
* `git clone git@github.com:wassan128/meltdowner.git`
* `go build meltdowner/main.go && mv -f main melt`

### How to use
#### Initialize
`./melt init`

#### Create new post
`./melt new "POST TITLE"`

#### Write post content
`vim source/YYYYMMDD_POST_TITLE.md`

#### Start server
`./melt server`

#### Generate static content
`./melt generate`

### License
MIT

### Author
wassan128

