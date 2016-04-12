gengen
------

Little go:generate tool for quick "generics" (using templates).

###Usage
```
Usage of gengen:
  -d string
        Directory to parse files from
  -o string
        Output filename
  -t string
        Template to render
```
Example go generate directive (taken from [gomaild2](https://github.com/trapped/gomaild2):
```
//go:generate gengen -d ./commands/ -t process.go.tmpl -o process.go
```