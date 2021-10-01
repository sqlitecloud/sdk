//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.1.0
//     //             ///   ///  ///    Date        : 2021/10/01
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : Auxiliary GO methods
//   ////                ///  ///                     related to the SQCloud
//     ////     //////////   ///                      class.
//        ////            ////                        ( Convenience API's )
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import "fmt"
import "errors"

type SQCloudKeyValues struct {
  Key         string
  Value       string
}

// BeginTransaction starts a transaction. 
// To finish a transaction, see: SQCloud.EndTransaction(), to undo a transaction, see: SQCloud.RollBackTransaction()
func (this *SQCloud) BeginTransaction() error {
  return this.Execute( "BEGIN TRANSACTION" )
}

// EndTransaction finishes a transaction and writes the changes to the disc. 
// To start a transaction, see: SQCloud.BeginTransaction(), to undo a transaction, see: SQCloud.RollBackTransaction().
func (this *SQCloud) EndTransaction() error {
  return this.Execute( "END TRANSACTION" )
}

// RollBackTransaction aborts a transaction without the writing the changes to the disc.
// To start a transaction, see: SQCloud.BeginTransaction(), to finish a transaction, see: SQCloud.EndTransaction().
func (this *SQCloud) RollBackTransaction() error {
  return this.Execute( "ROLLBACK TRANSACTION" )
}

func (this *SQCloud) GetAutocompleteTokens() ( tokens []string ) {
  if tables, err := this.ListTables(); err == nil {
    for _, table := range tables {
      if table != "sqlite_sequence" {
        tokens = append( tokens, table )
        for _, column := range this.ListColumns( table ) {
          tokens = append( tokens, fmt.Sprintf( "%s", column ) )
          tokens = append( tokens, fmt.Sprintf( "%s.%s", table, column ) )
        }
      }
    }
  }
  return
}

func (this *SQCloud) ListColumns( TableName string ) ( columns []string ) {
  res, err := this.Select( fmt.Sprintf( "pragma table_info( %s )", SQCloudEnquoteString( TableName ) ) ) 
  if res != nil {
    defer res.Free()

    if err ==nil {
      for row, err := res.GetFirstRow(); row != nil; row = row.Next() {
        switch {
        case err != nil:    break
        case row == nil:    break
        default:  
          switch val, err := row.GetString( 1 ); {
          case err != nil:  break
          default:          columns = append( columns, val )
          }
        }
      }
    }
  }
  return
}


// SelectSingleString executes the query in SQL and returns the result as string.
// The given query must return one single value. If no value or more than one values are selected by the query, nil is 
// returned and error describes the problem that occurred.
// See: SQCloud.SelectSingleInt64(), SQCloud.SelectStringList(), SQCloud.SelectKeyValues()
func (this *SQCloud) SelectSingleString( SQL string ) ( string, error ) {
  if result, err := this.Select( SQL ); result == nil {
    return "", errors.New( "ERROR: Query returned no result (-1)" )

  } else {
    defer result.Free()

    switch {
    case err != nil:                          return "", err
    case result.IsError():                    return "", errors.New( result.GetErrorAsString() )
    case result.IsString():                   return result.GetString()
    case !result.IsRowSet():                  return "", errors.New( "ERROR: Query returned an invalid result" )
    case result.GetNumberOfRows() != 1:       return "", errors.New( "ERROR: Query returned not exactly one row" )
    case result.GetNumberOfColumns() != 1:    return "", errors.New( "ERROR: Query returned not exactly one column" )
    case result.GetValueType_( 0, 0 ) != '+': return "", errors.New( "ERROR: Query returned not a string" )
    default:                                  return result.GetStringValue( 0, 0 )
} } }


// SelectSingleInt64 executes the query in SQL and returns the result as int64.
// The given query must return one single numerical value. If no value, more than one values or a non-numerical value is selected by the query, 0 is 
// returned and error describes the problem that occurred.
// See: SQCloud.SelectSingleString(), SQCloud.SelectStringList(), SQCloud.SelectKeyValues()
func (this *SQCloud) SelectSingleInt64( SQL string ) ( int64, error ) {
  if result, err := this.Select( SQL ); result == nil {
    return 0, errors.New( "ERROR: Query returned no result (-1)" )

  } else {
    defer result.Free()

    switch {
    case err != nil:                          return 0, err
    case result.IsError():                    return 0, errors.New( result.GetErrorAsString() )
    case result.IsInteger():                  return result.GetInt64()
    case !result.IsRowSet():                  return 0, errors.New( "ERROR: Query returned an invalid result" )
    case result.GetNumberOfRows() != 1:       return 0, errors.New( "ERROR: Query returned not exactly one row" )
    case result.GetNumberOfColumns() != 1:    return 0, errors.New( "ERROR: Query returned not exactly one column" )
    case result.GetValueType_( 0, 0 ) != ':': return 0, errors.New( "ERROR: Query returned not an integer" )
    default:                                  return result.GetInt64Value( 0, 0 )
} } }

// SelectStringList executes the query in SQL and returns the result as an array of strings.
// The given query must result in 0 or more rows and one column. 
// If no row was selected, an empty string array is returned.
// If more than one column was selected, an empty string array and an error describing the problem is returned.
// In all other cases, an array of strings is returned.
// See: SQCloud.SelectSingleString(), SQCloud.SelectSingleInt64(), SQCloud.SelectKeyValues()
func (this *SQCloud) SelectStringList( SQL string ) ( []string, error ) {
  if result, err := this.Select( SQL ); result == nil {
    return []string{}, errors.New( "ERROR: Query returned no result (-1)" )

  } else {
    defer result.Free()

    switch {
    case err != nil:                          return []string{}, err
    case result.IsError():                    return []string{}, errors.New( result.GetErrorAsString() )
    case result.IsString():                   return []string{ result.GetString_() }, nil
    case !result.IsRowSet():                  return []string{}, errors.New( "ERROR: Query returned an invalid result" )
    case result.GetNumberOfColumns() != 1:    return []string{}, errors.New( "ERROR: Query returned not exactly one column" )
    default:                                  
      stringList := []string{}
      for _, row := range result.Rows() {
        switch val, err := row.GetValue( 0 ); {
        case err != nil:                      return []string{}, err
        case val == nil:                      return []string{}, errors.New( "ERROR: Value in query result is invalid" )
        case !val.IsString():                 return []string{}, errors.New( "ERROR: Value in query result is not a string" )
        default:                              stringList = append( stringList, val.GetString() )
      } }
      return stringList, nil
} } }

// SelectKeyValues executes the query in SQL and returns the result as an array of SQCloudKeyValues.
// The given query must result in 0 or more rows and two columns. 
// If no row was selected, an empty array is returned.
// If less or more than two columns where selected, an empty array and an error describing the problem is returned.
// In all other cases, an array of SQCloudKeyValues is returned.
// See: SQCloud.SelectSingleString(), SQCloud.SelectSingleInt64(), SQCloud.SelectStringList()
// BUG(andreas): Rewrite and use map[string]string
func (this *SQCloud) SelectKeyValues( SQL string ) ( []SQCloudKeyValues, error ) {
  if result, err := this.Select( SQL ); result == nil {
    return []SQCloudKeyValues{}, errors.New( "ERROR: Query returned no result (-1)" )

  } else {
    defer result.Free()

    switch {
    case err != nil:                          return []SQCloudKeyValues{}, err
    case result.IsError():                    return []SQCloudKeyValues{}, errors.New( result.GetErrorAsString() )
    case !result.IsRowSet():                  return []SQCloudKeyValues{}, errors.New( "ERROR: Query returned an invalid result" )
    case result.GetNumberOfColumns() != 2:    return []SQCloudKeyValues{}, errors.New( "ERROR: Query returned not exactly one column" )
    default:                                  
    keyValueList := []SQCloudKeyValues{}
      for _, row := range result.Rows() {
        key, kerr := row.GetValue( 0 )
        val, verr := row.GetValue( 1 )
        switch {
        case kerr != nil:                     return []SQCloudKeyValues{}, kerr
        case verr != nil:                     return []SQCloudKeyValues{}, verr
        case key == nil:                      return []SQCloudKeyValues{}, errors.New( "ERROR: Key in query result is invalid" )
        case val == nil:                      return []SQCloudKeyValues{}, errors.New( "ERROR: Value in query result is invalid" )
        case !key.IsString():                 return []SQCloudKeyValues{}, errors.New( "ERROR: Key in query result is not a string" )
        case !val.IsString():                 return []SQCloudKeyValues{}, errors.New( "ERROR: Value in query result is not a string" )
        default:                              keyValueList = append( keyValueList, SQCloudKeyValues{ Key: key.GetString(), Value: val.GetString() } )
      } }
      return keyValueList, nil
} } }