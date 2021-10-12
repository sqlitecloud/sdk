//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.4
//     //             ///   ///  ///    Date        : 2021/10/05
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : Go Methods related to the
//   ////                ///  ///                     SQCloud class for managing
//     ////     //////////   ///                      the connection and executing
//        ////            ////                        queries.
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import "io"
import "fmt"
import "net"
import "time"
import "errors"
import "strconv"
import "github.com/pierrec/lz4"

type Chunk struct {
  DataBufferOffset  uint64
  LEN               uint64
  RAW               []byte
}

func( this *Chunk ) GetType()      byte   { return this.RAW[ 0 ]             }
func( this *Chunk ) IsCompressed() bool   { return this.GetType() == '%'     }
func( this* Chunk ) GetChunkSize() uint64 { return uint64( len( this.RAW ) ) }
func( this* Chunk ) GetData()      []byte {
  switch this.RAW {
  case nil: return []byte{ '_', ' ' }
  default:  return this.RAW[ this.DataBufferOffset : ]
} }

func( this* Chunk ) Uncompress() error {
  // %TLEN CLEN ULEN *0 NROWS NCOLS <Compressed DATA>
  // %TLEN CLEN ULEN /0 NROWS NCOLS <Compressed DATA>

  if this.RAW == nil      { return errors.New( "Nil pointer exception" ) }
  if !this.IsCompressed() { return nil                                   }

  var err           error

  var hStartIndex   uint64 = 1 // Index of the start of the uncompressed header in chunk (*0 NROWS NCOLS ...)
  var zStartIndex   uint64 = 0 // Index of the start of the compressed buffer in this chunk (<Compressed DATA...>)

  var LEN           uint64 = 0
  var lLEN          uint64 = 0

  var COMPRESSED    uint64 = 0
  var cLEN          uint64 = 0

  var UNCOMPRESSED  uint64 = 0
  var iUNCOMPRESSED int    = 0
  var uLEN          uint64 = 0

  LEN, lLEN, err = this.readUInt64At( hStartIndex )             // "%TLEN "
  hStartIndex += lLEN                                           // hStartIndex -> "CLEN ULEN *0 NROWS NCOLS <Compressed DATA...>"

  COMPRESSED, cLEN, err = this.readUInt64At( hStartIndex )      // "CLEN "
  hStartIndex += cLEN                                           // hStartIndex -> "ULEN *0 NROWS NCOLS <Compressed DATA...>"

  UNCOMPRESSED, uLEN, err = this.readUInt64At( hStartIndex )    // "ULEN "
  hStartIndex += uLEN                                           // hStartIndex -> "*0 NROWS NCOLS <Compressed DATA...>"

  zStartIndex = LEN - COMPRESSED + lLEN + 1                     // zStartIndex -> "<Compressed DATA...>"
  hLEN := zStartIndex - hStartIndex                             // = len( "*0 NROWS NCOLS " )

  newHeader := fmt.Sprintf( "%c%d %s", this.RAW[ hStartIndex ], UNCOMPRESSED + hLEN - 3, string( this.RAW[ hStartIndex + 3 : hStartIndex + hLEN ] ) )

  this.DataBufferOffset = uint64( len( newHeader ) )            // = len( "*200020 1000 2 " )

  buf := make( []byte, this.DataBufferOffset + UNCOMPRESSED )   // allocate memory
  copy( buf[ 0 : this.DataBufferOffset ], []byte( newHeader ) ) // copy the new header into the memory

  if iUNCOMPRESSED, err = lz4.UncompressBlock( this.RAW[ zStartIndex : ], buf[ this.DataBufferOffset : ] ); err != nil { return err }

  // Overwrite old Buffer with uncompressed one
  this.LEN              = uint64( iUNCOMPRESSED ) + hLEN - 3    // see: newHeader :=...
  this.RAW              = buf

  return nil
}

func (this *Chunk ) readUInt64At( offset uint64 ) ( uint64, uint64, error ) {
  if this.RAW == nil { return 0, 0, errors.New( "Nil chunk" ) }

  var zero        uint64 = uint64( '0' ) 
  var val         uint64 = 0
  var bytesRead   uint64 = 0
  var maxLEN      int = len( this.RAW ) - int( offset ) // 0...end of chunk                                       
  
  if maxLEN < 0   { maxLEN = 0  }
  if maxLEN > 20  { maxLEN = 20 }                       // MaxInt64 = 18446744073709551615 (len=20)

  for {
    if bytesRead == uint64( maxLEN ) { return val, bytesRead, nil }
    switch c := this.RAW[ bytesRead + offset ]; c {
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':  val = val * 10 + ( uint64( c ) - zero )
    case ' ':                                               return val, bytesRead + 1, nil
    default:                                                return 0, 0, errors.New( "Invalid rune" )
    }
    bytesRead++
  }
  return 0, 0, errors.New( "Overflow" )
}

func (this *Chunk ) readValueAt( offset uint64 ) ( Value, uint64, error ) {
  value :=  Value{ Type: this.RAW[ offset ], Buffer: nil }
  switch bytesRead, err := value.readBufferAt( this, offset + 1 ); {
  case err != nil:      return Value{ Type: 0, Buffer: nil }, 0, err
  case bytesRead == 0:  return Value{ Type: 0, Buffer: nil }, 0, errors.New( "End Of Chunk" )
  default:              return value, 1 + bytesRead, nil
  }
}

func (this *Value ) readBufferAt( chunk *Chunk, offset uint64 ) ( uint64, error ) {
  var bytesLeft uint64 = 0
  if chunk.GetChunkSize() > offset { bytesLeft = chunk.GetChunkSize() - offset }
  this.Buffer = nil

  switch this.Type {
  case '_': return 1, nil

  case '+', '!', '-', ':', ',', '$', '#', '^', '@', '|':
    var TRIM  uint64 = 0  // Trims if it is a the C-String

    switch this.Type {
    case ':':             // Space terminated INT
      // MaxInt64         = 18446744073709551615                                = 20 bytes
      // MinInt64         = -9223372036854775807                                = 20 bytes <- MAX LEN
      if bytesLeft > 20 { bytesLeft = 20 }
      fallthrough

    case ',':             // Space terminated FLOAT
      // MaxFloat32       = 3.40282346638528859811704183484516925440e+38        = 44 bytes
      // SmalestFloat32   = 1.401298464324817070923729583289916131280e-45       = 45 bytes
      // FaxFloat64       = 1.79769313486231570814527423731704356798070e+308    = 48 bytes
      // SmallestFloat64  = 4.9406564584124654417656879286822137236505980e-324  = 50 bytes <- MAX LEN
      if bytesLeft > 50 { bytesLeft = 50 }

      bytesRead := uint64( 0 );
      for ; bytesRead < bytesLeft; bytesRead++ {
        if chunk.RAW[ offset + bytesRead ] == ' ' { 
          bytesRead++
          break 
        }
        this.Buffer = chunk.RAW[ offset : offset + 1 + bytesRead ]
      }
      if len( this.Buffer ) == 0 { return 0, errors.New( "End Of Chunk" ) }
      return bytesRead, nil

    case '!':                 // Zero terminated C-String
      TRIM = 1                // Cut one byte off the buffer / dont copy the zero byte of the C string
      fallthrough

    default:                  // Everything else is a LEN Value (+!-$#^@)
      switch LEN, len, err := chunk.readUInt64At( offset ); {
      case err != nil:        return 0, err
      case len == 0:          return 0, errors.New( "LEN not found" )
      default:                this.Buffer = chunk.RAW[ offset + len : offset + len + LEN - TRIM ]
                              return len + LEN, nil
      }
    }
  }
  return 0, errors.New( "Unsuported type" )
}

func ( this *SQCloud ) sendString( data string ) ( int, error ) {
  var err         error
  var bytesSent   int
  var bytesToSend int

  if err := this.reconnect()                                                           ; err != nil { return 0, err }
  switch this.Timeout {
  case 0:   if err := ( *this.sock ).SetWriteDeadline( time.Time{} );                    err != nil { return 0, err }
  default:  if err := ( *this.sock ).SetWriteDeadline( time.Now().Add( this.Timeout ) ); err != nil { return 0, err }
  }
  

  rawBuffer  := []byte( fmt.Sprintf( "+%d %s", len( data ), data ) )
  bytesToSend = len( rawBuffer )

  if bytesSent, err = (*this.sock).Write( rawBuffer )                       ; err != nil { return bytesSent, err }
  if bytesSent != bytesToSend                                                            { return bytesSent, errors.New( "Partitial data sent" ) }

  return bytesSent, nil
}


func (this *SQCloud ) readNextRawChunk() ( *Chunk, error ) {
  // every chunk (except RAW JSON) starts with: (<type>)[data]_
  // (_=space)

  NULL     := Chunk{ 0, 0, []byte{ '_', ' ' } }
  chunk    := NULL
  snoop    := []byte{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }

  // Read first byte = Chunk Type
  switch this.Timeout {
  case 0:   if err := ( *this.sock ).SetReadDeadline( time.Time{} );                    err != nil { return &NULL, err }
  default:  if err := ( *this.sock ).SetReadDeadline( time.Now().Add( this.Timeout ) ); err != nil { return &NULL, err }
  }
  switch readCount, err := ( *this.sock ).Read( snoop[ 0 : 1 ] ); {

  case err == io.EOF:         return &NULL, errors.New( "EOF" )
  case err != nil:
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                              return &NULL, errors.New( "Timeout" )
    } else {
                              return &NULL, err
    }
  case readCount < 1:         return &NULL, errors.New( "No Data" )

  case snoop[ 0 ] == '{':     return &NULL, errors.New( "Not implmented" )

  default:
    // Reading second argument (NULL, INT/FLOAT, LEN) until first space
    for tokenLength := 1; tokenLength < len( snoop ); tokenLength++ {

      switch this.Timeout {
      case 0:   if err := ( *this.sock ).SetReadDeadline( time.Time{} );                    err != nil { return &NULL, err }
      default:  if err := ( *this.sock ).SetReadDeadline( time.Now().Add( this.Timeout ) ); err != nil { return &NULL, err }
      }
      switch readCount, err := ( *this.sock ).Read( snoop[ tokenLength : tokenLength + 1 ] ); {

      case err == io.EOF:     return &NULL, errors.New( "EOF" )
      case err != nil:
        if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                              return &NULL, errors.New( "Timeout" )
        } else {
                              return &NULL, err
        }
      case readCount < 1:     return &NULL, errors.New( "No Data" )

      case snoop[ tokenLength ] == ' ': // first space found, raw header = complete

        chunk.DataBufferOffset = 0

        switch snoop[ 0 ] {
        case '_':         // SCSP NULL
          chunk.LEN = 0
          chunk.RAW = snoop[ 0 : tokenLength + 1 ]
          return &chunk, nil

        case ':', ',':    // SCSP Integer, SCSP Float
          chunk.LEN = uint64( tokenLength )
          chunk.RAW = snoop[ 0 : tokenLength + 1 ]
          return &chunk, nil

        case '#': 
        fallthrough

        default:          // all other - except JSON RAW
          var LEN uint64
          if LEN, err = strconv.ParseUint( string( snoop[ 1 : tokenLength ] ), 10, 64 ); err != nil { return &NULL, err }

          tokenLength++
          chunk.DataBufferOffset = uint64( tokenLength )
          chunk.RAW              = make( []byte, uint64( tokenLength ) + LEN )

          // Copy the static snoop buffer into the new dynamic data buffer
          copy( chunk.RAW[ 0 : tokenLength ], snoop[ 0 : tokenLength ] )

          var totalBytesRead uint64 = 0
          for {

            switch this.Timeout {
            case 0:   if err := ( *this.sock ).SetReadDeadline( time.Time{} );                    err != nil { return &NULL, err }
            default:  if err := ( *this.sock ).SetReadDeadline( time.Now().Add( this.Timeout ) ); err != nil { return &NULL, err }
            }
            switch readCount, err := ( *this.sock ).Read( chunk.RAW[ uint64( tokenLength ) + totalBytesRead : ] ); {

            case err == io.EOF:
              totalBytesRead += uint64( readCount )
              if totalBytesRead == LEN {            return &chunk, nil }
                                                    return &NULL, errors.New( "EOF" )
            case err != nil:
              if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                                                    return &NULL, errors.New( "Timeout" )
              } else {
                                                    return &NULL, err
              }
            case totalBytesRead + uint64( readCount ) == LEN:
                                                    chunk.LEN = LEN
                                                    return &chunk, nil
            default:
              totalBytesRead += uint64( readCount )
              time.Sleep( 100 * time.Millisecond ) // wait a moment for the buffers to fill up again...
            }
          }
        }
      }
    }
    return &NULL, errors.New( "Snoop overflow" )
  }
}