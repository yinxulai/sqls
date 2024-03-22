# sqls

[![Go test](https://github.com/yinxulai/sqls/actions/workflows/test.yml/badge.svg)](https://github.com/yinxulai/sqls/actions/workflows/test.yml)

sqls is prepared for developers who like to write `sql` by themselves instead of using `ORM`. There is nothing special, it just simply helps you edit `sql` better.

> The meaning of sqls is `sql+s`, where `s` means `string` and `simple`, `small`

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
