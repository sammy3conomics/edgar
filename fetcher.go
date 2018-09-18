package edgar

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type fetcher struct {
}

func (f *fetcher) CompanyFolder(
	ticker string,
	fileTypes ...FilingType) (CompanyFolder, error) {

	comp := newCompany(ticker)

	for _, t := range fileTypes {
		comp.FilingLinks[t] = getFilingLinks(ticker, t)
	}
	return comp, nil
}

// Read from the reader and unmarshal from JSON to company folder
func (f *fetcher) CreateFolder(
	r io.Reader,
	fileTypes ...FilingType) (CompanyFolder, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	c := new(company)
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	// Get all the latest links for all the filing types
	for _, key := range fileTypes {
		c.FilingLinks[key] = getFilingLinks(c.Ticker(), key)
	}
	return c, nil
}

func NewFilingFetcher() FilingFetcher {
	return &fetcher{}
}