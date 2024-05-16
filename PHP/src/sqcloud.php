<?php
	/*
		MIT License
		
		Copyright (c) 2022-2024 SQLite Cloud, Inc.
		
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

	// v1.1.0:	added new rowset metadata v2
	//			removed ACK for rowset sent in chunk

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

	const ROWSET_CHUNKS_END = '/6 0 0 0 ';

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
		public $notnull = NULL;
		public $prikey = NULL;
		public $autoinc = NULL;
		
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
				print("decltype: ");
				print_r($this->decltype);
				print("\n");

				print("dbname: ");
				print_r($this->dbname);
				print("\n");
				
				print("tblname: ");
				print_r($this->tblname);
				print("\n");
				
				print("origname: ");
				print_r($this->origname);
				print("\n");

				print("notnull: ");
				print_r($this->notnull);
				print("\n");

				print("prikey: ");
				print_r($this->prikey);
				print("\n");

				print("autoinc: ");
				print_r($this->autoinc);
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
		const SDKVersion = '1.1.0';
		
		// User name is required unless connectionstring is provided
		public $username = '';
		// Password is required unless connection string is provided
		public $password = '';
		// Password is hashed
		public $password_hashed = false;
		// API key instead of username and password
		public $apikey = '';

		// Name of database to open
		public $database = '';
		// Optional query timeout passed directly to TLS socket
		public $timeout = NULL;
		// Socket connection timeout
		public $connect_timeout = 20;

		// Enable compression
		public $compression = false;
		// Tell the server to zero-terminate strings
		public $zerotext = false;
		// Database will be created in memory	
		public $memory = false;
		// Create the database if it doesn't exist?
		public $create = false;
		// Request for immediate responses from the server node without waiting for linerizability guarantees
		public $non_linearizable = false;
		// Connect using plain TCP port, without TLS encryption, NOT RECOMMENDED 
		public $insecure = false;
		// Accept invalid TLS certificates
		public $no_verify_certificate = false;
		
		// Certificates 
		public $tls_root_certificate = NULL;
		public $tls_certificate = NULL;
		public $tls_certificate_key = NULL;
		
		// Server should send BLOB columns
		public $noblob = false;
		// Do not send columns with more than max_data bytes
		public $maxdata = 0;
		// Server should chunk responses with more than maxRows
		public $maxrows = 0;
		// Server should limit total number of rows in a set to maxRowset
		public $maxrowset = 0;
		
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
				if ($this->no_verify_certificate) {
					stream_context_set_option($context, 'ssl', 'verify_peer ', false);
					stream_context_set_option($context, 'ssl', 'verify_peer_name ', false);
				}
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

		public function connectWithString ($connectionString) {
			// URL STRING FORMAT
    		// sqlitecloud://user:pass@host.com:port/dbname?timeout=10&key2=value2&key3=value3
			// or sqlitecloud://host.sqlite.cloud:8860/dbname?apikey=zIiAARzKm9XBVllbAzkB1wqrgijJ3Gx0X5z1A4m4xBA

			$params = parse_url($connectionString);
			if (!is_array($params)) {
				$this->errmsg = "Invalid connection string: {$connectionString}.";
				$this->errcode = -1;
				return false;
			}

			$options = [];
			$query = isset($params['query']) ? $params['query'] : '';
			parse_str($query, $options);
			foreach ($options as $option => $value) {
				$opt = strtolower($option);

				// prefix for certificate options
				if (strcmp($opt, "root_certificate") == 0 
				|| strcmp($opt, "certificate") == 0 
				|| strcmp($opt, "certificate_key") == 0) {
					$opt = "tls_" . $opt;
				}

				if (property_exists($this, $opt)) {
					if (filter_var($value, FILTER_VALIDATE_BOOLEAN, FILTER_NULL_ON_FAILURE) !== null) {
						$this->{$opt} = (bool) ($value);
					} else if (is_numeric($value)) {
						$this->{$opt} = (int) ($value);
					} else {
						$this->{$opt} = $value;
					}
				}
			}
			
			// apikey or username/password is accepted
			if (!$this->apikey) {
				$this->username = isset($params['user']) ? urldecode($params['user']) : '';
				$this->password = isset($params['pass']) ?  urldecode($params['pass']) : '';
			}
			
			$path = isset($params['path']) ? $params['path'] : '';
			$database = str_replace('/', '', $path);
			if ($database) {
				$this->database = $database;
			}
			
			$hostname = $params['host'];
			$port = isset($params['port']) ? (int)($params['port']) : null;

			if ($port) {
				return $this->connect($hostname, $port);
			}

			return $this->connect($hostname);
			
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
			
			if ($this->apikey) {
				$buffer .= "AUTH APIKEY {$this->apikey};";
			} 
			
			if ($this->username && $this->password) {
				$command = $this->password_hashed ? 'HASH' : 'PASSWORD';
				$buffer .= "AUTH USER {$this->username} {$command} {$this->password};";
			}
			
			if ($this->database) {
				if ($this->create && !$this->memory) {
					$buffer .= "CREATE DATABASE {$this->database} IF NOT EXISTS;";
				}
				$buffer .= "USE DATABASE {$this->database};";
			}
			
			if ($this->compression) {
				$buffer .= "SET CLIENT KEY COMPRESSION TO 1;";
			}
			
			if ($this->zerotext) {
				$buffer .= "SET CLIENT KEY ZEROTEXT TO 1;";
			}

			if ($this->non_linearizable) {
				$buffer .= "SET CLIENT KEY NONLINEARIZABLE TO 1;";
			}
		
			if ($this->noblob) {
				$buffer .= "SET CLIENT KEY NOBLOB TO 1;";
			}

			if ($this->maxdata) {
				$buffer .= "SET CLIENT KEY MAXDATA TO {$this->maxdata};";
			}

			if ($this->maxrows) {
				$buffer .= "SET CLIENT KEY MAXROWS TO {$this->maxrows};";
			}

			if ($this->maxrowset) {
				$buffer .= "SET CLIENT KEY MAXROWSET TO {$this->maxrowset};";	
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
				} elseif ($c == CMD_ROWSET_CHUNK) {
					// chunkes are completed when the buffer contains the end-of-chunk marker
					$isEndOfChunk = substr($buffer, -strlen(ROWSET_CHUNKS_END)) == ROWSET_CHUNKS_END;
					if (!$isEndOfChunk) continue;
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
		
		private function internal_uncompress_data ($buffer) {
			// %LEN COMPRESSED UNCOMPRESSED BUFFER
						
			// extract compressed size
			$space_index = strpos($buffer, ' ');
			$buffer = substr($buffer, $space_index + 1);

			// extract compressed size
			$space_index = strpos($buffer, ' ');
			$compressed_size = intval(substr($buffer, 0, $space_index));
			$buffer = substr($buffer, $space_index + 1);

			// extract decompressed size
			$space_index = strpos($buffer, ' ');
			$uncompressed_size = intval(substr($buffer, 0, $space_index));
			$buffer = substr($buffer, $space_index + 1);

			// extract data header
			$header = substr($buffer, 0, -$compressed_size);

			// extract compressed data
			$compressed_buffer = substr($buffer, -$compressed_size);

			$decompressed_buffer = $header . $this->lz4decode($compressed_buffer, 0);

			// sanity check result
			if (strlen($decompressed_buffer) != $uncompressed_size + strlen($header)) {
				return NULL;
			}

			return $decompressed_buffer;
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
				$buffer = $this->internal_uncompress_data ($buffer);
				if ($buffer == NULL) {
					$this->errcode = -1;
					$this->errmsg = "An error occurred while decompressing the input buffer of len {$blen}.";
					return false;
				}
				// after decompression length has changed
				$blen = strlen($buffer);
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
					// - When decompressed, LEN for ROWSET is *0
					//
            		// CMD_ROWSET_CHUNK:    /LEN IDX:VERSION ROWS COLS DATA
					//
        			$start = $this->internal_parse_rowset_signature($buffer, $len, $idx, $version, $nrows, $ncols);
					if ($start < 0) return false;
					
					// check for end-of-chunk condition
					if ($start == 0 && $version == 0) {
						$rowset = $this->rowset;
						$this->rowset = NULL;
						return $rowset;
					}
        			
        			$rowset = $this->internal_parse_rowset($buffer, $start, $idx, $version, $nrows, $ncols);

					// continue parsing next chunk in the buffer
					if ($buffer[0] == CMD_ROWSET_CHUNK) {
						$buffer = substr($buffer, $len + strlen("/{$len} "));
						if ($buffer) {
							return $this->internal_parse_buffer($buffer, strlen($buffer));
						}
					}

					return $rowset;
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
			if ($buffer == ROWSET_CHUNKS_END) {
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
			
			// parse column names (header is guarantee to contain column names)
			$rowset->colname = array();
			for ($i = 0; $i < $ncols; $i++) {
				$len = $this->internal_parse_number($buffer, $cstart, $unused, $start);
				$value = substr($buffer, $cstart, $len);
				array_push($rowset->colname, $value);
				$start = $cstart + $len;
			}
			
			if ($rowset->version == 1) return $start;
			
			// if version != 2 returns an error because rowset version is not supported
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

			// parse not null flags
			$rowset->notnull = array();
			for ($i = 0; $i < $ncols; $i++) {
				$value = $this->internal_parse_number($buffer, $cstart, $unused, $start);
				array_push($rowset->notnull, $value);
				$start = $cstart;
			}

			// parse primary key flags
			$rowset->prikey = array();
			for ($i = 0; $i < $ncols; $i++) {
				$value = $this->internal_parse_number($buffer, $cstart, $unused, $start);
				array_push($rowset->prikey, $value);
				$start = $cstart;
			}

			// parse autoincrement flags
			$rowset->autoinc = array();
			for ($i = 0; $i < $ncols; $i++) {
				$value = $this->internal_parse_number($buffer, $cstart, $unused, $start);
				array_push($rowset->autoinc, $value);
				$start = $cstart;
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

			return $rowset;
		}
		
		// MARK: -
		
		function __destruct() {
        	$this->disconnect();
    	}
	}

?>