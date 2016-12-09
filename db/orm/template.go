package orm

import "strings"

type CodeTemplate string

var (
	// 实体仓储接口模板
	TPL_ENTITY_REP_INTERFACE = CodeTemplate(
		`// auto generate by gof (http://github.com/jsix/gof)
        package {{.VAR.IRepoPkgName}}

        import(
            "{{.VAR.ModelPkgIRepo}}"
        )

        type I<R> interface{
            // Get <E>
            Get<R2>(primary interface{})*<E2>
            // Select <E>
            Select<R2>(where string,v ...interface{})[]*<E2>
            // Save <E>
            Save<R2>(v *<E2>)(int,error)
            // Delete <E>
            Delete<R2>(primary interface{}) error
            // Batch Delete <E>
            BatchDelete<R2>(where string,v ...interface{})(int64,error)
        }`)

	// <R> : 仓储类的名称
	// <R2> : 函数后添加的仓储类名称
	// <E> : 实体
	// <E2> : 包含包名的实体
	// <Ptr> : 仓库类对象引用
	TPL_ENTITY_REP = CodeTemplate(
		`// auto generate by gof (http://github.com/jsix/gof)
            package {{.VAR.RepoPkgName}}
            import(
                "log"
                "{{.VAR.ModelPkg}}"
                "database/sql"
                "github.com/jsix/gof/db/orm"
            )

            type <R> struct{
                _orm orm.Orm
            }

            // Create new <R>
            func New<R>(o orm.Orm)*<R>{
                return &<R>{
                    _orm:o,
                }
            }

            // Get <E>
            func (<Ptr> *<R>) Get<R2>(primary interface{})*<E2>{
                e := <E2>{}
                err := <Ptr>._orm.Get(primary,&e)
                if err == nil{
                    return &e
                }
                if err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return nil
            }

            // Select <E>
            func (<Ptr> *<R>) Select<R2>(where string,v ...interface{})[]*<E2> {
                list := []*<E2>{}
                err := <Ptr>._orm.Select(&list,where,v...)
                if err != nil && err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return list
            }

            // Save <E>
            func (<Ptr> *<R>) Save<R2>(v *<E2>)(int,error){
                id,err := orm.Save(<Ptr>._orm,v,v.<PK>)
                if err != nil && err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return id,err
            }

            // Delete <E>
            func (<Ptr> *<R>) Delete<R2>(primary interface{}) error {
                err := <Ptr>._orm.DeleteByPk(<E2>{}, primary)
                if err != nil && err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return err
            }

			// Batch Delete <E>
            func (<Ptr> *<R>) BatchDelete<R2>(where string,v ...interface{})(int64,error) {
                r,err := <Ptr>._orm.Delete(<E2>{},where,v...)
                if err != nil && err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return r,err
            }

            `)
)

func (g CodeTemplate) String() string {
	return string(g)
}

func (g CodeTemplate) Replace(s, r string, n int) CodeTemplate {
	return CodeTemplate(strings.Replace(string(g), s, r, n))
}

func resolveRepTag(g CodeTemplate) CodeTemplate {
	return g.Replace("<Ptr>", "{{.Ptr}}", -1).
		Replace("<E>", "{{.E}}", -1).
		Replace("<E2>", "{{.E2}}", -1).
		Replace("<R>", "{{.R}}", -1).
		Replace("<R2>", "{{.R2}}", -1).
		Replace("<PK>", "{{.PK}}", -1)
}

func init() {
	TPL_ENTITY_REP = resolveRepTag(TPL_ENTITY_REP)
	TPL_ENTITY_REP_INTERFACE = resolveRepTag(TPL_ENTITY_REP_INTERFACE)
}
