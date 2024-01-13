package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceRepositoryDeployKey(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-project-%v"
		}

		resource "bitbucketserver_repository" "test" {
			project = bitbucketserver_project.test.key
			name = "repo"
		}

		resource "bitbucketserver_repository_deploy_key" "test" {
			project = bitbucketserver_project.test.key
			repository = bitbucketserver_repository.test.slug
			label = "newLabelForTest"
			key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQD1F2KTOi8ptpqkXdcARgYy27uWRCav1isJUB3Hz59MDM6CSAXa+HCgUNDV+6gXQ1eR24nr1Efb7AkEM8LmvMkNlQqpAsPEPxVlndA0KSRLXb3mWruJzkZtrJo1HhVDffuKnevPLFkB75wyBX3jn1FNtc2qfc2i80mu8TfviZnOPYnrRa8B2S9Q8IlUtXAyihQ1G3fn5r5nxrw3QQGrD3cckB9nCZpzAPn6hlDHVk6n5efoWAUxM5AhKln0OLsA1HwjGJN7/dbDPu1nSLuiJSAISSSg4E4SNdfnr3FOhTA79AKNzsor2/EXdbG+f+S4op3s3wtt05zwkHLXRSqKYoT31RFiV9d1XIav+dGTvXCgc6DNG6rbogE6ZbugDZXdcHOAoNs7IUDFbtI/HGKS7CStxAUEchoxM8HDYXmhYt1kUEpmP3g2ckILHGoGPEkOCYGPqx5HbDvAXAJVk3DdSOibCckR2FK2qEoCbMgnPUX84CqNJPHBZ24AaE8htE6TKr0= user@somewhere"
			permission = "REPO_READ"
		}
	`, projectKey, projectKey)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucketserver_repository_deploy_key.test", "id"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_deploy_key.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_repository_deploy_key.test", "repository", "repo"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_deploy_key.test", "key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQD1F2KTOi8ptpqkXdcARgYy27uWRCav1isJUB3Hz59MDM6CSAXa+HCgUNDV+6gXQ1eR24nr1Efb7AkEM8LmvMkNlQqpAsPEPxVlndA0KSRLXb3mWruJzkZtrJo1HhVDffuKnevPLFkB75wyBX3jn1FNtc2qfc2i80mu8TfviZnOPYnrRa8B2S9Q8IlUtXAyihQ1G3fn5r5nxrw3QQGrD3cckB9nCZpzAPn6hlDHVk6n5efoWAUxM5AhKln0OLsA1HwjGJN7/dbDPu1nSLuiJSAISSSg4E4SNdfnr3FOhTA79AKNzsor2/EXdbG+f+S4op3s3wtt05zwkHLXRSqKYoT31RFiV9d1XIav+dGTvXCgc6DNG6rbogE6ZbugDZXdcHOAoNs7IUDFbtI/HGKS7CStxAUEchoxM8HDYXmhYt1kUEpmP3g2ckILHGoGPEkOCYGPqx5HbDvAXAJVk3DdSOibCckR2FK2qEoCbMgnPUX84CqNJPHBZ24AaE8htE6TKr0= user@somewhere"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_deploy_key.test", "label", "newLabelForTest"),
				),
			},
		},
	})
}
