package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

func TestAccScalewayRdbInstanceBeta(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewayRdbInstanceBetaDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_rdb_instance_beta main {
						name = "test-rdb"
						node_type = "db-dev-s"
						engine = "PostgreSQL-11"
						is_ha_cluster = true
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance_beta", "minimal" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayRdbBetaExists("scaleway_rdb_instance_beta.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "name", "test-rdb"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "node_type", "db-dev-s"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "engine", "PostgreSQL-11"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "is_ha_cluster", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "disable_backup", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "user_name", "my_initial_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "password", "thiZ_is_v&ry_s3cret"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "tags.1", "scaleway_rdb_instance_beta"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "tags.2", "minimal"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance_beta.main", "endpoint_ip"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance_beta.main", "endpoint_port"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance_beta.main", "certificate"),
				),
			},
			{
				Config: `
					resource scaleway_rdb_instance_beta main {
						name = "test-rdb"
						node_type = "db-dev-m"
						engine = "PostgreSQL-11"
						is_ha_cluster = true
						disable_backup = false
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance_beta", "minimal" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayRdbBetaExists("scaleway_rdb_instance_beta.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "name", "test-rdb"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "node_type", "db-dev-m"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "engine", "PostgreSQL-11"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "is_ha_cluster", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "disable_backup", "false"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "user_name", "my_initial_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "password", "thiZ_is_v&ry_s3cret"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "tags.1", "scaleway_rdb_instance_beta"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance_beta.main", "tags.2", "minimal"),
				),
			},
		},
	})
}

func testAccCheckScalewayRdbBetaExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		rdbAPI, region, ID, err := rdbAPIWithRegionAndID(testAccProvider.Meta(), rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = rdbAPI.GetInstance(&rdb.GetInstanceRequest{
			InstanceID: ID,
			Region:     region,
		})

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayRdbInstanceBetaDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "scaleway_rdb_instance_beta" {
			continue
		}

		rdbAPI, region, ID, err := rdbAPIWithRegionAndID(testAccProvider.Meta(), rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = rdbAPI.GetInstance(&rdb.GetInstanceRequest{
			InstanceID: ID,
			Region:     region,
		})

		// If no error resource still exist
		if err == nil {
			return fmt.Errorf("Instance (%s) still exists", rs.Primary.ID)
		}

		// Unexpected api error we return it
		if !is404Error(err) {
			return err
		}
	}

	return nil
}
