/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 28 January 2020 - 18:48:45
** @Filename:				Deletor.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Tuesday 28 January 2020 - 18:49:26
*******************************************************************************/

package			postgre

import			"database/sql"
import			_ "github.com/lib/pq"

/******************************************************************************
**	DELETOR
*******************************************************************************/
type	S_Deletor struct {
	PGR			*sql.DB
	QueryTable	string
	QueryWhere	string
}
func	NewDeletor(PGR *sql.DB) (*S_Deletor){
	return &S_Deletor{PGR: PGR}
}
func	(q *S_Deletor) Into(table string) *S_Deletor {
	q.QueryTable = `DELETE FROM ` + table
	return q
}
func	(q *S_Deletor) Where(assert string) *S_Deletor {
	q.QueryWhere = `WHERE ` + assert
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
	rows, err := stmt.Query()
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