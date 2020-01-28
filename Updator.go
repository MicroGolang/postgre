/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 28 January 2020 - 17:57:16
** @Filename:				Updator.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Tuesday 28 January 2020 - 17:59:04
*******************************************************************************/

package			postgre

import			"strconv"
import			"database/sql"
import			_ "github.com/lib/pq"

type	S_Updator struct {
	PGR			*sql.DB
	QueryValues	string
	QueryTable	string
	QueryWhere	string
	Arguments	[]interface{}
}
type	S_UpdatorSetter struct {
	key	string
	value string
}
type	S_UpdatorWhere struct {
	key	string
	value string
}
func	NewUpdator(PGR *sql.DB) (*S_Updator){
	return &S_Updator{PGR: PGR}
}
func	(q *S_Updator) Set(values []S_UpdatorSetter) *S_Updator {
	q.QueryValues = `SET `
	for index, each := range values {
		if (index > 0) {q.QueryValues += `, `}
		q.QueryValues += each.key + `=`
		q.QueryValues += `$` + strconv.Itoa(index + 1)
		q.Arguments = append(q.Arguments, each.value)
	}
	return q
}
func	(q *S_Updator) Into(table string) *S_Updator {
	q.QueryTable = `UPDATE ` + table
	return q
}
func	(q *S_Updator) Where(asserts ...S_UpdatorWhere) *S_Updator {
	initialIndex := len(q.Arguments)
	q.QueryWhere = `WHERE `
	for index, each := range asserts {
		if (index > 0) {q.QueryWhere += `, `}
		q.QueryWhere += each.key + `=`
		q.QueryWhere += `$` + strconv.Itoa(index + initialIndex + 1)
		q.Arguments = append(q.Arguments, each.value)
	}
	return q
}
func	(q *S_Updator) Do() (error) {
	tx, err := q.PGR.Begin()
	if (err != nil) {
		return err
	}

	/**************************************************************************
	**	Assert the query string
	**************************************************************************/
	query := q.QueryTable + ` ` + q.QueryValues + ` ` + q.QueryWhere
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	/**************************************************************************
	**	Perfom the query
	**************************************************************************/
	rows, err := stmt.Query(q.Arguments...)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
