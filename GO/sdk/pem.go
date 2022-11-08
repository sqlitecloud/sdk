//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/09/02
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : Default client certificate.
//   ////                ///  ///                     This is a public key.
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

const PEM = `-----BEGIN CERTIFICATE-----
MIID6zCCAtOgAwIBAgIUI0lTm5CfVf3mVP8606CkophcyB4wDQYJKoZIhvcNAQEL
BQAwgYQxCzAJBgNVBAYTAklUMQswCQYDVQQIDAJNTjEQMA4GA1UEBwwHVmlhZGFu
YTEbMBkGA1UECgwSU1FMaXRlIENsb3VkLCBJbmMuMRQwEgYDVQQDDAtTUUxpdGVD
bG91ZDEjMCEGCSqGSIb3DQEJARYUbWFyY29Ac3FsaXRlY2xvdWQuaW8wHhcNMjEw
ODI1MTAwMTI0WhcNMzEwODIzMTAwMTI0WjCBhDELMAkGA1UEBhMCSVQxCzAJBgNV
BAgMAk1OMRAwDgYDVQQHDAdWaWFkYW5hMRswGQYDVQQKDBJTUUxpdGUgQ2xvdWQs
IEluYy4xFDASBgNVBAMMC1NRTGl0ZUNsb3VkMSMwIQYJKoZIhvcNAQkBFhRtYXJj
b0BzcWxpdGVjbG91ZC5pbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
ALnqTqnKgXadcZb4bkHWIrF7BEaPzS8ADUMvrmlP4hVOwg6x4rw33aSTfcXSf6/U
6HzqUgW7lu/Qg/O1WyvdTyseCRbopysPLfU3Hg2bOpcP7ZmgsB3xmPm0tXB/QNvb
sHbMGOGvWVKNCTPuemBuMVLAYNyEC5DWAxOG7IVz+arK2/+QeBH0+PLstVTSvUVy
eu1dacjx9kPEWO0gEwgxyYAeTmgYMRSEcicLF7egxoSS2kzUOLyMkWeV92tP+mzC
NKGgQoG4WnSrsE9ZcY3MiIdb0EnN+nR0VOBFejsTyJm7A+Ab6edEuvNmUTbraKJL
jRKZzUt1r4x3GC+j+UVCQp0CAwEAAaNTMFEwHQYDVR0OBBYEFPGRk8InGXhigfm4
teCLDtYSGu7XMB8GA1UdIwQYMBaAFPGRk8InGXhigfm4teCLDtYSGu7XMA8GA1Ud
EwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAGPk14p4H6vtJsZdgY2sVS2G
T8Ir4ukue79zRFAoCYfkfSW1A+uZbOK/YIjC1CetMXIde1SGvIMcoo1NjrKiqLls
srUN9SmmLihVzurtnoC5ScoUdRQtae8NBWXJnObxK7uXGBYamfs/x0nq1j7DV6Qc
TkuTmJvkcWGcJ9fBOzHgzi+dV+7Y98LP48Pyj/mAzI2icw+I5+DMzn2IktzFf0G7
Sjox3HYOoj2uG2669CLAnw6rkHESbi5imasC9FxWBVxWrnNd0icyiDb1wfBc5W9N
otHL5/wB1MaAmCIcQjIxEshj8pSYTecthitmrneimikFf4KFK0YMvGgKrCLmJsg=
-----END CERTIFICATE-----`

/*
openssl x509 -text -noout -in ../SSL/CA.pem

Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            23:49:53:9b:90:9f:55:fd:e6:54:ff:3a:d3:a0:a4:a2:98:5c:c8:1e
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: C=IT, ST=MN, L=Viadana, O=SQLite Cloud, Inc., CN=SQLiteCloud/emailAddress=marco@sqlitecloud.io
        Validity
            Not Before: Aug 25 10:01:24 2021 GMT
            Not After : Aug 23 10:01:24 2031 GMT
        Subject: C=IT, ST=MN, L=Viadana, O=SQLite Cloud, Inc., CN=SQLiteCloud/emailAddress=marco@sqlitecloud.io
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:b9:ea:4e:a9:ca:81:76:9d:71:96:f8:6e:41:d6:
                    22:b1:7b:04:46:8f:cd:2f:00:0d:43:2f:ae:69:4f:
                    e2:15:4e:c2:0e:b1:e2:bc:37:dd:a4:93:7d:c5:d2:
                    7f:af:d4:e8:7c:ea:52:05:bb:96:ef:d0:83:f3:b5:
                    5b:2b:dd:4f:2b:1e:09:16:e8:a7:2b:0f:2d:f5:37:
                    1e:0d:9b:3a:97:0f:ed:99:a0:b0:1d:f1:98:f9:b4:
                    b5:70:7f:40:db:db:b0:76:cc:18:e1:af:59:52:8d:
                    09:33:ee:7a:60:6e:31:52:c0:60:dc:84:0b:90:d6:
                    03:13:86:ec:85:73:f9:aa:ca:db:ff:90:78:11:f4:
                    f8:f2:ec:b5:54:d2:bd:45:72:7a:ed:5d:69:c8:f1:
                    f6:43:c4:58:ed:20:13:08:31:c9:80:1e:4e:68:18:
                    31:14:84:72:27:0b:17:b7:a0:c6:84:92:da:4c:d4:
                    38:bc:8c:91:67:95:f7:6b:4f:fa:6c:c2:34:a1:a0:
                    42:81:b8:5a:74:ab:b0:4f:59:71:8d:cc:88:87:5b:
                    d0:49:cd:fa:74:74:54:e0:45:7a:3b:13:c8:99:bb:
                    03:e0:1b:e9:e7:44:ba:f3:66:51:36:eb:68:a2:4b:
                    8d:12:99:cd:4b:75:af:8c:77:18:2f:a3:f9:45:42:
                    42:9d
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Subject Key Identifier:
                F1:91:93:C2:27:19:78:62:81:F9:B8:B5:E0:8B:0E:D6:12:1A:EE:D7
            X509v3 Authority Key Identifier:
                keyid:F1:91:93:C2:27:19:78:62:81:F9:B8:B5:E0:8B:0E:D6:12:1A:EE:D7

            X509v3 Basic Constraints: critical
                CA:TRUE
    Signature Algorithm: sha256WithRSAEncryption
         63:e4:d7:8a:78:1f:ab:ed:26:c6:5d:81:8d:ac:55:2d:86:4f:
         c2:2b:e2:e9:2e:7b:bf:73:44:50:28:09:87:e4:7d:25:b5:03:
         eb:99:6c:e2:bf:60:88:c2:d4:27:ad:31:72:1d:7b:54:86:bc:
         83:1c:a2:8d:4d:8e:b2:a2:a8:b9:6c:b2:b5:0d:f5:29:a6:2e:
         28:55:ce:ea:ed:9e:80:b9:49:ca:14:75:14:2d:69:ef:0d:05:
         65:c9:9c:e6:f1:2b:bb:97:18:16:1a:99:fb:3f:c7:49:ea:d6:
         3e:c3:57:a4:1c:4e:4b:93:98:9b:e4:71:61:9c:27:d7:c1:3b:
         31:e0:ce:2f:9d:57:ee:d8:f7:c2:cf:e3:c3:f2:8f:f9:80:cc:
         8d:a2:73:0f:88:e7:e0:cc:ce:7d:88:92:dc:c5:7f:41:bb:4a:
         3a:31:dc:76:0e:a2:3d:ae:1b:6e:ba:f4:22:c0:9f:0e:ab:90:
         71:12:6e:2e:62:99:ab:02:f4:5c:56:05:5c:56:ae:73:5d:d2:
         27:32:88:36:f5:c1:f0:5c:e5:6f:4d:a2:d1:cb:e7:fc:01:d4:
         c6:80:98:22:1c:42:32:31:12:c8:63:f2:94:98:4d:e7:2d:86:
         2b:66:ae:77:a2:9a:29:05:7f:82:85:2b:46:0c:bc:68:0a:ac:
         22:e6:26:c8
*/
