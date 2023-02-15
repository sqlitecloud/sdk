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
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/digitalocean/godo"
)

const (
	imageSlug      = "ubuntu-22-04-x64"
	doProviderName = "DigitalOcean"
)

var (
	digitalocean   CloudProvider
	regionLocation map[string]Coordinates
)

type CloudProvider interface {
	CreateNode(nodeCreateReq *CloudNodeCreateRequest) (*CloudNode, error)
	DestroyNode(cloudNode *CloudNode) error
}

// DropletCreateRequest represents a request to create a Droplet.
type CloudNodeCreateRequest struct {
	JobUUID       string
	Name          string
	Region        string
	Size          string
	Type          string
	ProjectUUID   string
	NodeID        int
	Hostname      string
	Domain        string
	Port          int
	ClusterPort   int
	ClusterConfig string
	Token         string
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type CloudNode struct {
	JobUUID        string
	Name           string
	Region         string
	Size           string
	Type           string
	ProjectUUID    string
	NodeID         int
	Hostname       string
	Domain         string
	AddrV4         string
	AddrV6         string
	Port           int
	Location       Coordinates
	Provider       string
	DropletID      int
	DomainRecordID int
}

func (node *CloudNode) FullyQualifiedDomainName() string {
	return fmt.Sprintf("%s.%s", node.Hostname, node.Domain)
}

type CloudProviderDigitalOcean struct {
	doclient *godo.Client
}

func init() {
	regionLocation = map[string]Coordinates{
		"nyc1": {40.7966743, -74.0334953},
		"nyc2": {40.7414619, -74.0052546},
		"nyc3": {40.8299598, -74.128822},
		"sfo1": {37.723698, -122.4002447},
		"sfo2": {37.7887409, -122.3927261},
		"sfo3": {37.3758049, -121.9745899},
		"ams2": {52.2933512, 4.9428649},
		"ams3": {52.3030946, 4.920444}, // DigitalOcean didn't publicly list where AMS3 is, we expect AMS3 will be another Equinix facility, there are 7 Equinix Data centres locations in Amsterdam.
		"sgp1": {1.3214987, 103.6931552},
		"lon1": {51.5224355, -0.6312062},
		"fra1": {50.1196098, 8.7360115},
		"tor1": {43.6508826, -79.3639686},
		"blr1": {12.8396177, 77.6593288},
		"syd1": {-33.847927, 150.6517896}, // geo location of sydney, not found any info for the datacenter
	}
}

func initCloudProviderDigitalOcean() {
	if cfg.Section("digitalocean").HasKey("token") {
		digitalocean = NewCloudProviderDigitalOcean(cfg.Section("digitalocean").Key("token").String())
	}
}

// NewService - our constructor function
func NewCloudProviderDigitalOcean(token string) *CloudProviderDigitalOcean {
	cpdo := &CloudProviderDigitalOcean{
		doclient: godo.NewFromToken(token),
	}

	return cpdo
}

var cloudConfigTemplate = template.Must(template.New("").Parse(`#cloud-config
users:
- name: demo
  groups: sudo
  shell: /bin/bash
  sudo: ALL=(ALL) NOPASSWD:ALL
  ssh_authorized_keys:
    - ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGkf64lQR2QEDT5hef+S0opIXwJ/bkihpHkzJW/IFu96 andrea@sqlitecloud.io
write_files:
  - path: /etc/sqlitecloud/node.ini
    content: |
      uuid = {{.ProjectUUID}}
      cluster_node_id = {{.NodeID}}
      base_path = /var/lib/sqlitecloud
      listening_port = {{.Port}} 
      cluster_port = {{.ClusterPort}}
      cluster_config = {{.ClusterConfig}}
      backup_node_id = 1
      backup_config = {"checkpoint-interval": "1m", "replicas": [{"url": "s3://sqlc-dev.s3.us-west-004.backblazeb2.com", "access-key-id": "004460c459c4bae0000000001", "secret-access-key": "K004ItQqhbR3Npy/S/qObqiYh/bLZas"}]}
  - path: /etc/systemd/system/sqlitecloud.service
    content: |
      [Unit]
      Description=SQLiteCloud
      After=network.target
      StartLimitIntervalSec=0
      #OnFailure=unit-status-mail@%n.service
      [Service]
      Type=simple
      #Restart=always
      #RestartSec=1
      ExecStart=/usr/local/sbin/sqlitecloud --config /etc/sqlitecloud/node.ini
      #KillMode=process
      [Install]
      WantedBy=multi-user.target
runcmd:
  - cd /root
  - mkdir -p /var/lib/sqlitecloud/{{.Port}}
  - wget --no-check-certificate 'https://docs.google.com/uc?export=download&id=1a8gqZA_R-m0R4BZ_F_7JT0_kuwVgsgtc' -O sqlitecloud-v0.9.8-linux-amd64.tar.gz
  - tar xvzf sqlitecloud-v0.9.8-linux-amd64.tar.gz
  - mv -t /usr/local/sbin/ sqlitecloud libraft.so
  - wget --no-check-certificate 'https://docs.google.com/uc?export=download&id=1iNd2GMwfEkCvqfp_dcqn8sURt7JfyfiD' -O litestream-v0.3.9-enc.tar.gz
  - tar xvzf litestream-v0.3.9-enc.tar.gz
  - mv litestream /usr/bin/
  - snap install core; snap refresh core; snap install --classic certbot; ln -s /snap/bin/certbot /usr/bin/certbot
  - certbot certonly --standalone -d {{.Hostname}}.{{.Domain}} --non-interactive --agree-tos -m certbot@sqlitecloud.io > certbot.log 2> certbot.err.log
  - ln -s /etc/letsencrypt/live/{{.Hostname}}.{{.Domain}}/privkey.pem /var/lib/sqlitecloud/{{.Port}}/certificate_key.pem
  - ln -s /etc/letsencrypt/live/{{.Hostname}}.{{.Domain}}/fullchain.pem /var/lib/sqlitecloud/{{.Port}}/certificate.pem
  - systemctl start sqlitecloud`))

//   - DIGITALOCEAN_TOKEN={{.Token}}
//   - curl -X POST -H "Content-Type:\ application/json" -H "Authorization:\ Bearer $DIGITALOCEAN_TOKEN" -d "{\"type\":\"A\",\"name\":\"{{.Hostname}}\",\"data\":\"$(curl -s ifconfig.me)\",\"priority\":null,\"port\":null,\"ttl\":1800,\"weight\":null,\"flags\":null,\"tag\":null}" "https://api.digitalocean.com/v2/domains/{{.Domain}}/records" -o "sqlite_cloud_dns_records.log"

//   - sed -i "s/` + dropletAddrPortStringPlaceholder + `/$(curl -s ifconfig.me):{{.Port}}/g" /etc/sqlitecloud/node.ini

func setDefaults(nodeCreateReq *CloudNodeCreateRequest) {
	if nodeCreateReq.Port == 0 {
		nodeCreateReq.Port = 8860
	}

	if nodeCreateReq.ClusterPort == 0 {
		nodeCreateReq.ClusterPort = nodeCreateReq.Port + 1000
	}

	if nodeCreateReq.Token == "" && cfg.Section("digitalocean").HasKey("token") {
		nodeCreateReq.Token = cfg.Section("digitalocean").Key("token").String()
	}
}

func (this *CloudProviderDigitalOcean) CreateNode(nodeCreateReq *CloudNodeCreateRequest) (*CloudNode, error) {
	setDefaults(nodeCreateReq)

	cloudConfigBuf := new(bytes.Buffer)
	cloudConfigTemplate.Execute(cloudConfigBuf, nodeCreateReq)

	// SQLiteWeb.Logger.Infof("Cloud Config:\n\n\n%s\n\n\n", cloudConfigBuf.String())

	// 1. Create Droplet Request
	createDropletRequest := &godo.DropletCreateRequest{
		Name:   nodeCreateReq.Name,
		Region: nodeCreateReq.Region, // "nyc3",
		Size:   nodeCreateReq.Size,   // "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			Slug: imageSlug,
		},
		Tags:     []string{"sqlitecloud", "test"},
		SSHKeys:  []godo.DropletCreateSSHKey{{Fingerprint: "f0:42:56:b6:23:2a:72:0a:47:94:f4:08:10:32:fb:8d"}},
		UserData: cloudConfigBuf.String(),
	}

	ctxDroplet, ctxDropletCancel := context.WithTimeout(context.Background(), cloudRequestTimeout)
	defer ctxDropletCancel()

	newDroplet, resp, err := this.doclient.Droplets.Create(ctxDroplet, createDropletRequest)
	if err != nil {
		err = fmt.Errorf("Error: cannot create a digitalocean droplet: %s", err.Error())
		return nil, err
	}
	if newDroplet == nil {
		err = errors.New("Error: cannot create a digitalocean droplet")
		return nil, err
	}

	cloudNode := &CloudNode{
		JobUUID:     nodeCreateReq.JobUUID,
		Name:        newDroplet.Name,
		Region:      newDroplet.Region.Slug,
		Size:        newDroplet.Size.Slug,
		Type:        nodeCreateReq.Type,
		ProjectUUID: nodeCreateReq.ProjectUUID,
		NodeID:      nodeCreateReq.NodeID,
		Hostname:    nodeCreateReq.Hostname,
		Domain:      nodeCreateReq.Domain,
		Port:        nodeCreateReq.Port,
		Location:    regionLocation[newDroplet.Region.Slug],
		Provider:    doProviderName,
		DropletID:   newDroplet.ID,
	}

	if len(resp.Links.Actions) != 1 {
		err = fmt.Errorf("Error: invalid response links: %d", len(resp.Links.Actions))
		return cloudNode, err
	}

	SQLiteWeb.Logger.Infof("Droplet created %s %s %d: %s", nodeCreateReq.Name, nodeCreateReq.JobUUID, newDroplet.ID, resp.Links.Actions[0].HREF)

	// 2. Poll the Actions endpoint to get the Droplet addresses
	ctx := context.TODO()
	timeout := time.NewTimer(pollingTimeout)
	for {
		action, _, err := this.doclient.Actions.Get(ctx, resp.Links.Actions[0].ID)
		if err != nil {
			err = fmt.Errorf("Error: cannot get actions for created droplet: %s", err.Error())
			return cloudNode, err
		}

		completed := false
		switch action.Status {
		case "errored":
			err = errors.New("Error: create droplet action status is errored")
			return cloudNode, err
		case "completed":
			completed = true
		}

		if completed {
			break
		}

		select {
		case <-timeout.C:
			err = fmt.Errorf("Error: cannot get new droplet's info before timeout %v", pollingTimeout)
			return cloudNode, err
		default:
			// non-blocking select
		}

		time.Sleep(pollingSleep)
	}

	newDroplet, resp, err = this.doclient.Droplets.Get(ctx, newDroplet.ID)
	if err != nil {
		err = fmt.Errorf("Error: cannot get new droplet's info: %s", err.Error())
		return cloudNode, err
	}

	addrV4 := ""
	for _, addr := range newDroplet.Networks.V4 {
		if addr.Type == "public" {
			addrV4 = addr.IPAddress
			break
		}
	}

	addrV6 := ""
	for _, addr := range newDroplet.Networks.V6 {
		if addr.Type == "public" {
			addrV6 = addr.IPAddress
			break
		}
	}

	cloudNode.AddrV4 = addrV4
	cloudNode.AddrV6 = addrV6

	// 3. Add the DNS record
	ctxDnsRecord, ctxDnsRecordCancel := context.WithTimeout(context.Background(), cloudRequestTimeout)
	defer ctxDnsRecordCancel()

	createDomainRecordRequest := &godo.DomainRecordEditRequest{
		Type: "A",
		Name: nodeCreateReq.Hostname,
		Data: addrV4,
	}

	if domainRecord, _, err := this.doclient.Domains.CreateRecord(ctxDnsRecord, nodeCreateReq.Domain, createDomainRecordRequest); err != nil {
		err = fmt.Errorf("Error: cannot create new domain record: %s.%s %s", nodeCreateReq.Hostname, nodeCreateReq.Domain, err.Error())
		return cloudNode, err
	} else {
		cloudNode.DomainRecordID = domainRecord.ID
	}

	return cloudNode, nil
}

func (this *CloudProviderDigitalOcean) DestroyNode(cloudNode *CloudNode) error {
	if cloudNode == nil || cloudNode.DropletID == 0 {
		return nil
	}

	if cloudNode.Provider != doProviderName {
		return fmt.Errorf("Error: expecting provider %s, got %s (id: %d)", doProviderName, cloudNode.Provider, cloudNode.DropletID)
	}

	ctx := context.TODO()
	if _, err := this.doclient.Droplets.Delete(ctx, cloudNode.DropletID); err != nil {
		return fmt.Errorf("Error: cannot delete droplet %d: %s", cloudNode.DropletID, err.Error())
	}

	if cloudNode.DomainRecordID != 0 {
		ctx := context.TODO()
		if _, err := this.doclient.Domains.DeleteRecord(ctx, cloudNode.Domain, cloudNode.DomainRecordID); err != nil {
			return fmt.Errorf("Error: cannot delete domain record %d %s.%s: %s", cloudNode.DomainRecordID, cloudNode.Hostname, cloudNode.Domain, err.Error())
		}
	}

	return nil
}
