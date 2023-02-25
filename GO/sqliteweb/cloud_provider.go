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
	pollingSleep        = 5 * time.Second
	pollingTimeout      = 10 * time.Minute
	cloudRequestTimeout = 5 * time.Minute
)

type CloudRegion string
type CloudSize string

const (
	// https://docs.digitalocean.com/reference/api/api-reference/#tag/Regions
	CloudRegionNewYork1      CloudRegion = "New York 1"
	CloudRegionNewYork3      CloudRegion = "New York 3"
	CloudRegionSanFrancisco3 CloudRegion = "San Francisco 3"
	CloudRegionAmsterdam3    CloudRegion = "Amsterdam 3"
	CloudRegionLondon1       CloudRegion = "London 1"
	CloudRegionFrankfurt1    CloudRegion = "Frankfurt 1"
	CloudRegionSingapore1    CloudRegion = "Singapore 1"
	CloudRegionToronto1      CloudRegion = "Toronto 1"
	CloudRegionBangalore1    CloudRegion = "Bangalore 1"
	CloudRegionSydney1       CloudRegion = "Sydney 1"
	// not available regions:
	// NewYork2 CloudRegion = "New York 2"
	// SanFrancisco1 CloudRegion = "San Francisco 1"
	// SanFrancisco2 CloudRegion = "San Francisco 2"
	// Amsterdam2 CloudRegion = "Amsterdam 2"
)

const (
	CloudSize_1_1_25 CloudSize = "1VCPU/1GB/25GB"
	CloudSize_1_2_50 CloudSize = "1VCPU/2GB/50GB"
	CloudSize_2_2_60 CloudSize = "2VCPU/2GB/60GB"
	// CloudSize_2_16_300 CloudSize = "2VCPU/16GB/300GB"
)

var CloudRegions []CloudRegion
var CloudSizes []CloudSize

func init() {
	CloudRegions = []CloudRegion{CloudRegionNewYork1, CloudRegionNewYork3, CloudRegionSanFrancisco3, CloudRegionAmsterdam3, CloudRegionLondon1, CloudRegionFrankfurt1, CloudRegionSingapore1, CloudRegionToronto1, CloudRegionBangalore1, CloudRegionSydney1}
	CloudSizes = []CloudSize{CloudSize_1_1_25, CloudSize_1_2_50, CloudSize_2_2_60}
}

func createNode(userid int, name string, region CloudRegion, size CloudSize, nodetype string, projectuuid string, nodeid int) (string, error) {
	nodeshortid := ""
	uniqueNodeID := int64(0)
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
		sql := fmt.Sprintf("INSERT INTO Node (project_uuid, node_id, shortuuid, name) VALUES ('%s', %d, '%s', '%s') RETURNING id;", projectuuid, nodeid, sid, name)
		res, err, errcode, _, _ := dashboardcm.ExecuteSQL("auth", sql)
		if err == nil && res.GetNumberOfColumns() == 1 && res.GetNumberOfRows() == 1 {
			nodeshortid = sid
			uniqueNodeID = res.GetInt64Value_(0, 0)
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
	sql := fmt.Sprintf("INSERT INTO Jobs (uuid, name, status, steps, progress, user_id, node_id) VALUES ('%s', 'Create Node %s', 'Creating droplet', 2, 0, %d, %d)", jobuuid, name, userid, uniqueNodeID) // ; SELECT uuid FROM Jobs WHERE rowid = last_insert_rowid();
	_, err, _, _, _ = dashboardcm.ExecuteSQL("auth", sql)
	if err != nil {
		return "", fmt.Errorf("Cannot create the job %s: %s", jobuuid, err.Error())
	}

	cloudProvider := digitalocean
	sqlitecloudPort := 8860

	go func() {
		clusterConfig, isFirstNode, err := clusterConfig(projectuuid, nodeid, fmt.Sprintf("%s.%s", nodeshortid, dropletDomain), sqlitecloudPort)
		if err != nil {
			SQLiteWeb.Logger.Error(err.Error())
			sql := fmt.Sprintf("UPDATE Jobs SET status = '%s', error = 1, stamp = DATETIME('now') WHERE uuid = '%s'", err.Error(), jobuuid)
			authExecSQL(sql)
			return
		}

		nodeAddedToClusterConf := false

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
			NewCluster:    isFirstNode,
			Tags:          []string{"sqlitecloud", projectuuid},
		}

		cloudNode, err := cloudProvider.CreateNode(nodeCreateReq)
		if err != nil {
			SQLiteWeb.Logger.Errorf("digitalocean CreateNode: %s", err.Error())
			sql := fmt.Sprintf("UPDATE Jobs SET status = '%s', error = 1, stamp = DATETIME('now') WHERE uuid = '%s'", err.Error(), jobuuid)
			authExecSQL(sql)
			deleteCloudNode_(cloudProvider, cloudNode, uniqueNodeID, nodeAddedToClusterConf)
			return
		}
		SQLiteWeb.Logger.Debugf("Droplet completed %s %s %d", cloudNode.Name, cloudNode.JobUUID, cloudNode.DropletID)

		if !isFirstNode {
			// execute the ADD NODE command to the cluster
			addCommand := fmt.Sprintf("ADD NODE %d ADDRESS %s:%d", cloudNode.NodeID, cloudNode.FullyQualifiedDomainName(), cloudNode.Port)
			_, err, _, _, _ = dashboardcm.ExecuteSQL(projectuuid, addCommand)
			if err != nil {
				err = fmt.Errorf("Cannot add the new node %d to the cluster %s: %s", cloudNode.NodeID, cloudNode.JobUUID, err.Error())
				SQLiteWeb.Logger.Error(err.Error())
				sql := fmt.Sprintf("UPDATE Jobs SET status = '%s', error = 1, stamp = DATETIME('now') WHERE uuid = '%s'", err.Error(), jobuuid)
				authExecSQL(sql)
				deleteCloudNode_(cloudProvider, cloudNode, uniqueNodeID, nodeAddedToClusterConf)
				return
			}
			SQLiteWeb.Logger.Debugf("Node added to the cluster %d %s", cloudNode.NodeID, cloudNode.JobUUID)
			nodeAddedToClusterConf = true
		}

		sql = fmt.Sprintf("UPDATE Jobs SET progress = progress + 1, status = 'System setup', stamp = DATETIME('now') WHERE uuid = '%s'", nodeCreateReq.JobUUID)
		authExecSQL(sql)

		sql = fmt.Sprintf("UPDATE Node SET hostname = '%s', type = '%s', provider = '%s', image = '%s', region = '%s', addr4 = '%s', addr6 = '%s', port = %d, latitude = %f, longitude = %f, droplet_id = %d, domain_record_id = %d WHERE id = %d", cloudNode.FullyQualifiedDomainName(), cloudNode.Type, "DigitalOcean", cloudNode.Size, cloudNode.Region, cloudNode.AddrV4, cloudNode.AddrV6, cloudNode.Port, cloudNode.Location.Latitude, cloudNode.Location.Longitude, cloudNode.DropletID, cloudNode.DomainRecordID, uniqueNodeID)
		authExecSQL(sql)

		adminUser, adminPasswd, tmpAdminUser, tmpAdminPasswd := getProjectAdminCredentials(projectuuid, isFirstNode)

		conn, err := waitForConnection(cloudNode, tmpAdminUser, tmpAdminPasswd)
		if err != nil {
			SQLiteWeb.Logger.Error(err.Error())
			sql := fmt.Sprintf("UPDATE Jobs SET status = '%s', error = 1, stamp = DATETIME('now') WHERE uuid = '%s'", err.Error(), jobuuid)
			authExecSQL(sql)
			deleteCloudNode_(cloudProvider, cloudNode, uniqueNodeID, nodeAddedToClusterConf)
			return
		}

		if err := setServerAdminCredentials(conn, adminUser, adminPasswd, isFirstNode); err != nil {
			SQLiteWeb.Logger.Error(err.Error())
			sql := fmt.Sprintf("UPDATE Jobs SET status = '%s', error = 1, stamp = DATETIME('now') WHERE uuid = '%s'", err.Error(), jobuuid)
			authExecSQL(sql)
			deleteCloudNode_(cloudProvider, cloudNode, uniqueNodeID, nodeAddedToClusterConf)
			return
		}

		sql = fmt.Sprintf("UPDATE Jobs SET progress = progress + 1, status = 'Completed', stamp = DATETIME('now') WHERE uuid = '%s'", nodeCreateReq.JobUUID)
		authExecSQL(sql)

		sql = fmt.Sprintf("UPDATE Node SET active = 1 WHERE id = %d", uniqueNodeID)
		authExecSQL(sql)
	}()

	return jobuuid, nil
}

func deleteCloudNode_(cloudProvider CloudProvider, cloudNode *CloudNode, uniqueNodeID int64, nodeAddedToClusterConf bool) {
	if err := deleteCloudNode(cloudProvider, cloudNode, uniqueNodeID, nodeAddedToClusterConf); err != nil {
		SQLiteWeb.Logger.Error(err.Error())
		return
	}
}

func deleteCloudNode(cloudProvider CloudProvider, cloudNode *CloudNode, uniqueNodeID int64, nodeAddedToClusterConf bool) error {
	if nodeAddedToClusterConf {
		removeCommand := fmt.Sprintf("REMOVE NODE %d", cloudNode.NodeID)
		if _, err, _, _, _ := dashboardcm.ExecuteSQL(cloudNode.ProjectUUID, removeCommand); err != nil {
			return err
		}
	}

	if err := cloudProvider.DeleteNode(cloudNode); err != nil {
		return err
	}

	// sql := fmt.Sprintf("DELETE FROM Node WHERE id = %d", uniqueNodeID)
	// authExecSQL(sql)

	return nil
}

func deleteNode(userid int, uniqueNodeID int64, clusterNodeID int, projectUUID string) error {
	cloudProvider := digitalocean

	cloudNode := &CloudNode{
		NodeID:      clusterNodeID,
		ProjectUUID: projectUUID,
		Domain:      dropletDomain,
		Provider:    cloudProvider.ProviderName(),
	}

	sql := fmt.Sprintf("SELECT shortuuid, droplet_id, domain_record_id, name FROM Node WHERE id = %d AND active = 1", uniqueNodeID)
	if res, err, _, _, _ := dashboardcm.ExecuteSQL("auth", sql); err != nil {
		return fmt.Errorf("Error n delete node %d: %s", uniqueNodeID, err.Error())
	} else if !res.IsRowSet() {
		return fmt.Errorf("Error in delete node %d: Invalid response ", uniqueNodeID)
	} else if res.GetNumberOfRows() == 0 {
		return fmt.Errorf("Error in delete node %d: 0 rows", uniqueNodeID)
	} else {
		cloudNode.Hostname = res.GetStringValue_(0, 0)
		cloudNode.DropletID = int(res.GetInt32Value_(0, 1))
		cloudNode.DomainRecordID = int(res.GetInt32Value_(0, 2))
		cloudNode.Name = res.GetStringValue_(0, 3)
	}

	jobuuid := uuid.New().String()
	sql = fmt.Sprintf("INSERT INTO Jobs (uuid, name, status, steps, progress, user_id, node_id) VALUES ('%s', 'Delete Node %s', 'Deleting droplet', 1, 0, %d, %d)", jobuuid, cloudNode.Name, userid, uniqueNodeID) // ; SELECT uuid FROM Jobs WHERE rowid = last_insert_rowid();
	if _, err, _, _, _ := dashboardcm.ExecuteSQL("auth", sql); err != nil {
		return fmt.Errorf("Cannot create the job %s: %s", jobuuid, err.Error())
	}

	if err := deleteCloudNode(cloudProvider, cloudNode, uniqueNodeID, true); err != nil {
		sql = fmt.Sprintf("UPDATE Jobs SET error = 1, status = '%s' WHERE uuid = '%s'", err.Error(), jobuuid)
		authExecSQL(sql)
		return err
	}

	sql = fmt.Sprintf("UPDATE Node SET active = 0 WHERE id = %d", uniqueNodeID)
	authExecSQL(sql)

	sql = fmt.Sprintf("UPDATE Jobs SET progress = steps, status = 'Completed' WHERE uuid = '%s'", jobuuid)
	authExecSQL(sql)

	return nil
}

func getProjectAdminCredentials(projectuuid string, isFirstNode bool) (adminUser string, adminPassword string, tmpAdminUser string, tmpAdminPassword string) {
	adminUser, adminPassword, tmpAdminUser, tmpAdminPassword = "admin", "admin", "admin", "admin"

	sql := fmt.Sprintf("SELECT admin_username, admin_password FROM Project WHERE uuid = '%s'", projectuuid)
	if res, err, _, _, _ := dashboardcm.ExecuteSQL("auth", sql); err != nil {
		SQLiteWeb.Logger.Errorf("Cannot get admin credentials %s", err.Error())
	} else {
		adminUser = res.GetStringValue_(0, 0)
		adminPassword = res.GetStringValue_(0, 1)

		if !isFirstNode {
			tmpAdminUser = adminUser
			tmpAdminPassword = adminPassword
		}
	}

	return
}

func setServerAdminCredentials(conn *sqlitecloud.SQCloud, adminUser string, adminPasswd string, isFirstNode bool) error {
	if isFirstNode {
		if adminUser != "" && adminUser != "admin" {
			if err := conn.ExecuteArray("RENAME USER admin TO ?", []interface{}{adminUser}); err != nil {
				return err
			}

		}
		if adminPasswd != "" && adminPasswd != "admin" {
			if err := conn.ExecuteArray("SET MY PASSWORD ?", []interface{}{adminPasswd}); err != nil {
				return err
			}
		}
	}

	return nil
}

func authExecSQL(sql string) {
	_, err, _, _, _ := dashboardcm.ExecuteSQL("auth", sql)
	if err != nil {
		SQLiteWeb.Logger.Errorf("Error '%s' for query '%s'", err.Error(), sql)
	}
}

func waitForConnection(cloudNode *CloudNode, user string, passwd string) (conn *sqlitecloud.SQCloud, err error) {
	timeout := time.NewTimer(pollingTimeout)
	for {
		connectionString := fmt.Sprintf("sqlitecloud://%s:%s@%s:%d?timeout=5", user, passwd, cloudNode.FullyQualifiedDomainName(), cloudNode.Port)
		conn, err = sqlitecloud.Connect(connectionString)
		SQLiteWeb.Logger.Debugf("sqlitecloud.Connect %s err %v", cloudNode.FullyQualifiedDomainName(), err)
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

func clusterConfig(projectuuid string, nodeid int, hostname string, port int) (clusterConfig string, isFirstNode bool, err error) {
	sql := fmt.Sprintf("SELECT node_id as id, hostname || ':' || port as 'public' FROM Node WHERE project_uuid = '%s' AND node_id != %d AND active = 1", projectuuid, nodeid)
	listNodes, err, _, _, _ := dashboardcm.ExecuteSQL("auth", sql)
	if err != nil {
		err = fmt.Errorf("Cannot get project's nodes %s: %s", projectuuid, err.Error())
		return
	}

	listNodesObj, err := ResultToObj(listNodes)
	if err != nil {
		err = fmt.Errorf("Cannot get project's nodes %s: %s", projectuuid, err.Error())
		return
	}

	listNodesRowset, ok := listNodesObj.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("Invalid nodesobj response: %v", listNodesObj)
		return
	}

	listNodesRowsMap, ok := listNodesRowset["rows"].([]map[string]interface{})
	if !ok {
		// add the new node itself, it is a one node cluster
		// otherwise the newly added node must not be included in the cluster config
		listNodesRowsMap = []map[string]interface{}{{"id": nodeid, "public": fmt.Sprintf("%s:%d", hostname, port)}}
	}

	clusterConfigB, err := json.Marshal(listNodesRowsMap)
	if err != nil {
		return
	}

	clusterConfig = string(clusterConfigB)
	isFirstNode = listNodes.GetNumberOfRows() == 0
	return
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
