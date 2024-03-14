# sqls

sqls 是针对喜欢自己编写 `sql` 而不是使用 `ORM` 的开发者准备的，没有任何特殊的东西，只是简单的帮助你更好的编辑 `sql`

> sqls 的含义是 `sql+s`，其中 `s` 的含义是 `string` 以及 `simple`、`small`

## 安装

```bash
go get github.com/yinxulai/sqls
```

## 使用

```go
s := Begin()
s.SELECT("P.ID, P.USERNAME, P.PASSWORD, P.FULL_NAME")
s.SELECT("P.LAST_NAME, P.CREATED_ON, P.UPDATED_ON")
s.FROM("PERSON P")
s.FROM("ACCOUNT A")
s.INNER_JOIN("DEPARTMENT D on D.ID = P.DEPARTMENT_ID")
s.INNER_JOIN("COMPANY C on D.COMPANY_ID = C.ID")
s.WHERE("P.ID = A.ID")
s.WHERE("P.FIRST_NAME like ?")
s.OR()
s.WHERE("P.LAST_NAME like ?")
s.GROUP_BY("P.ID")
s.HAVING("P.LAST_NAME like ?")
s.OR()
s.HAVING("P.FIRST_NAME like ?")
s.ORDER_BY("P.ID")
s.ORDER_BY("P.FULL_NAME")
sqlstring := s.String()
// sqlstring is:
// SELECT P.ID, P.USERNAME, P.PASSWORD, P.FULL_NAME, P.LAST_NAME, P.CREATED_ON, P.UPDATED_ON
// FROM PERSON P, ACCOUNT A
// INNER JOIN DEPARTMENT D on D.ID = P.DEPARTMENT_ID
// INNER JOIN COMPANY C on D.COMPANY_ID = C.ID
// WHERE (P.ID = A.ID AND P.FIRST_NAME like ?) OR (P.LAST_NAME like ?)
// GROUP BY P.ID
// HAVING (P.LAST_NAME like ?) OR (P.FIRST_NAME like ?)
// ORDER BY P.ID, P.FULL_NAME
```
