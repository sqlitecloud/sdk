//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/08/13
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : Auxiliary GO methods
//   ////                ///  ///                     related to the SQCloud
//     ////     //////////   ///                      class.
//        ////            ////          
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

///// Convenience API's


// BeginTransaction starts a transaction. 
// To finish a transaction, see: SQCloud.EndTransaction(), to undo a transaction, see: SQCloud.RollBackTransaction()
// BUG(andreas): ToDo
func (this *SQCloud) BeginTransaction() {
}

// EndTransaction finishes a transaction and writes the changes to the disc. 
// To start a transaction, see: SQCloud.BeginTransaction(), to undo a transaction, see: SQCloud.RollBackTransaction().
// BUG(andreas): ToDo
func (this *SQCloud) EndTransaction() {
}

// RollBackTransaction aborts a transaction without the writing the changes to the disc.
// To start a transaction, see: SQCloud.BeginTransaction(), to finish a transaction, see: SQCloud.EndTransaction().
// BUG(andreas): ToDo
func (this *SQCloud) RollBackTransaction() {
  
}

// AutoCommit enables/disables the immediate writing of changes to the disc.
// BUG(andreas): ToDo
func (this *SQCloud) AutoCommit( Enabled bool ) {
  
}

func (this *SQCloud) GetAutocompleteTokens() ( tokens []string ) {
  if tables, err := this.ListTables(); err == nil {
    for _, table := range tables {
      tokens = append( tokens, table )
      for _, column := range this.ListColumns( table ) {
        tokens = append( tokens, fmt.Sprintf( "%s", column ) )
        tokens = append( tokens, fmt.Sprintf( "%s.%s", table, column ) )
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
  result, err := this.Select( SQL )
  if result != nil {
    defer result.Free()

    if err == nil {
      if result.GetNumberOfColumns() == 1 && result.GetNumberOfRows() == 1 {
        return result.GetStringValue( 0, 0 )
      }
      return "", errors.New( "ERROR: Query returned not exactly one value (-1)" )
    } 
    return "", err
  }
  return "", errors.New( "ERROR: Query returned no result (-1)" )
}

// SelectSingleInt64 executes the query in SQL and returns the result as int64.
// The given query must return one single numerical value. If no value, more than one values or a non-numerical value is selected by the query, 0 is 
// returned and error describes the problem that occurred.
// See: SQCloud.SelectSingleString(), SQCloud.SelectStringList(), SQCloud.SelectKeyValues()
func (this *SQCloud) SelectSingleInt64( SQL string ) ( int64, error ) {
  result, err := this.Select( SQL )
  if result != nil {
    defer result.Free()

    if err == nil {
      if result.GetNumberOfColumns() == 1 && result.GetNumberOfRows() == 1 {
        return result.GetInt64Value( 0, 0 )
      }
      return 0, errors.New( "ERROR: Query returned not exactly one value (-1)" )
    } 
    return 0, err
  }
  return 0, errors.New( "ERROR: Query returned no result (-1)" )
}

// SelectStringList executes the query in SQL and returns the result as an array of strings.
// The given query must result in 0 or more rows and one column. 
// If no row was selected, an empty string array is returned.
// If more than one column was selected, an empty string array and an error describing the problem is returned.
// In all other cases, an array of strings is returned.
// See: SQCloud.SelectSingleString(), SQCloud.SelectSingleInt64(), SQCloud.SelectKeyValues()
func (this *SQCloud) SelectStringList( SQL string ) ( []string, error ) {
  result, err := this.Select( SQL )

  if result != nil {
    defer result.Free()
    
    if stringList := []string{}; err == nil {
      for _, row := range result.Rows() {
        switch val, err := row.GetValue( 0 ); {
        case err != nil:  return []string{}, errors.New( "Invalid server response" )
        default:          stringList = append( stringList, val.GetString() )
        }
      }
      return stringList, nil
    }
    return []string{}, err  
  }
  return []string{}, nil
}

// SelectKeyValues executes the query in SQL and returns the result as an array of SQCloudKeyValues.
// The given query must result in 0 or more rows and two columns. 
// If no row was selected, an empty array is returned.
// If less or more than two columns where selected, an empty array and an error describing the problem is returned.
// In all other cases, an array of SQCloudKeyValues is returned.
// See: SQCloud.SelectSingleString(), SQCloud.SelectSingleInt64(), SQCloud.SelectStringList()
// BUG(andreas): Rewrite and use map[string]string
func (this *SQCloud) SelectKeyValues( SQL string ) ( []SQCloudKeyValues, error ) {
  result, err := this.Select( SQL )

  if result != nil {
    defer result.Free()

    if keyValueList := []SQCloudKeyValues{}; err == nil {
      for _, row := range result.Rows() {
        key, kerr := row.GetValue( 0 )
        val, verr := row.GetValue( 1 )
        switch {
        case kerr != nil, verr != nil: return []SQCloudKeyValues{}, errors.New( "Invalid server response" )
        default: keyValueList = append( keyValueList, SQCloudKeyValues{ Key: key.GetString(), Value: val.GetString() } )
        }
      }
      return keyValueList, nil
    }
    return []SQCloudKeyValues{}, err
  }
  return []SQCloudKeyValues{}, nil
}