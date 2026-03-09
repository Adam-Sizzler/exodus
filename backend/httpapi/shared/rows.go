package shared

// RowScanner allows using a shared scan function for sql.Rows and sql.Row.
type RowScanner interface {
	Scan(dest ...interface{}) error
}
