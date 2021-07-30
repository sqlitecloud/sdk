package sqlitecloud

import "fmt"
//import "os"
//import "bufio"
import "strings"
import "errors"
//import "time"
import "strconv"

// Helper functions

func SQCloudEnquoteString( Token string ) string {
	Token = strings.Replace( Token, "\"", "\"\"", -1 )
	if strings.Contains( Token, " " ) {
		return fmt.Sprintf( "\"%s\"", Token )
	}
	return Token
}

func parseBool( value string, defaultValue bool ) ( bool, error ) {
	switch strings.ToUpper( strings.TrimSpace( value ) ) {
	  case "FALSE", "DISABLED", "0": return false, nil
	  case "TRUE", "ENABLED", "1"  : return true, nil
	  case ""                      : return defaultValue, nil
	  default                      : return false, errors.New( "ERROR: Not a Boolean value" )
	}	
}

func parseInt( value string, defaultValue int, minValue int, maxValue int ) ( int, error ) {
	println( "ParseInt = " + value )
	value = strings.TrimSpace( value )
	if value == "" {
		return defaultValue, nil
	}
	if v, err := strconv.Atoi( value ); err == nil {
		if v < minValue {
			return 0, errors.New( "ERROR: The given Number is too small" )
		}
		if v > maxValue {
			return 0, errors.New( "ERROR: The given Number is too large" )
		}
		return v, nil
	} else {
		return 0, err
	}
}

func parseString( value string, defaultValue string ) ( string, error ) {
	value = strings.TrimSpace( value )
	if value == "" {
		if defaultValue != "" {
			return defaultValue, nil
		}
		return "", errors.New( "ERROR: Empty value" )
	}
	return value, nil
}


