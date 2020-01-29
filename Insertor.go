/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 28 January 2020 - 17:48:49
** @Filename:				Insertor.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Wednesday 29 January 2020 - 15:32:44
*******************************************************************************/

package			postgre

import			"strconv"
import			"database/sql"
import			_ "github.com/lib/pq"

type	S_Insertor struct {
	PGR			*sql.DB
	QueryAs		string
	QueryValues	string
	QueryTable	string
	QueryWhere	string
	Arguments	[]interface{}
}
type	S_InsertorWhere struct {
	Key	string
	Value string
}
func	NewInsertor(PGR *sql.DB) (*S_Insertor){
	return &S_Insertor{PGR: PGR}
}
func	(q *S_Insertor) Values(values ...S_InsertorWhere) *S_Insertor {
	q.QueryAs += `(`
	q.QueryValues = `VALUES (`
	for index, value := range values {
		if (index > 0) {
			q.QueryValues += `, `
			q.QueryAs += `, `
		}
		q.QueryValues += `$` + strconv.Itoa(index + 1)
		q.QueryAs += value.Key
		q.Arguments = append(q.Arguments, value.Value)
	}
	q.QueryValues += `)`
	q.QueryAs += `)`
	return q
}
func	(q *S_Insertor) Into(table string) *S_Insertor {
	q.QueryTable = `INSERT INTO ` + table
	return q
}
func	(q *S_Insertor) Do() (string, error) {
	tx, err := q.PGR.Begin()
	if (err != nil) {
		return ``, err
	}

	/**************************************************************************
	**	Assert the query string
	**************************************************************************/
	query := q.QueryTable + ` ` + q.QueryAs + ` ` + q.QueryValues + ` RETURNING ID`
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return ``, err
	}
	defer stmt.Close()

	/**************************************************************************
	**	Perfom the query
	**************************************************************************/
	rows, err := stmt.Query(q.Arguments...)
	if err != nil {
		tx.Rollback()
		return ``, err
	}

	ID := ``
	rows.Next()
	rows.Scan(&ID)
	rows.Close()

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return ``, err
	}
	return ID, nil
}