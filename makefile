.PHONY = generate-mocks


# ############################
# TEST
# ############################

generate-mocks:
	#go get github.com/golang/mock/gomock
	#go install github.com/golang/mock/mockgen
	mockgen -destination=service/test_mocks/repo_mock.go -package=test_mocks github.com/lutomas/PR00B121-TDD/service Repo
