package orm

import "strings"

type GenTemplate string

var (
	// <R> : 仓储类的名称
	// <R2> : 函数后添加的仓储类名称
	// <E> : 实体
	// <T> : 仓库类对象引用
	TPL_ENTITY_REP = GenTemplate(
		`// auto generate by gof (http://github.com/jsix/gof)
        namespace rep
        import(
            "database/sql"
            "github.com/jsix/gof/db/orm"
        )
        type <R> struct{
            _orm orm.ORM
        }

        // Create new <R>
        func New<R>(o orm.ORM)*<R>{
            return &<R>{
                _orm:o,
            }
        }

        // Get <E>
        func (<T> *<R>) Get<R2>(primary interface{})*<E>{
            e := <E>{}
            err := <T>._orm.Get(primary,&e)
            if err == nil{
                return &e
            }
            if err != sql.ErrNoRow{
              log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
            }
            return nil
        }

        // Save <E>
        func (<T> *<R>) Save<R2>(v *<E>)(int,error){
            id,err := orm.Save(<T>._orm,v,v.<PK>)
            if err != nil && err != sql.ErrNoRow{
              log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
            }
            return id,err
        }

        // Delete <E>
        func (<T> *<R>) Delete<R2>(primary interface{}) error {
            err := <T>._orm.DeleteByPk(<E>{}, primary)
            if err != nil && err != sql.ErrNoRow{
              log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
            }
            return err
        }

        // Select <E>
        func (<T> *<R>) Select<R2>(where string,v ...interface{})[]*<E> {
            list := []*<E>{}
            err := <T>._orm.Select(&list,where,v...)
            if err != nil && err != sql.ErrNoRow{
              log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
            }
            return list
        }
        `)
)

func (g GenTemplate) String() string {
	return string(g)
}

func (g GenTemplate) Replace(s, r string, n int) GenTemplate {
	return GenTemplate(strings.Replace(string(g), s, r, n))
}

func init() {
	TPL_ENTITY_REP = TPL_ENTITY_REP.Replace("<T>", "{{.T}}", -1).
		Replace("<E>", "{{.E}}", -1).
		Replace("<R>", "{{.R}}", -1).
		Replace("<R2>", "{{.R2}}", -1).
		Replace("<PK>", "{{.PK}}", -1)

}
