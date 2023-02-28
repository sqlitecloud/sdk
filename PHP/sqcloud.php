<?php
	/*
		MIT License
		
		Copyright (c) 2022 SQLite Cloud, Inc.
		
		Permission is hereby granted, free of charge, to any person obtaining a copy
		of this software and associated documentation files (the "Software"), to deal
		in the Software without restriction, including without limitation the rights
		to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
		copies of the Software, and to permit persons to whom the Software is
		furnished to do so, subject to the following conditions:
		
		The above copyright notice and this permission notice shall be included in all
		copies or substantial portions of the Software.
		
		THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
		IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
		FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
		AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
		LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
		OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
		SOFTWARE.
	*/

	// namespace SQCloud;
	
	const CMD_STRING		=	'+';
	const CMD_ZEROSTRING	=	'!';
	const CMD_ERROR			=	'-';
	const CMD_INT			=	':';
	const CMD_FLOAT			=	',';
	const CMD_ROWSET		=	'*';
	const CMD_ROWSET_CHUNK	=	'/';
	const CMD_JSON			=	'#';
	const CMD_RAWJSON		=	'{';
	const CMD_NULL			=	'_';
	const CMD_BLOB			=	'$';
	const CMD_COMPRESSED	=	'%';
	const CMD_PUBSUB		=	'|';
	const CMD_COMMAND		=	'^';
	const CMD_RECONNECT		=	'@';
	const CMD_ARRAY			=	'=';
	
	class SQLiteCloudRowset {
		public $nrows = 0;
		public $ncols = 0;
		public $version = 0;
		public $data = NULL;
		
		// version 2 only
		public $colname = NULL;
		public $decltype = NULL;
		public $dbname = NULL;
		public $tblname = NULL;
		public $origname = NULL;
		
		private function compute_index ($row, $col) {
			if ($row < 0 || $row >= $this->nrows) return -1;
			if ($col < 0 || $col >= $this->ncols) return -1;
			return $row*$this->ncols+$col;
		}
		
		public function value ($row, $col) {
			$index = $this->compute_index($row, $col);
			if ($index < 0) return NULL;
			return $this->data[$index];
		}
		
		public function name ($col) {
			if ($col < 0 || $col >= $this->ncols) return NULL;
			return $this->colname[$col];
		}
		
		public function dump () {
			print("version: {$this->version}\n");
			print("nrows: {$this->nrows}\n");
			print("ncols: {$this->ncols}\n");
			
			print("colname: ");
			print_r($this->colname);
			print("\n");
			
			if ($this->version == 2) {
				print("dbname: ");
				print_r($this->dbname);
				print("\n");
				
				print("tblname: ");
				print_r($this->tblname);
				print("\n");
				
				print("origname: ");
				print_r($this->origname);
				print("\n");
			}
			
			if ($this->data && count($this->data) > 0) {
				print("data: ");
				print_r($this->data);
				print("\n");
			}
		}
	}
	
	class SQLiteCloud {
		public const SDKVersion = '1.0.0';
		
		public $username = '';
		public $password = '';
		public $database = '';
		public $timeout = NULL;
		public $connect_timeout = 20;
		public $compression = false;
		public $sqlitemode = false;
		public $zerotext = false;	
		public $insecure = false;
		public $tls_root_certificate = NULL;
		public $tls_certificate = NULL;
		public $tls_certificate_key = NULL;
		
		public $errmsg = NULL;
		public $errcode = 0;
		public $xerrcode = 0;
		
		private $socket = NULL;
		private $isblob = false;
		private $rowset = NULL;
		
		// PUBLIC
		public function connect ($hostname = "localhost", $port = 8860) {
			$ctx = ($this->insecure) ? 'tcp' : 'tls';
			$address = "{$ctx}://{$hostname}:{$port}";
			
			// check setup context for TLS connection
			$context = NULL;
			if (!$this->insecure) {
				$context = stream_context_create();
				if ($this->tls_root_certificate) stream_context_set_option($context, 'ssl', 'cafile', $this->tls_root_certificate);
				if ($this->tls_certificate) stream_context_set_option($context, 'ssl', 'local_cert', $this->tls_certificate);
				if ($this->tls_certificate_key) stream_context_set_option($context, 'ssl', 'local_pk', $this->tls_certificate_key);
			}
			
			// connect to remote socket
			$socket = stream_socket_client($address, $this->errcode, $this->errmsg, $this->connect_timeout, STREAM_CLIENT_CONNECT, $context);
			if (!$socket) {
				if ($this->errcode == 0) {
					// if the value returned in errcode is 0 and stream_socket_client returned false, it is an indication
					// that the error occurred before the connect() call. This is most likely due to a problem initializing
					// the socket
					$extmsg = ($this->insecure) ? '(before connecting to remote host)' : '(possibly wrong TLS certificate)';
					$this->errmsg = "An error occurred while initializing the socket {$extmsg}.";
					$this->errcode = -1;
				}
				return false;
			}
			
			$this->socket = $socket;
			if ($this->internal_config_apply() == false) return false;
			
			return true;
		}
		
		public function disconnect () {
			$this->internal_clear_error();
			if ($this->socket) fclose($this->socket);
			$this->socket = NULL;
		}
		
		public function execute ($command) {
			return $this->internal_run_command($command);
		}
		
		public function sendblob ($blob) {
			$this->isblob = true;
			$rc = $this->internal_run_command($blob);
			$this->isblob = false;
			return $rc;
		}
		
		// MARK: -
		
		// PRIVATE
		
		// lz4decode function from http://heap.ch/blog/2019/05/18/lz4-decompression/
		/*
			MIT License
			
			Copyright (c) 2019 Stephan J. Müller
			
			Permission is hereby granted, free of charge, to any person obtaining a copy
			of this software and associated documentation files (the "Software"), to deal
			in the Software without restriction, including without limitation the rights
			to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
			copies of the Software, and to permit persons to whom the Software is
			furnished to do so, subject to the following conditions:
			
			The above copyright notice and this permission notice shall be included in all
			copies or substantial portions of the Software.
			
			THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
			IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
			FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
			AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
			LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
			OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
			SOFTWARE.
		*/
		private function lz4decode($in, $offset = 0, $header = '') {
		  $len = strlen($in);
		  $out = $header;
		  $i = $offset;
		  $take = function() use ($in, &$i) {
			return ord($in[$i++]);
		  };
		  $addOverflow = function(&$sum) use ($take) {
			do {
			  $sum += $summand = $take();
			} while ($summand === 0xFF);
		  };
		  while ($i < $len) {
			$token = $take();
			$nLiterals = $token >> 4;
			if ($nLiterals === 0xF) $addOverflow($nLiterals);
			$out .= substr($in, $i, $nLiterals);
			$i += $nLiterals;
			if ($i === $len) break;
			$offset = $take() | $take() << 8;
			$matchlength = $token & 0xF;
			if ($matchlength === 0xF) $addOverflow($matchlength);
			$matchlength += 4;
			$j = strlen($out) - $offset;
			while ($matchlength--) {
			  $out .= $out[$j++];
			}
		  }
		  return $out;
		}
		
		private function internal_config_apply () {
			if ($this->timeout > 0) stream_set_timeout($this->socket, $this->timeout);
			
			$buffer = '';
			if ((strlen($this->username) > 0) && (strlen($this->password) > 0)) {
				$buffer .= "AUTH USER {$this->username} PASSWORD {$this->password};";
			}
			
			if (strlen($this->database) > 0) {
				$buffer .= "USE DATABASE {$this->database};";
			}
			
			if ($this->compression) {
				$buffer .= "SET CLIENT KEY COMPRESSION TO 1;";
			}
			
			if ($this->sqlitemode) {
				$buffer .= "SET CLIENT KEY SQLITE TO 1;";
			}
			
			if ($this->zerotext) {
				$buffer .= "SET CLIENT KEY ZEROTEXT TO 1;";
			}
			
			if (strlen($buffer) > 0) {
				$result = $this->internal_run_command($buffer);
				if ($result === false) return false;
			}
			
			return true;
		}
		
		private function internal_run_command ($buffer) {
			$this->internal_clear_error();
			
			if ($this->internal_socket_write($buffer) === false) return false;
			return $this->internal_socket_read();
		}
		
		private function internal_setup_pubsub ($buffer) {
			return true;
		}
		
		private function internal_reconnect ($buffer) {
			return true;
		}
		
		private function internal_parse_array ($buffer) {
			// extract the number of values in the array
			$start = 0;
    		$n = $this->internal_parse_number($buffer, $start, $unused, 0);
    		
    		// loop to parse each individual value
			$r = array();
			for ($i=0; $i < $n; ++$i) {
				$cellsize = 0;
				$len = strlen($buffer) - $start;
				$value = $this->internal_parse_value($buffer, $len, $cellsize, $start);
				$start += $cellsize;
				array_push($r, $value);
			}
			
			return $r;
		}
		
		private function internal_clear_error () {
			$this->errmsg = NULL;
			$this->errcode = 0;
			$this->xerrcode = 0;
		}
		
		private function internal_socket_write ($buffer) {
			// compute header
			$delimit = ($this->isblob) ? '$' : '+';
			$len = ($buffer) ? strlen($buffer) : 0;
			$header = "{$delimit}{$len} ";
			
			// write header and buffer
			if (fwrite($this->socket, $header) === false) return false;
			if ($len == 0) return true;
			if (fwrite($this->socket, $buffer) === false) return false;
			
			return true;
		}
		
		private function internal_socket_read () {
			$buffer = "";
			$len = 8192;
			
			$nread = 0;
			while (true) {
				// read from socket
				$temp = fread($this->socket, $len);
				if ($temp === false) return false;
								
				// update buffers
				$buffer .= $temp;
				$nread += strlen($temp);
				
				// get first character
				$c = $buffer[0];
				
				// check if command does not have an explicit length
				if (($c == CMD_INT) || ($c == CMD_FLOAT) || ($c == CMD_NULL)) {
					// command is terminated by a space character
					if ($buffer[$nread-1] != ' ') continue;
				} else {
					$cstart = 0;
					$n = $this->internal_parse_number($buffer, $cstart);
					
					$can_be_zerolength = ($c == CMD_BLOB) || ($c == CMD_STRING);
					if ($n == 0 && !$can_be_zerolength) continue;
					
					// check exit condition
					if ($n + $cstart != $nread) continue;
				}
				
				return $this->internal_parse_buffer($buffer, $nread);
			}
			
			return false;
		}
		
		private function internal_uncompress_data ($buffer, $blen) {
			// %LEN COMPRESSED UNCOMPRESSED BUFFER
			
			$tlen = 0;	// total length
			$clen = 0;	// compressed length
			$ulen = 0;	// uncompressed length
			$hlen = 0;	// raw header length
			$seek1 = 0;
			
			$start = 1;
			$counter = 0;
			for ($i = 0; $i < $blen; $i++) {
				if ($buffer[$i] != ' ') continue;
				++$counter;
				
				$data = substr($buffer, $start, $i-$start);
				$start = $i + 1;
				
				if ($counter == 1) {
					$tlen = intval($data);
					$seek1 = $start;
				}
				else if ($counter == 2) {
					$clen = intval($data);
				}
				else if ($counter == 3) {
					$ulen = intval($data);
					break;
				}
			}
			
			// sanity check header values
			if ($tlen == 0 || $clen == 0 || $ulen == 0 || $start == 1 || $seek1 == 0) return NULL;
			
			// copy raw header
			$hlen = $start - $seek1;
			$header = substr($buffer, $start, $hlen);
			
			// compute index of the first compressed byte
			$start += $hlen;
			
			// perform real decompression in pure PHP code
			$clone = $this->lz4decode($buffer, $start, $header);
			
			// sanity check result
			if (strlen($clone) != $ulen + $hlen) return NULL;
			
			return $clone;
		}
		
		private function internal_parse_value ($buffer, &$len, &$cellsize = NULL, $index = 0) {
			if ($len <= 0) return NULL;
			
			// handle special NULL value case
			if (is_null($buffer) || $buffer[$index] == CMD_NULL) {
				$len = 0;
				if (!is_null($cellsize)) $cellsize = 2;
				return NULL;
			}
			
			$cstart = $index;
			$blen = $this->internal_parse_number($buffer, $cstart, $unused, $index+1);
			
			// handle decimal/float cases
			if (($buffer[$index] == CMD_INT) || ($buffer[$index] == CMD_FLOAT)) {
				$nlen = $cstart - $index;
        		$len = $nlen - 2;
        		if (!is_null($cellsize)) $cellsize = $nlen;
				return substr($buffer, $index+1, $len);
			}
			
			$len = ($buffer[$index] == CMD_ZEROSTRING) ? $blen - 1 : $blen;
    		if (!is_null($cellsize)) $cellsize = $blen + $cstart - $index;
    		
    		return substr($buffer, $cstart, $len);
		}
		
		private function internal_parse_buffer ($buffer, $blen) {
			// possible return values:
			// true 	=> OK
			// false 	=> error
			// integer
			// double
			// string
			// array
			// object
			// NULL
			
			// check OK value
			if (strcmp($buffer, '+2 OK') == 0) return true;
		
			// check for compressed result
			if ($buffer[0] == CMD_COMPRESSED) {
				$buffer = $this->internal_uncompress_data ($buffer, $blen);
				if ($buffer == NULL) {
					$this->errcode = -1;
					$this->errmsg = 'An error occurred while decompressing the input buffer of len {$len}.';
					return false;
				}
			}
			
			// first character contains command type
			switch ($buffer[0]) {
				case CMD_ZEROSTRING:
        		case CMD_RECONNECT:
        		case CMD_PUBSUB:
        		case CMD_COMMAND:
        		case CMD_STRING:
        		case CMD_ARRAY:
        		case CMD_BLOB:
        		case CMD_JSON: {
        			$cstart = 0;
        			$len = $this->internal_parse_number($buffer, $cstart);
        			if ($len == 0) return "";
        			
        			if ($buffer[0] == CMD_ZEROSTRING) --$len;
        			$clone = substr($buffer, $cstart, $len);
        		
        			if ($buffer[0] == CMD_COMMAND) return $this->internal_run_command($clone);
            		else if ($buffer[0] == CMD_PUBSUB) return $this->internal_setup_pubsub($clone);
					else if ($buffer[0] == CMD_RECONNECT) return $this->internal_reconnect($clone);
					else if ($buffer[0] == CMD_ARRAY) return $this->internal_parse_array($clone);
            
        			return $clone;
        		}
        		
        		case CMD_ERROR: {
        			// -LEN ERRCODE:EXTCODE ERRMSG
        			$cstart = 0; $cstart2 = 0;
        			$len = $this->internal_parse_number($buffer, $cstart);
        			$clone = substr($buffer, $cstart);
        			
        			$extcode = 0;
        			$errcode = $this->internal_parse_number($clone, $cstart2, $extcode, 0);
        			$this->errcode = $errcode;
        			$this->xerrcode = $extcode;
        			
        			$len -= $cstart2;
        			$this->errmsg = substr($clone, $cstart2);
        			
        			return false;
        		}
        			
        		case CMD_ROWSET:
        		case CMD_ROWSET_CHUNK: {
        			// CMD_ROWSET:          *LEN 0:VERSION ROWS COLS DATA
            		// CMD_ROWSET_CHUNK:    /LEN IDX:VERSION ROWS COLS DATA
        			$start = $this->internal_parse_rowset_signature($buffer, $len, $idx, $version, $nrows, $ncols);
					if ($start < 0) return false;
					
					// check for end-of-chunk condition
					if ($start == 0 && $version == 0) {
						$rowset = $this->rowset;
						$this->rowset = NULL;
						return $rowset;
					}
        			
					// continue parsing
        			return $this->internal_parse_rowset($buffer, $start, $idx, $version, $nrows, $ncols);
        		}
        			
        		case CMD_NULL:
        			return NULL;
        			
        		case CMD_INT:
        		case CMD_FLOAT: {
        			$clone = $this->internal_parse_value($buffer, $blen);
        			if (is_null($clone)) return 0;
        			if ($buffer[0] == CMD_INT) return intval($clone);
        			return floatval($clone);
        		}
        			
        		case CMD_RAWJSON:
        			return NULL;
			}
			
			return NULL;
		}
		
		private function internal_parse_number ($buffer, &$cstart, &$extcode = NULL, $index = 1) {
			$value = 0;
			$extvalue = 0;
			$isext = false;
			$blen = strlen($buffer);
			
			// from 1 to skip the first command type character
			for ($i = $index; $i < $blen; $i++) {
				$c = $buffer[$i];
				
				// check for optional extended error code (ERRCODE:EXTERRCODE)
        		if ($c == ':') {$isext = true; continue;}
        		
        		// check for end of value
				if ($c == ' ') {
					$cstart = $i + 1;
            		if (!is_null($extcode)) $extcode = $extvalue;
					return $value;
        		}
        		
        		// compute numeric value
        		if ($isext) $extvalue = ($extvalue * 10) + ((int)$buffer[$i]);
        		else $value = ($value * 10) + ((int)$buffer[$i]);
			}
			
			return 0;
		}
		
		// MARK: -
		
		function internal_parse_rowset_signature ($buffer, &$len, &$idx, &$version, &$nrows, &$ncols) {
			// ROWSET:          *LEN 0:VERS NROWS NCOLS DATA
			// ROWSET in CHUNK: /LEN IDX:VERS NROWS NCOLS DATA
			
			// check for end-of-chunk condition
			if ($buffer == '/6 0 0 0 ') {
				$version = 0;
				return 0;
			}
			
			$start = 1;
			$counter = 0;
			$n = strlen($buffer);
			for ($i = 0; $i < $n; $i++) {
				if ($buffer[$i] != ' ') continue;
				++$counter;
				
				$data = substr($buffer, $start, $i-$start);
				$start = $i + 1;
				
				if ($counter == 1) {
					$len = intval($data);
				}
				else if ($counter == 2) {
					// idx:vers
					$values = explode(":", $data);
					$idx = intval($values[0]);
					$version = intval($values[1]);
				}
				else if ($counter == 3) {
					$nrows = intval($data);
				}
				else if ($counter == 4) {
					$ncols = intval($data);
					return $start;
				}
				else return -1;
			}
			return -1;
		}
		
		function internal_parse_rowset_header ($rowset, $buffer, $start) {
			$ncols = $rowset->ncols;
			
			// parse column names
			$rowset->colname = array();
			for ($i = 0; $i < $ncols; $i++) {
				$len = $this->internal_parse_number($buffer, $cstart, $unused, $start);
				$value = substr($buffer, $cstart, $len);
				array_push($rowset->colname, $value);
				$start = $cstart + $len;
			}
			
			if ($rowset->version == 1) return $start;
			
			// if version != 2 returns an error
			if ($rowset->version != 2) return -1;
			
			// parse declared types
			$rowset->decltype = array();
			for ($i = 0; $i < $ncols; $i++) {
				$len = $this->internal_parse_number($buffer, $cstart, $unused, $start);
				$value = substr($buffer, $cstart, $len);
				array_push($rowset->decltype, $value);
				$start = $cstart + $len;
			}
			
			// parse database names
			$rowset->dbname = array();
			for ($i = 0; $i < $ncols; $i++) {
				$len = $this->internal_parse_number($buffer, $cstart, $unused, $start);
				$value = substr($buffer, $cstart, $len);
				array_push($rowset->dbname, $value);
				$start = $cstart + $len;
			}
			
			// parse table names
			$rowset->tblname = array();
			for ($i = 0; $i < $ncols; $i++) {
				$len = $this->internal_parse_number($buffer, $cstart, $unused, $start);
				$value = substr($buffer, $cstart, $len);
				array_push($rowset->tblname, $value);
				$start = $cstart + $len;
			}
			
			// parse column original names
			$rowset->origname = array();
			for ($i = 0; $i < $ncols; $i++) {
				$len = $this->internal_parse_number($buffer, $cstart, $unused, $start);
				$value = substr($buffer, $cstart, $len);
				array_push($rowset->origname, $value);
				$start = $cstart + $len;
			}
			
			return $start;
		}
		
		function internal_parse_rowset_values ($rowset, $buffer, $start, $bound) {
    		// loop to parse each individual value
			for ($i=0; $i < $bound; ++$i) {
				$cellsize = 0;
				$len = strlen($buffer) - $start;
				$value = $this->internal_parse_value($buffer, $len, $cellsize, $start);
				$start += $cellsize;
				array_push($rowset->data, $value);
			}
		}
		
		function internal_parse_rowset($buffer, $start, $idx, $version, $nrows, $ncols) {
			$rowset = NULL;
			$n = $start;
			$ischunk = ($buffer[0] == CMD_ROWSET_CHUNK);
			
			// idx == 0 means first (and only) chunk for rowset
			// idx == 1 means first chunk for chunked rowset
			$first_chunk = ($ischunk) ? ($idx == 1) : ($idx == 0);
			if ($first_chunk) {
				$rowset = new SQLiteCloudRowset();
				$rowset->nrows = $nrows;
				$rowset->ncols = $ncols;
				$rowset->version = $version;
				$rowset->data = array();
				if ($ischunk) $this->rowset = $rowset;
				$n = $this->internal_parse_rowset_header($rowset, $buffer, $start);
				if ($n <= 0) return NULL;
			} else {
				$rowset = $this->rowset;
				$rowset->nrows += $nrows;
			}
			
			// parse values
			$this->internal_parse_rowset_values($rowset, $buffer, $n, $nrows*$ncols);
			
			if ($ischunk) {
				if ($this->internal_socket_write("OK") === false) {
					$this->rowset = NULL;
					return false;
				}
				return $this->internal_socket_read();
			}
			
			return $rowset;
		}
		
		// MARK: -
		
		function __destruct() {
        	$this->disconnect();
    	}
	}

?>