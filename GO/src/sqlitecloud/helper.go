package sqlitecloud

import "fmt"
//import "os"
//import "bufio"
import "strings"
import "errors"
//import "time"
import "strconv"

// Helper functions

// SQCloudEnquoteString enquotes the given string if necessary and returns the result as a news created string.
// If the given string contains a '"', the '"' is properly escaped.
// If the given string contains one or more spaces, the whole string is enquoted with '"'
func SQCloudEnquoteString( Token string ) string {
  Token = strings.Replace( Token, "\"", "\"\"", -1 )
  if strings.Contains( Token, " " ) {
    return fmt.Sprintf( "\"%s\"", Token )
  }
	if Token == "" {
		return "\"\""
	}
  return Token
}

// parseBool parses the given string value and tries to interpret the value as a bool.
// true is returned if the value is "TRUE", "ENABLED" or 1.
// false is returned if the value is "FALSE", "DISABLED" or 0.
// The specified defaultValue is returned, if the given string value was an emptry string. 
// An error is returned in any other case.
func parseBool( value string, defaultValue bool ) ( bool, error ) {
  switch strings.ToUpper( strings.TrimSpace( value ) ) {
    case "FALSE", "DISABLED", "0": return false, nil
    case "TRUE", "ENABLED", "1"  : return true, nil
    case ""                      : return defaultValue, nil
    default                      : return false, errors.New( "ERROR: Not a Boolean value" )
  } 
}

// parseInt parses the given string value and tries to extract its value as an int value.
// If value was an empty string, the defaultValue is evaluated instead.
// If the given string value does not resemble a numeric value or its numeric value is smaler than minValue or exceeds maxValue, an error describing the problem is returned.
func parseInt( value string, defaultValue int, minValue int, maxValue int ) ( int, error ) {
  // println( "ParseInt = " + value )
  value = strings.TrimSpace( value )
  if value == "" {
    value = fmt.Sprintf( "%d", defaultValue )
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

// parseString returns a non empty string.
// The given string value is trimmed.
// If the given string is an empty string, the defaultValue is evaluated instead.
// If the given string and the defaultValue are emptry strings, an error is returned.
func parseString( value string, defaultValue string ) ( string, error ) {
  value = strings.TrimSpace( value )
  if value == "" {
    value = strings.TrimSpace( defaultValue )
  }
  if value == "" {
    return "", errors.New( "ERROR: Empty value" )
  }
  return value, nil
}