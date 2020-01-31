/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Wednesday 29 January 2020 - 16:54:02
** @Filename:				DEBUG.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Friday 31 January 2020 - 12:35:53
*******************************************************************************/

package postgre

import			"strconv"
import			"reflect"
import			"database/sql"
import			_ "github.com/lib/pq"
import			"github.com/microgolang/logs"

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
func	(q *S_Selector) All(receptacle interface{}) (interface{}, error) {
	tx, err := q.PGR.Begin()
	if (err != nil) {
		return nil, err
	}

	/**************************************************************************
	**	Assert the query string
	**************************************************************************/
	query := q.QuerySelect + ` ` + q.QueryFrom + ` ` + q.QueryWhere
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	/**************************************************************************
	**	Perfom the query
	**************************************************************************/
	rows, err := stmt.Query(q.Arguments...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var sliceItem reflect.Type
	items := reflect.TypeOf(receptacle).Elem()
	if items.Kind() == reflect.Ptr {
        items = items.Elem()
    }
    if items.Kind() == reflect.Slice {
		sliceItem = items
        items = items.Elem()
	}
	
	var myTypes []interface{}
	var myRows [][]interface{}
	for j := 0; j < items.NumField(); j++ {
		if (items.Field(j).Type.String() == `int`) {
			randomValue := 0
			myTypes = append(myTypes, &randomValue)
		} else if (items.Field(j).Type.String() == `string`) {
			randomValue := ``
			myTypes = append(myTypes, &randomValue)
		}
	}

	for rows.Next() {
		err = rows.Scan(myTypes...)
		myRows = append(myRows, myTypes)

		if err != nil {
			tx.Rollback()
			logs.Pretty(err)
			return nil, err
		}
	}
	rows.Close()

	receptArry := reflect.MakeSlice(sliceItem, len(myRows), len(myRows))
	for _, row := range myRows {
		for index := 0; index < receptArry.Len(); index++ {
			for j := 0; j < receptArry.Index(index).NumField(); j++ {
				typeOf := receptArry.Index(index).Field(j).Type().String()

				if (typeOf == `int`) {
					receptArry.Index(index).Field(j).Set(reflect.ValueOf(*(row[j]).(*int)))
				} else if (typeOf == `string`) {
					receptArry.Index(index).Field(j).Set(reflect.ValueOf(*(row[j]).(*string)))
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return receptArry.Interface(), nil
}
