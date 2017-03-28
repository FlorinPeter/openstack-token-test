package main

import (
  "fmt"
  "os"
  "time"

  "github.com/gophercloud/gophercloud"
  "github.com/gophercloud/gophercloud/openstack"
  "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
  "github.com/gophercloud/gophercloud/pagination"
)

var (
  provider *gophercloud.ProviderClient
  region string = "" // appears to be set in the openstack cloud_config, not auto-discovered
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

func main() {
  var err error
  fmt.Println("starting openstack-token test")

  fmt.Println("creating NewClient")
  provider, err = openstack.NewClient(os.Getenv("OS_AUTH_URL"))
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

  for {
    fmt.Printf("Starting loop at %s", time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
    err = doLoop()
    if err != nil {
      fmt.Println("failed on loop: %v", err)
      os.Exit(1)
    }

    fmt.Println("now sleeping for an hour...")
    time.Sleep(1 * time.Hour)

  }

}

func doLoop() (err error) {
  fmt.Println("doLoop called")

  fmt.Println("creating compute handle")
  compute, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
    Region: region,
  })
  if err != nil {
    return
  }

  fmt.Println("starting pagination")
  name_filter := ""
  opts := servers.ListOpts{
    Name:   name_filter,
    Status: "ACTIVE",
  }
  pager := servers.List(compute, opts)

  names := make([]string, 0)
  err = pager.EachPage(func(page pagination.Page) (bool, error) {
    sList, err := servers.ExtractServers(page)
    if err != nil {
      return false, err
    }
    for i := range sList {
      names = append(names, sList[i].Name)
    }
    return true, nil
  })
  if err != nil {
    return err
  }
  fmt.Println("pagination complete")

  fmt.Printf("Found %v instances matching %v: %v\n",
    len(names), name_filter, names)

  return
}