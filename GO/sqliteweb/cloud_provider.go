//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
//     //             ///   ///  ///    Date        : 2023/01/24
//    ///             ///   ///  ///    Author      : Andreas Donetti
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description :
//   ////                ///  ///
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	sqlitecloud "github.com/sqlitecloud/sdk"
	"github.com/teris-io/shortid"
)

const (
	pollingSleep   = 5 * time.Second
	pollingTimeout = 10 * time.Minute
)

func createNode(userid int, name string, region string, size string, nodetype string, projectuuid string, nodeid int) (string, error) {
	nodeshortid := ""
	nid := int64(0)
	nattempts := 0
	maxretry := 100

	shortidgen, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return "", fmt.Errorf("Cannot create the ShortId generator: %s", err.Error())
	}

	for {
		var sid string = ""
		var err error = nil

		for {
			sid, err = shortidgen.Generate()
			if err != nil {
				return "", fmt.Errorf("Cannot generate new shortid: %s", err.Error())
			}

			// the shortuuid is used as the hostname for the DNS record
			// hostname labels may contain only the ASCII letters 'a' through 'z' (in a case-insensitive manner),
			// the digits '0' through '9', and the hyphen ('-')
			sid = strings.ToLower(sid)
			if validateShortUUID(sid) {
				break
			}
		}

		nattempts += 1
		sql := fmt.Sprintf("INSERT INTO Node (project_uuid, node_id, shortuuid, name, created) VALUES ('%s', %d, '%s', '%s', 0) RETURNING id;", projectuuid, nodeid, sid, name)
		res, err, errcode, _, _ := dashboardcm.ExecuteSQL("auth", sql)
		if err == nil && res.GetNumberOfColumns() == 1 && res.GetNumberOfRows() == 1 {
			nodeshortid = sid
			nid = res.GetInt64Value_(0, 0)
			break
		}

		// (19) SQLITE_CONSTRAINT
		if errcode != 19 {
			return "", fmt.Errorf("Cannot insert new node (%d/%d) %s %s %s: %s (%d)", nattempts, 5, projectuuid, name, sid, err.Error(), errcode)
		}

		if nattempts >= maxretry {
			return "", fmt.Errorf("Cannot insert new node (attempts %d/%d) %s %s %s: %s", nattempts, maxretry, projectuuid, name, sid, err.Error())
		}
	}

	jobuuid := uuid.New().String()
	sql := fmt.Sprintf("INSERT INTO Jobs (uuid, name, status, steps, progress, user_id, node_id) VALUES ('%s', 'Create Node %s', 'Creating droplet', 2, 0, %d, %d)", jobuuid, name, userid, nid) // ; SELECT uuid FROM Jobs WHERE rowid = last_insert_rowid();
	_, err, _, _, _ = dashboardcm.ExecuteSQL("auth", sql)
	if err != nil {
		return "", fmt.Errorf("Cannot create the job %s: %s", jobuuid, err.Error())
	}

	cloudProvider := digitalocean
	sqlitecloudPort := 8860

	go func() {
		clusterConfig, nnodes, err := clusterConfig(projectuuid, nodeid, fmt.Sprintf("%s.%s", nodeshortid, dropletDomain), sqlitecloudPort)
		if err != nil {
			SQLiteWeb.Logger.Error(err.Error())
			sql := fmt.Sprintf("UPDATE Jobs SET status = '%s' WHERE uuid = '%s'", err.Error(), jobuuid)
			authExecSQL(sql)
			return
		}

		// run CreateNode asynchronously in a goroutine
		nodeCreateReq := &CloudNodeCreateRequest{
			JobUUID:       jobuuid,
			Name:          name,
			Region:        region,
			Size:          size,
			Type:          nodetype,
			ProjectUUID:   projectuuid,
			NodeID:        nodeid,
			Hostname:      nodeshortid,
			Domain:        dropletDomain,
			Port:          sqlitecloudPort,
			ClusterConfig: clusterConfig,
		}

		cloudNode, err := cloudProvider.CreateNode(nodeCreateReq)
		if err != nil {
			SQLiteWeb.Logger.Errorf("digitalocean CreateNode: %s", err.Error())
			sql := fmt.Sprintf("UPDATE Jobs SET status = '%s' WHERE uuid = '%s'", err.Error(), jobuuid)
			authExecSQL(sql)
			return
		}
		SQLiteWeb.Logger.Debugf("Droplet completed %s %s %d", cloudNode.Name, cloudNode.JobUUID, cloudNode.DropletID)

		if nnodes > 0 {
			// execute the ADD NODE command to the cluster
			addCommand := fmt.Sprintf("ADD NODE %d ADDRESS %s:%d", cloudNode.NodeID, cloudNode.AddrV4, cloudNode.Port)
			_, err, _, _, _ = dashboardcm.ExecuteSQL(projectuuid, addCommand)
			if err != nil {
				SQLiteWeb.Logger.Errorf("Cannot add the new node %d to the cluster %s: %s", cloudNode.NodeID, cloudNode.JobUUID, err.Error())
			}
			SQLiteWeb.Logger.Debugf("Node added to the cluster %d %s", cloudNode.NodeID, cloudNode.JobUUID)
		}

		sql = fmt.Sprintf("UPDATE Jobs SET progress = progress + 1, status = 'System setup' WHERE uuid = '%s'", nodeCreateReq.JobUUID)
		authExecSQL(sql)

		sql = fmt.Sprintf("UPDATE Node SET hostname = '%s', type = '%s', provider = '%s', image = '%s', region = '%s', addr4 = '%s', addr6 = '%s', port = %d, latitude = %f, longitude = %f, created = 1 WHERE id = %d", cloudNode.FullyQualifiedDomainName(), cloudNode.Type, "DigitalOcean", cloudNode.Size, cloudNode.Region, cloudNode.AddrV4, cloudNode.AddrV6, cloudNode.Port, cloudNode.Location.Latitude, cloudNode.Location.Longitude, nid)
		authExecSQL(sql)

		conn, err := waitForConnection(cloudNode)
		if err != nil {
			SQLiteWeb.Logger.Error(err.Error())
			sql := fmt.Sprintf("UPDATE Jobs SET status = '%s' WHERE uuid = '%s'", err.Error(), jobuuid)
			authExecSQL(sql)
			return
		}

		if nnodes == 0 {
			res, err, _, _, _ := dashboardcm.ExecuteSQL("auth", "SELECT admin_username, admin_password FROM Project")
			if err != nil {
				SQLiteWeb.Logger.Errorf("Cannot get admin credentials %s", err.Error())
			}
			if v := res.GetStringValue_(0, 0); v != "" && v != "admin" {
				conn.ExecuteArray("RENAME USER admin TO ?", []interface{}{v})
			}
			if v := res.GetStringValue_(0, 1); v != "" && v != "admin" {
				conn.ExecuteArray("SET MY PASSWORD ?", []interface{}{v})
			}
		}

		sql = fmt.Sprintf("UPDATE Jobs SET progress = progress + 1, status = 'completed' WHERE uuid = '%s'", nodeCreateReq.JobUUID)
		authExecSQL(sql)

		sql = fmt.Sprintf("UPDATE Node SET created = 1 WHERE id = %d", nid)
		authExecSQL(sql)
	}()

	return jobuuid, nil
}

func authExecSQL(sql string) {
	_, err, _, _, _ := dashboardcm.ExecuteSQL("auth", sql)
	if err != nil {
		SQLiteWeb.Logger.Errorf("Error '%s' for query '%s'", err.Error(), sql)
	}
}

func waitForConnection(cloudNode *CloudNode) (conn *sqlitecloud.SQCloud, err error) {
	timeout := time.NewTimer(pollingTimeout)
	for {
		connectionString := fmt.Sprintf("sqlitecloud://admin:admin@%s:%d?timeout=5", cloudNode.FullyQualifiedDomainName(), cloudNode.Port)
		conn, err = sqlitecloud.Connect(connectionString)
		SQLiteWeb.Logger.Debugf("sqlitecloud.Connect %s err %v", connectionString, err)
		if err == nil && conn != nil {
			break
		}

		select {
		case <-timeout.C:
			err = fmt.Errorf("Cannot connect to sqlitecloud service before timeout %v", pollingTimeout)
			return
		default:
			// non-blocking select
		}

		time.Sleep(pollingSleep)
	}

	return
}

func clusterConfig(projectuuid string, nodeid int, hostname string, port int) (string, int, error) {
	sql := fmt.Sprintf("SELECT node_id as id, hostname || ':' || port as 'public' FROM Node WHERE project_uuid = '%s' AND node_id != %d AND created = 1", projectuuid, nodeid)
	listNodes, err, _, _, _ := dashboardcm.ExecuteSQL("auth", sql)
	if err != nil {
		return "", 0, fmt.Errorf("Cannot get project's nodes %s: %s", projectuuid, err.Error())
	}

	listNodesObj, err := ResultToObj(listNodes)
	if err != nil {
		return "", 0, fmt.Errorf("Cannot get project's nodes %s: %s", projectuuid, err.Error())
	}

	listNodesRowset, ok := listNodesObj.(map[string]interface{})
	if !ok {
		return "", 0, fmt.Errorf("Invalid nodesobj response: %v", listNodesObj)
	}

	listNodesRowsMap, ok := listNodesRowset["rows"].([]map[string]interface{})
	if !ok {
		// add the new node itself, it is a one node cluster
		// otherwise the newly added node must not be included in the cluster config
		listNodesRowsMap = []map[string]interface{}{{"id": nodeid, "public": fmt.Sprintf("%s:%d", hostname, port)}}
	}

	clusterConfig, err := json.Marshal(listNodesRowsMap)
	if err != nil {
		return "", 0, err
	}

	return string(clusterConfig), int(listNodes.GetNumberOfRows()), nil
}

const dropletDomain = "sqlite.cloud"

// hostname labels may contain only the ASCII letters 'a' through 'z' (in a case-insensitive manner),
// the digits '0' through '9', and the hyphen ('-'). The original specification of hostnames in RFC 952,
// mandated that labels could not start with a digit or with a hyphen, and must not end with a hyphen.
// However, a subsequent specification (RFC 1123) permitted hostname labels to start with digits.
// In addition, we want the first char to be a letter.
func validateShortUUID(shortuuid string) bool {
	var firstrune rune
	for _, c := range shortuuid {
		firstrune = c
		break
	}

	switch {
	case len(shortuuid) != 9:
		// sanity check on standard shortuuid len
		return false
	case !unicode.IsLetter(firstrune):
		// first string must be a letter
		return false
	case strings.Count(shortuuid, "-") > 1:
		// no more than on hyphen
		return false
	case strings.Contains(shortuuid, "_"):
		// must not contains _
		return false
	case shortuuid[len(shortuuid)-1:] == "-":
		// must not end with -
		return false
	default:
		return true
	}
}
