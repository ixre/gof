package generator

var (
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

            // Create new <R>
            func New<R>(o orm.Orm)*<R>{
                return &<R>{
                    _orm:o,
                }
            }

            type <R> struct{
                _orm orm.Orm
            }

            // Get <E>
            func (<Ptr> *<R>) Get<R2>_(primary interface{})*<E2>{
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

            // GetBy <E>
            func (<Ptr> *<R>) Get<R2>By_(where string,v ...interface{})*<E2>{
                e := <E2>{}
                err := <Ptr>._orm.GetBy(&e,where,v...)
                if err == nil{
                    return &e
                }
                if err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return nil
            }

            // Select <E>
            func (<Ptr> *<R>) Select<R2>_(where string,v ...interface{})[]*<E2> {
                list := []*<E2>{}
                err := <Ptr>._orm.Select(&list,where,v...)
                if err != nil && err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return list
            }

            // Save <E>
            func (<Ptr> *<R>) Save<R2>_(v *<E2>)(int,error){
                id,err := orm.Save(<Ptr>._orm,v,int(v.<PK>))
                if err != nil && err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return id,err
            }

            // Delete <E>
            func (<Ptr> *<R>) Delete<R2>_(primary interface{}) error {
                err := <Ptr>._orm.DeleteByPk(<E2>{}, primary)
                if err != nil && err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return err
            }

            // Batch Delete <E>
            func (<Ptr> *<R>) BatchDelete<R2>_(where string,v ...interface{})(int64,error) {
                r,err := <Ptr>._orm.Delete(<E2>{},where,v...)
                if err != nil && err != sql.ErrNoRows{
                  log.Println("[ Orm][ Error]:",err.Error(),"; Entity:<E>")
                }
                return r,err
            }

            `)
)
