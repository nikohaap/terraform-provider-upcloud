package upcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestUpcloudInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUpcloudInstanceInstanceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("upcloud_instance.my-instance", "zone"),
					resource.TestCheckResourceAttrSet("upcloud_instance.my-instance", "hostname"),
					resource.TestCheckResourceAttr(
						"upcloud_instance.my-instance", "zone", "fi-hel1"),
					resource.TestCheckResourceAttr(
						"upcloud_instance.my-instance", "hostname", "debian.example.com"),
				),
			},
		},
	})
}

func testUpcloudInstanceInstanceConfig() string {
	return fmt.Sprintf(`
		resource "upcloud_instance" "my-instance" {
			zone     = "fi-hel1"
			hostname = "debian.example.com"
		
		
			storage_devices = [{
				size    = 10
				action  = "clone"
				storage = "01000000-0000-4000-8000-000020030100"
			},
				{
					action  = "attach"
					storage = "01000000-0000-4000-8000-000020010301"
					type    = "cdrom"
				},
				{
					action = "create"
					size   = 10
					tier   = "maxiops"
				},
			]
		}
`)
}

func TestAccInstance_changePlan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithSmallInstancePlan,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"upcloud_instance.my-instance", "plan", "1xCPU-2GB"),
				),
			},
			{
				Config: testAccPlanConfigUpdateInstancePlan,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"upcloud_instance.my-instance", "plan", "2xCPU-4GB"),
				),
			},
		},
	})
}

const testAccInstanceConfigWithSmallInstancePlan = `
resource "upcloud_instance" "my-instance" {
			zone     = "fi-hel1"
			hostname = "debian.example.com"
			plan     = "1xCPU-2GB"
		
		
			storage_devices = [{
				size    = 10
				action  = "clone"
				storage = "01000000-0000-4000-8000-000020030100"
			},
				{
					action  = "attach"
					storage = "01000000-0000-4000-8000-000020010301"
					type    = "cdrom"
				},
				{
					action = "create"
					size   = 10
					tier   = "maxiops"
				},
			]
		}
`

const testAccPlanConfigUpdateInstancePlan = `
resource "upcloud_instance" "my-instance" {
			zone     = "fi-hel1"
			hostname = "debian.example.com"
			plan     = "2xCPU-4GB"
		
		
			storage_devices = [{
				size    = 10
				action  = "clone"
				storage = "01000000-0000-4000-8000-000020030100"
			},
				{
					action  = "attach"
					storage = "01000000-0000-4000-8000-000020010301"
					type    = "cdrom"
				},
				{
					action = "create"
					size   = 10
					tier   = "maxiops"
				},
			]
		}
`

func Test_instanceRestartIsRequired(t *testing.T) {
	type args struct {
		storageDevices []interface{}
	}
	cases := []struct {
		name string
		args args
		want bool
	}{
		{"Instance reboot is not required if there's not any valid backup rules", args{
			storageDevices: []interface{}{
				map[string]interface{}{
					"id":          "1",
					"action":      "clone",
					"backup_rule": map[string]interface{}{},
				},
			},
		}, false},
		{"Instance reboot is required if there's at least one valid backup rule", args{
			storageDevices: []interface{}{
				map[string]interface{}{
					"id":     "1",
					"action": "clone",
					"backup_rule": map[string]interface{}{
						"interval":  "test-interval",
						"time":      "test-time",
						"retention": "test-retention",
					},
				},
				map[string]interface{}{
					"id":          "2",
					"action":      "clone",
					"backup_rule": map[string]interface{}{},
				},
				map[string]interface{}{
					"id":     "3",
					"action": "clone",
					"backup_rule": map[string]interface{}{
						"interval":  "test-interval",
						"time":      "test-time",
						"retention": "test-retention",
					},
				},
			},
		}, true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := instanceRestartIsRequired(tt.args.storageDevices); got != tt.want {
				t.Errorf("instanceRestartIsRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}
