package tcpServer

import (
	"conf"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type MysqlClient struct {
	DBSql      *sql.DB
	DBName     string
	TableName  string
	Password   string
	UserName   string
	StmtIns    *sql.Stmt
	StmtSearch *sql.Stmt
	StmtUpdate *sql.Stmt
}

type Col struct {
	AutoID   string
	Account  string
	Password string
	Nickname string
	Extend   string
}

var MC *MysqlClient

func init() {
	MC = &MysqlClient{
		DBSql:     nil,
		DBName:    conf.MysqlDB,
		TableName: conf.MysqlTable,
		UserName:  conf.MysqlUser,
		Password:  conf.MysqlPass,
	}

	var err error
	MC.DBSql, err = sql.Open("mysql", MC.UserName+":"+MC.Password+"@/"+MC.DBName)
	if err != nil {
		panic(err.Error())
	}

	MC.StmtIns, err = MC.DBSql.Prepare("insert into " + MC.TableName + "(account,password,nickname) values(? ,? ,?)")
	if err != nil {
		panic(err)
	}

	MC.StmtSearch, err = MC.DBSql.Prepare("select * from " + MC.TableName + " where account = ?")
	if err != nil {
		panic(err)
	}

	MC.StmtUpdate, err = MC.DBSql.Prepare("update " + MC.TableName + " set nickname = ? where account = ?")
	if err != nil {
		panic(err)
	}

	err = MC.DBSql.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func conCheck() {
	if MC.DBSql == nil || time.Now().Second()%10 == 0 {
		err := MC.DBSql.Ping()
		if err != nil {
			panic(err)
		}
	}
}

func (c *MysqlClient) InsertUser(col *Col) error {
	_, err := c.StmtIns.Exec(col.Account, col.Nickname, col.Password)
	return err
}

func (c *MysqlClient) GetUser(account string) *Col {
	row := c.StmtSearch.QueryRow(account)
	col := &Col{}
	row.Scan(&col.AutoID, &col.Account, &col.Password, &col.Nickname, &col.Extend)

	if col.AutoID != "" {
		return col
	}
	return nil
}

func (c *MysqlClient) UpdateNickname(changeValue string, account string) error {
	_, err := c.StmtUpdate.Exec(changeValue, account)
	return err
}
