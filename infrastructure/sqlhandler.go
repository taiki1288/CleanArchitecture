package infrastructure

// DBは外部パッケージを使うので一番外側のinfrastructure層になる。
import (
	"Clean-Architecture/interfaces/database"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)


type SqlHandler struct {
	// *Sql.DBはdatabase/sqlパッケージのDB接続に必要なtype
	Conn *sql.DB
}

// NewSqlHandlerの戻り値はinterfaceのSqlHandler
func NewSqlHandler() database.SqlHandler {
	conn, err := sql.Open("mysql", "root:@/CleanArchitecture")
	if err != nil {
		panic(err.Error)
	}
	// newは指定した方のポインタ型を生成する関数。(構造体の初期化)
	sqlHandler := new(SqlHandler)
	sqlHandler.Conn = conn
	return sqlHandler
}

func (handler *SqlHandler) Execute(statement string, args ...interface{}) (database.Result, error) {
	res := SqlResult{}
	// Execは、行を返さずにクエリを実行する。
	// 引数は(クエリ, クエリ内のプレースホルダパラメータ)
	result, err := handler.Conn.Exec(statement, args...)
	if err != nil {
		return res, err
	}
	res.Result = result
	return res, nil
}

func (handler *SqlHandler) Query(statement string, args ...interface{}) (database.Row, error) {
	// Queryは、行を返すクエリ(通常は SELECT)を実行する。
	// 引数は(クエリ, クエリ内のプレースホルダパラメータ)
	rows, err := handler.Conn.Query(statement, args...)
	if err != nil {
		return new(SqlRow), err
	}
	// newは指定した方のポインタ型を生成する関数。(構造体の初期化)
	row := new(SqlRow)
	row.Rows = rows
	return row, nil
}

// LastInsertId, RowsAffectedメソッドの呼び出しに必要な構造体
type SqlResult struct {
	// Result型
	Result sql.Result
}

func (r SqlResult) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

func (r SqlResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

// Scan, Next, Closeメソッドの呼び出しに必要な構造体
type SqlRow struct {
	// Rows型
	Rows *sql.Rows
}

func (r SqlRow) Scan(dest ...interface{}) error {
	// 現在の行の列をdestが指す値にコピーする。
	return r.Rows.Scan(dest...)
}

func (r SqlRow) Next() bool {
	// Scan メソッドで読み取るため，次の結果行を準備します。成功した場合は trueになる。
	return r.Rows.Next()
}

func (r SqlRow) Close() error {
	// Rowsを閉じ、それ以上の列挙はできなくなる。
	return r.Rows.Close()
}