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
Example go generate directive (taken from [gomaild2](https://github.com/trapped/gomaild2)):
```
//go:generate gengen -d ./commands/ -t process.go.tmpl -o process.go
```

`gengen` parses a directory (`-d`) for Go code, then renders your template (`-t`) exposing the parsed data as `.Packages` (`map[name string]*ast.Package`).
This helps reduce "pseudo-generic" boilerplate code, such as monolithic `switch` statements (see [gomaild2](https://github.com/trapped/gomaild2)'s [SMTP command processor](https://github.com/trapped/gomaild2/blob/master/smtp/process.go.tmpl)).