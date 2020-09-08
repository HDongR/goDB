package models

import (
	"database/sql"
	"errors"
	"fmt"
	"goDB/config"
	"log"
	"strings"
	"time"
)

type User struct {
	Id      int64       `form:"id"`
	Loginid string      `form:"loginid"`
	Passwd  string      `form:"passwd"`
	Name    string      `form:"name"`
	Date    string      `form:"date"`
	Extra   interface{} `form:"extra"`
}

type UserManager struct {
	Conn   *sql.DB
	Result *sql.Result
	Prefix string
	Index  string
}

func NewUserManager(conn *sql.DB) *UserManager {
	var item UserManager

	if conn == nil {
		item.Conn = NewConnection()
	} else {
		item.Conn = conn
	}

	item.Prefix = "u"
	item.Index = ""

	return &item
}

func (p *UserManager) Close() {
	if p.Conn != nil {
		p.Conn.Close()
	}
}

func (p *UserManager) GetLast(items *[]User) *User {
	if items == nil {
		return nil
	} else if len(*items) == 0 {
		return nil
	} else {
		return &(*items)[0]
	}
}

func (p *UserManager) SetIndex(index string) {
	p.Index = index
}

func (p *UserManager) GetQuery() string {
	ret := ""

	tableName := "user_tb"
	if config.Database == "mssql" || config.Database == "sqlserver" {
		tableName = config.Owner + ".user_tb"
	}

	str := "select u_id, u_loginid, u_passwd, u_name, u_date from " + tableName + " "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ") "
	}

	return ret
}

func (p *UserManager) GetQuerySelect() string {
	ret := ""

	tableName := "user_tb"
	if config.Database == "mssql" || config.Database == "sqlserver" {
		tableName = config.Owner + ".user_tb"
	}

	str := "select count(*) from " + tableName + " "

	if p.Index == "" {
		ret = str
	} else {
		ret = str + " use index(" + p.Index + ") "
	}

	return ret
}

func (p *UserManager) Insert(item *User) error {
	if p.Conn == nil {
		return errors.New("Connection Error")
	}

	if item.Date == "" {
		t := time.Now()
		item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}

	tableName := "user_tb"
	if config.Database == "mssql" || config.Database == "sqlserver" {
		tableName = config.Owner + ".user_tb"
	}

	query := "insert into " + tableName + " (u_loginid, u_passwd, u_name, u_date) values (?, ?, ?, ?)"
	res, err := p.Conn.Exec(query, item.Loginid, item.Passwd, item.Name, item.Date)
	if err == nil {
		p.Result = &res
	} else {
		log.Println(item)
		log.Println(err)
		p.Result = nil
	}

	return err
}
func (p *UserManager) Delete(id int64) error {
	if p.Conn == nil {
		return errors.New("Connection Error")
	}

	tableName := "user_tb"
	if config.Database == "mssql" || config.Database == "sqlserver" {
		tableName = config.Owner + ".user_tb"
	}
	query := "delete from " + tableName + " where u_id = ?"
	_, err := p.Conn.Exec(query, id)

	return err
}
func (p *UserManager) Update(item *User) error {
	if p.Conn == nil {
		return errors.New("Connection Error")
	}

	tableName := "user_tb"
	if config.Database == "mssql" || config.Database == "sqlserver" {
		tableName = config.Owner + ".user_tb"
	}

	query := "update " + tableName + " set u_loginid = ?,u_passwd = ?,u_name = ?,u_date = ? where u_id = ?"
	_, err := p.Conn.Exec(query, item.Loginid, item.Passwd, item.Name, item.Date, item.Id)

	return err
}

func (p *UserManager) GetIdentity() int64 {
	if p.Result == nil {
		return 0
	}

	id, err := (*p.Result).LastInsertId()

	if err != nil {
		return 0
	} else {
		return id
	}
}

func (p *UserManager) Get(id int64) *User {
	if p.Conn == nil {
		return nil
	}

	query := p.GetQuery() + " where u_id = ?"

	rows, err := p.Conn.Query(query, id)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return nil
	}

	defer rows.Close()

	if !rows.Next() {
		return nil
	}

	var item User
	err = rows.Scan(&item.Id, &item.Loginid, &item.Passwd, &item.Name, &item.Date)

	if err != nil {
		return nil
	} else {
		return &item
	}
}

func (p *UserManager) GetList(page int, pagesize int, order string) *[]User {
	if p.Conn == nil {
		return nil
	}

	startpage := (page - 1) * pagesize
	query := p.GetQuery()

	var rows *sql.Rows
	var err error

	if page > 0 && pagesize > 0 {
		if order == "" {
			order = "u_id desc"
		} else {
			order = "u_" + order
		}
		query += " order by " + order
		if config.Database == "mysql" {
			query += " limit ? offset ?"
			rows, err = p.Conn.Query(query, pagesize, startpage)
		} else if config.Database == "mssql" || config.Database == "sqlserver" {
			query += "OFFSET ? ROWS FETCH NEXT ? ROWS ONLY"
			rows, err = p.Conn.Query(query, startpage, pagesize)
		}
	} else {
		if order == "" {
			order = "u_id"
		} else {
			order = "u_" + order
		}
		query += " order by " + order
		rows, err = p.Conn.Query(query)
	}

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return nil
	}

	defer rows.Close()

	var items []User

	for rows.Next() {
		var item User
		err = rows.Scan(&item.Id, &item.Loginid, &item.Passwd, &item.Name, &item.Date)

		items = append(items, item)
	}

	if err != nil {
		return nil
	} else {
		return &items
	}
}

func (p *UserManager) GetCount() int {
	if p.Conn == nil {
		return 0
	}

	query := p.GetQuerySelect()

	rows, err := p.Conn.Query(query)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return 0
	}

	defer rows.Close()

	if !rows.Next() {
		return 0
	}

	cnt := 0
	err = rows.Scan(&cnt)

	if err != nil {
		return 0
	} else {
		return cnt
	}
}

func (p *UserManager) GetListInID(ids []int, page int, pagesize int, order string) *[]User {
	if p.Conn == nil {
		return nil
	}

	startpage := (page - 1) * pagesize
	query := p.GetQuery()

	var rows *sql.Rows
	var err error

	query = query + " where u_id in (" + strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ", ", -1), "[]") + ")"

	if page > 0 && pagesize > 0 {
		if order == "" {
			order = "u_id desc"
		} else {
			order = "u_" + order
		}
		query += " order by " + order
		if config.Database == "mysql" {
			query += " limit ? offset ?"
			rows, err = p.Conn.Query(query, pagesize, startpage)
		} else if config.Database == "mssql" || config.Database == "sqlserver" {
			query += "OFFSET ? ROWS FETCH NEXT ? ROWS ONLY"
			rows, err = p.Conn.Query(query, startpage, pagesize)
		}
	} else {
		if order == "" {
			order = "u_id"
		} else {
			order = "u_" + order
		}
		query += " order by " + order
		rows, err = p.Conn.Query(query)
	}

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return nil
	}

	defer rows.Close()

	var items []User

	for rows.Next() {
		var item User
		err = rows.Scan(&item.Id, &item.Loginid, &item.Passwd, &item.Name, &item.Date)

		items = append(items, item)
	}

	if err != nil {
		return nil
	} else {
		return &items
	}
}

func (p *UserManager) GetCountInID(ids []int) int {
	if p.Conn == nil {
		return 0
	}

	query := p.GetQuerySelect()

	query = query + " where u_id in (" + strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ", ", -1), "[]") + ")"

	rows, err := p.Conn.Query(query)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return 0
	}

	defer rows.Close()

	if !rows.Next() {
		return 0
	}

	cnt := 0
	err = rows.Scan(&cnt)

	if err != nil {
		return 0
	} else {
		return cnt
	}
}

func (p *UserManager) GetByLoginid(loginid string) *User {

	if p.Conn == nil {
		return nil
	}

	query := p.GetQuery() + " where 1=1 "
	var params []interface{}

	if loginid != "" {
		query += " and u_loginid = ?"
		params = append(params, loginid)
	}

	rows, err := QueryArray(p.Conn, query, params)

	if err != nil {
		log.Printf("query error : %v, %v\n", err, query)
		return nil
	}

	defer rows.Close()

	if rows.Next() {
		var item User
		err = rows.Scan(&item.Id, &item.Loginid, &item.Passwd, &item.Name, &item.Date)

		return &item
	} else {
		return nil
	}
}
