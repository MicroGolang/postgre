/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 28 January 2020 - 18:48:45
** @Filename:				Deletor.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Tuesday 28 January 2020 - 19:04:34
*******************************************************************************/

package			postgre

import			"strconv"
import			"database/sql"
import			_ "github.com/lib/pq"

/******************************************************************************
**	DELETOR
*******************************************************************************/
type	S_Deletor struct {
	PGR			*sql.DB
	QueryTable	string
	QueryWhere	string
	Arguments	[]interface{}
}
type	S_DeletorWhere struct {
	Key	string
	Value string
}
func	NewDeletor(PGR *sql.DB) (*S_Deletor){
	return &S_Deletor{PGR: PGR}
}
func	(q *S_Deletor) Into(table string) *S_Deletor {
	q.QueryTable = `DELETE FROM ` + table
	return q
}
func	(q *S_Deletor) Where(asserts ...S_DeletorWhere) *S_Deletor {
	q.QueryWhere = `WHERE `
	for index, each := range asserts {
		if (index > 0) {q.QueryWhere += `, `}
		q.QueryWhere += each.Key + `=`
		q.QueryWhere += `$` + strconv.Itoa(index + 1)
		q.Arguments = append(q.Arguments, each.Value)
	}
	return q
}
func	(q *S_Deletor) Do() (error) {
	tx, err := q.PGR.Begin()
	if (err != nil) {
		return err
	}

	/**************************************************************************
	**	Assert the query string
	**************************************************************************/
	query := q.QueryTable + ` ` + q.QueryWhere
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