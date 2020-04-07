package DbClient

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/gentwolf-shen/gohelper/httpclient"
)

type DbDriver struct {
	dbServer string
	auth     *Auth
}

func NewDbDriver(dbServer, appKey, appSecret string) *DbDriver {
	dbDriver := &DbDriver{}
	dbDriver.dbServer = dbServer

	dbDriver.auth = &Auth{}
	dbDriver.auth.AppKey = appKey
	dbDriver.auth.AppSecret = appSecret

	return dbDriver
}

/**
查询数据，返回多条记录
*/
func (this DbDriver) Query(item *SqlMessage) ([]map[string]string, error) {
	b, err := this.send("query", item)
	if err != nil {
		return nil, err
	}

	var rows []map[string]string
	err = json.Unmarshal(b, &rows)

	return rows, err
}

/**
查询数据，单条记录
*/
func (this DbDriver) QueryRow(item *SqlMessage) (map[string]string, error) {
	rows, err := this.Query(item)
	if err != nil {
		return nil, err
	}

	if len(rows) > 0 {
		return rows[0], nil
	}

	return nil, nil
}

/**
查询记录，一个字段
*/
func (this DbDriver) QueryScalar(name string, item *SqlMessage) (string, error) {
	row, err := this.QueryRow(item)
	if err != nil {
		return "", err
	}

	if len(row) > 0 {
		return row[name], nil
	}

	return "", nil
}

/**
更新数据
*/
func (this DbDriver) Update(item *SqlMessage) (int64, error) {
	b, err := this.send("update", item)
	if err != nil {
		return 0, err
	}

	n, _ := strconv.ParseInt(strings.TrimSpace(string(b)), 10, 64)
	return n, nil
}

/**
添加数据
*/
func (this DbDriver) Insert(item *SqlMessage) (int64, error) {
	return this.Update(item)
}

/**
删除数据
*/
func (this DbDriver) Delete(item *SqlMessage) (int64, error) {
	return this.Update(item)
}

/**
事务处理，多条SQL必须是操作同一数据库
SQL: UPDATE、INSERT、DELETE
*/
func (this DbDriver) Transaction(items []*SqlMessage) (bool, error) {
	if _, err := this.send("transaction", items); err != nil {
		return false, err
	}

	return true, nil
}

/**
事务处理，多条SQL必须是操作同一数据库
SQL: UPDATE、INSERT、DELETE
*/
func (this DbDriver) TransactionV2(items *BatchSqlMessage) (bool, error) {
	if _, err := this.send("v2/transaction", items); err != nil {
		return false, err
	}

	return true, nil
}

/**
批量查询
*/
func (this DbDriver) BatchQuery(items []*SqlMessage) ([][]map[string]string, error) {
	b, err := this.send("batch/query", items)
	if err != nil {
		return nil, err
	}

	var rows [][]map[string]string
	err = json.Unmarshal(b, &rows)

	return rows, err
}

/**
发送数据操作命令
*/
func (this DbDriver) send(method string, item interface{}) ([]byte, error) {
	headers := make(map[string]string, 1)
	sql := ""

	switch item.(type) {
	case *SqlMessage:
		sql = item.(*SqlMessage).Sql
		item.(*SqlMessage).Token = this.auth.GetAuthToken(sql)
	case *BatchSqlMessage:
		sql = item.(*BatchSqlMessage).Sql
		item.(*BatchSqlMessage).Token = this.auth.GetAuthToken(sql)
	case []*SqlMessage:
		sql = item.([]*SqlMessage)[0].Sql
		item.([]*SqlMessage)[0].Token = this.auth.GetAuthToken(sql)
	}

	data, _ := json.Marshal(item)
	b, err := httpclient.PostToBody(this.dbServer+"/"+method, data, headers)
	if err != nil {
		str := sql + "\n\thttp code " + err.Error() + "\n\tmsg " + string(b)
		return b, errors.New(str)
	}

	return b, nil
}
