# MeltDowner
* convert markdown to html
* simple static site generator

## How to build
* `git clone git@github.com:wassan128/meltdowner.git`
* `go build -o melt meltdowner/main.go`

## How to use
### Initialize
`./melt init`

### Create new post
`./melt new "POST TITLE"`

### Write post content
`vim source/YYYYMMDD_POST_TITLE.md`

### Start server
`./melt server`

### Generate static content
`./melt generate`

## License
MIT

## Author
wassan128

