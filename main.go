package main

import (
  "fmt"
  "os"
  "github.com/gophercloud/gophercloud"
)

var (
  provider *gophercloud.ProviderClient
  region string = ""
)

func authOptions() gophercloud.AuthOptions {
  return gophercloud.AuthOptions{
    IdentityEndpoint: os.Getenv("OS_AUTH_URL"), // cfg.Global.AuthUrl,
    Username:         os.Getenv("OS_USERNAME"), // cfg.Global.Username,
    UserID:           os.Getenv("OS_USER_ID"), // known to be unset, will be "" // cfg.Global.UserId,
    Password:         os.Getenv("OS_PASSWORD"), // cfg.Global.Password,
    TenantID:         os.Getenv("OS_TENANT_ID"), // cfg.Global.TenantId,
    TenantName:       os.Getenv("OS_TENANT_NAME"), // known to be unset, will be "" // cfg.Global.TenantName,
    DomainID:         os.Getenv("OS_USER_DOMAIN_ID"), // known to be unset, will be "" // cfg.Global.DomainId,
    DomainName:       os.Getenv("OS_USER_DOMAIN_NAME"), // cfg.Global.DomainName,

    // Persistent service, so we need to be able to renew tokens.
    AllowReauth: true,
  }
}

func GetZone() (cloudprovider.Zone, error) {
  md, err := getMetadata()
  if err != nil {
    return cloudprovider.Zone{}, err
  }

  zone := cloudprovider.Zone{
    FailureDomain: "usw1",
    Region:        os.region,
  }
  glog.V(1).Infof("Current zone is %v", zone)

  return zone, nil
}

func main() {
  fmt.Println("starting openstack-token test")

  fmt.Println("creating NewClient")
  provider, err := openstack.NewClient(cfg.Global.AuthUrl)
  if err != nil {
    fmt.Println("failed to create provider: %v", err)
    os.Exit(1)
  }

  fmt.Println("authenticating Provider")
  err = openstack.Authenticate(provider, authOptions())
  if err != nil {
    fmt.Println("failed to authenticate: %v", err)
    os.Exit(1)
  }
}

