package generator

var TPL_REPO_FACTORY = CodeTemplate(`
        package repo
        {{$var := .VAR}}
		import(
		    "github.com/jsix/gof/db/orm"
		    "{{$var.ModelPkg}}"
		    "{{$var.RepoPkg}}"
		    "{{$var.IRepoPkg}}"
		)

		type repoFactory struct{
			o orm.Orm
		{{range $i,$tb := .Tables}}
		    _{{$tb.Name}}_repo {{$var.IRepoPkgName}}.I{{$tb.Title}}Repo{{end}}
		}

		func NewRepoFactory(o orm.Orm)*repoFactory{
			r := &repoFactory{
				o:o,
			}
			return r.init()
		}
		func (r *repoFactory) init()*repoFactory{
        {{range $i,$tb := .Tables}}
		    r.o.Mapping({{$var.ModelPkgName}}.{{$tb.Title}}{},"{{$tb.Name}}"){{end}}
		    return r
		}
		{{range $i,$tb := .Tables}}
		func (r *repoFactory) Get{{$tb.Title}}Repo(){{$var.IRepoPkgName}}.I{{$tb.Title}}Repo{
		    if r._{{$tb.Name}}_repo == nil{
		        r._{{$tb.Name}}_repo = {{$var.RepoPkgName}}.New{{$tb.Title}}Repo(r.o)
		    }
		    return r._{{$tb.Name}}_repo
		}
		{{end}}
`)
