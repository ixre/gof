package db

type (
	// 表
	Table struct {
		// 名称
		Name string
		// 注释
		Comment string
		// 引擎
		Engine string
		// 架构
		Schema string
		// 编码
		Charset string
		// 列
		Columns []*Column
	}

	// 列
	Column struct {
		Name    string
		IsPk    bool
		IsAuto  bool
		NotNull bool
		DbType  string
		Length  int
		Comment string
		Type    int
	}
)



var (
	//TypeUnknown = int(reflect.Invalid)
	//TypeString  = int(reflect.String)
	//TypeBoolean = int(reflect.Bool)
	//TypeInt16   = int(reflect.Int16)
	//TypeInt32   = int(reflect.Int32)
	//TypeInt64   = int(reflect.Int64)
	//TypeFloat32 = int(reflect.Float32)
	//TypeFloat64 = int(reflect.Float64)

	TypeUnknown = 0
	TypeString  = 1
	TypeBoolean = 2
	TypeInt16   = 3
	TypeInt32   = 4
	TypeInt64   = 5
	TypeFloat32 = 6
	TypeFloat64 = 7

	TypeBytes    = 14
	TypeDateTime = 15
	TypeDecimal  = 16
)