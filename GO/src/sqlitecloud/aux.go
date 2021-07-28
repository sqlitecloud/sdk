package sqlitecloud

import "fmt"
import "os"
import "bufio"
import "strings"
//import "time"
import "errors"

///// Convenience API's

func (this *SQCloud) ExecuteFiles( FilePathes []string ) error {
	for _, file := range FilePathes {
		err := this.ExecuteFile( file )
		if( err != nil ) {
			return err
		}
	}
	return nil
}

func (this *SQCloud) ExecuteFile( FilePath string ) error {
	file, err := os.Open( FilePath )
	if err == nil {
		defer file.Close()

		line := bufio.NewScanner( file )
		for line.Scan() {
			if strings.ToUpper( line.Text() ) != ".PROMPT" {
				fmt.Println( ">> %s\r\n", line.Text() )
				this.Execute( line.Text() )
				continue
			}
			return nil
		}
		return line.Err()
	}
	return err
}

func (this *SQCloud) BeginTransaction() {

}
func (this *SQCloud) EndTransaction() {
	
}
func (this *SQCloud) RollBackTransaction() {
	
}

func (this *SQCloud) AutoCommit( Enabled bool ) {
	
}

func (this *SQCloud) Compress( Enabled bool ) error {
	switch Enabled {
	  case false: return this.Execute( "SET KEY CLIENT_COMPRESSION TO 0" )
	  default: 		return this.Execute( "SET KEY CLIENT_COMPRESSION TO 1" )
	}
}




func (this *SQCloud) SelectSingleString( SQL string ) ( string, error ) {
	result, err := this.SelectStringList( SQL )
	if err != nil {
		switch len( result ) {
			case 0:  return ""         , errors.New( "ERROR: Query returned no value (-1)" )
			case 1:  return result[ 1 ], nil
			default: return ""         , errors.New( "ERROR: Query returned too many values (-1)" ) 
		}
	}
	return "", err
}
func (this *SQCloud) SelectSingleInt64( SQL string ) ( int64, error ) {
	result, err := this.Select( SQL )
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 1 {
				if result.GetNumberOfRows() == 1 {
					val := result.CGetInt64Value( 0, 0 )
					result.Free()
					return val, nil
				}
			}
			result.Free()
			return 0, errors.New( "ERROR: Query returned not exactly one value (-1)" )
		}
		return 0, errors.New( "ERROR: Query returned no result (-1)" )
	}
	return 0, err
}

func (this *SQCloud) SelectStringList( SQL string ) ( []string, error ) {
	stringList := []string{}
	result, err := this.Select( SQL )
	//result.Dump( 150 )
	//println( result.GetNumberOfColumns() )
	//println( result.GetNumberOfRows() )
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 1 {
				rows :=result.GetNumberOfRows()
				for row := uint( 0 ); row < rows; row++ {
					// println( result.CGetStringValue( row, 0 ) )
					stringList = append( stringList, result.CGetStringValue( row, 0 ) )
				}
				result.Free()
				return stringList, nil
			}
			result.Free()
			return []string{}, errors.New( "ERROR: Query returned not 1 Column (-1)" )
		}
		return []string{}, errors.New( "ERROR: Query returned no result (-1)" )
	}
	return []string{}, err	
}

func (this *SQCloud) SelectKeyValues( SQL string ) ( []SQCloudKeyValues, error ) {
	keyValueList := []SQCloudKeyValues{}
	result, err := this.Select( SQL )
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 2 {
				rows :=result.GetNumberOfRows()
				for row := uint( 0 ); row < rows; row++ {
					keyValueList = append( keyValueList, SQCloudKeyValues{ Key: result.CGetStringValue( row, 0 ), Value: result.CGetStringValue( row, 1 ) } )
				}
				result.Free()
				return keyValueList, nil
			}
			result.Free()
			return []SQCloudKeyValues{}, errors.New( "ERROR: Query returned not 2 Columns (-1)" )
		}
		return []SQCloudKeyValues{}, errors.New( "ERROR: Query returned no result (-1)" )
	}
	return []SQCloudKeyValues{}, err
}