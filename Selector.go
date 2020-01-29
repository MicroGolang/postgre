/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 28 January 2020 - 18:49:41
** @Filename:				Selector.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Wednesday 29 January 2020 - 16:09:56
*******************************************************************************/

package			postgre

import			"strconv"
import			"database/sql"
import			_ "github.com/lib/pq"

type	S_Selector struct {
	PGR			*sql.DB
	QuerySelect	string
	QueryFrom	string
	QueryWhere	string
	Arguments	[]interface{}
}
type	S_SelectorWhere struct {
	Key	string
	Value string
}
func	NewSelector(PGR *sql.DB) (*S_Selector){
	return &S_Selector{PGR: PGR}
}
func	(q *S_Selector) Select(toSelect ...string) *S_Selector {
	q.QuerySelect = `SELECT `
	for index, selected := range toSelect {
		if (index > 0) {q.QuerySelect += `, `}
		q.QuerySelect += selected
	}
	return q
}
func	(q *S_Selector) From(table string) *S_Selector {
	q.QueryFrom = `FROM ` + table
	return q
}
func	(q *S_Selector) Where(asserts ...S_SelectorWhere) *S_Selector {
	q.QueryWhere = `WHERE `
	for index, each := range asserts {
		if (index > 0) {q.QueryWhere += ` AND `}
		q.QueryWhere += each.Key + `=`
		q.QueryWhere += `$` + strconv.Itoa(index + 1)
		q.Arguments = append(q.Arguments, each.Value)
	}
	return q
}
func	(q *S_Selector) One(receptacle ...interface{}) (error) {
	tx, err := q.PGR.Begin()
	if (err != nil) {
		return err
	}

	/**************************************************************************
	**	Assert the query string
	**************************************************************************/
	query := q.QuerySelect + ` ` + q.QueryFrom + ` ` + q.QueryWhere + ` LIMIT 1`
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


	for rows.Next() {
		err = rows.Scan(receptacle...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	rows.Close()

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}