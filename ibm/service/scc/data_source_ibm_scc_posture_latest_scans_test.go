// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc_test

import (
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMSccPostureListLatestScansDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMSccPostureListLatestScansDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_scc_posture_latest_scans.list_latest_scans", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_posture_latest_scans.list_latest_scans", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_posture_latest_scans.list_latest_scans", "last.#"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_posture_latest_scans.list_latest_scans", "latest_scans.#"),
				),
			},
		},
	})
}

func testAccCheckIBMSccPostureListLatestScansDataSourceConfigBasic() string {
	return `
		data "ibm_scc_posture_latest_scans" "list_latest_scans" {
		}
	`
}
