package repository_external

import "github.com/go-resty/resty/v2"

func (repository *ExternalRepositoryImpl) GetAllSubtaskRepository(cardId string) (*resty.Response, error) {
	return repository.
		RestyConfig.
		R().
		Get(repository.JiraConfig.Url + "/api/3/issue/" + cardId)
}