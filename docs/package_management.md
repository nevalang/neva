# Package management

Package is something that can be download and/or run.
Is a set of modules and `pkg.yml` file called _package descriptor_.
It has structure:

```yaml
deps: # this section maps git repos and tags with aliases
  shared: # now `shared` can be used in `import` section
    repo: "github.com/emil14/respect-shared" # assumes there is a pkg.yml
    v: 0.0.2 # with this tag

# this section defines packages's scope
import:
  # local import, starts with the './'
  me: ./mul_by_square
  # global import - always have two parts splitted by '/' delimeter
  sqrt: shared/square # assumes there is `square` export in the pkg.yml in `shared` package

# this section allowes reuse modules from this package
export:
  - me
  - sqrt

# if pkg.yml has this - the whole package can be used as one big module
root: me
```
