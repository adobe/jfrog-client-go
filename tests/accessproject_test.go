package tests

import (
	"reflect"
	"testing"

	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/stretchr/testify/assert"
)

func TestAccessProject(t *testing.T) {
	initAccessTest(t)
	t.Run("create-update-delete", testAccessProjectCreateUpdateDelete)
}

func testAccessProjectCreateUpdateDelete(t *testing.T) {
	projectParams := getTestProjectParams()
	err := testsAccessProjectService.Create(projectParams)
	defer deleteProjectAndAssert(t, projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	projectParams.ProjectDetails.Description += "123"
	projectParams.ProjectDetails.StorageQuotaBytes += 123
	projectParams.ProjectDetails.SoftLimit = &trueValue
	projectParams.ProjectDetails.AdminPrivileges.ManageMembers = &falseValue
	projectParams.ProjectDetails.AdminPrivileges.ManageResources = &trueValue
	projectParams.ProjectDetails.AdminPrivileges.IndexResources = &falseValue
	err = testsAccessProjectService.Update(projectParams)
	assert.NoError(t, err)
	updatedProject, err := testsAccessProjectService.Get(projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	assert.NotNil(t, updatedProject)
	if !reflect.DeepEqual(projectParams.ProjectDetails, *updatedProject) {
		t.Error("Unexpected project details built. Expected: `", projectParams.ProjectDetails, "` Got `", *updatedProject, "`")
	}
}

func deleteProjectAndAssert(t *testing.T, projectKey string) {
	err := testsAccessProjectService.Delete(projectKey)
	assert.NoError(t, err)
}

func getTestProjectParams() services.ProjectParams {
	adminPrivileges := services.AdminPrivileges{
		ManageMembers:   &trueValue,
		ManageResources: &falseValue,
		IndexResources:  &trueValue,
	}
	runNumberSuffix := getRunId()[len(getRunId())-3 : len(getRunId())]
	projectDetails := services.Project{
		DisplayName:       "testProject" + runNumberSuffix,
		Description:       "My Test Project",
		AdminPrivileges:   &adminPrivileges,
		SoftLimit:         &falseValue,
		StorageQuotaBytes: 1073741825, // needs to be higher than 1073741824
		ProjectKey:        "tstprj" + runNumberSuffix,
	}
	return services.ProjectParams{
		ProjectDetails: projectDetails,
	}
}
