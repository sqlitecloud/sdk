package sqlitecloud

import "fmt"
import "os"
import "bufio"
import "strings"
// import "errors"

func (this *SQCloud) Use( Database string ) error {
	_, err := this.Execute( fmt.Sprintf( "USE DATABASE %s", Database ) )
	return err
}

func (this *SQCloudResult ) Dump( MaxLine uint ) {
	this.bridge_Dump( MaxLine )
}



func (this *SQCloud) Compress( Enabled bool ) error {
	enabled := 0
	if Enabled {
		enabled = 1
	}
	_, err := this.Execute( fmt.Sprintf( "SET KEY CLIENT_COMPRESSION TO %d", enabled ) )
	return err
}

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


// func (this *SQCloud ) GetError() ( int, error ) {
// 	if this.bridge_IsError() {
// 		return this.bridge_GetErrorCode(), errors.New( this.bridge_GetErrorMessage() )
// 	}	
// 	return 0, nil
// }