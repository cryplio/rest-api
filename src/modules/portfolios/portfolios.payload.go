package portfolios

// PortfolioPayload represents the information of a portfolio that can be
// returned to the user
type PortfolioPayload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ExportPublic returns a PortfolioPayload containing only the fields that are safe to
// be seen by anyone
func (p *Portfolio) ExportPublic() *PortfolioPayload {
	return &PortfolioPayload{
		ID:   p.ID,
		Name: p.Name,
	}
}

// ExportPrivate returns a PortfolioPayload containing all the fields
func (p *Portfolio) ExportPrivate() *PortfolioPayload {
	return p.ExportPublic()
}

// PortfoliosPayload represents a list of Portfolio that can be
// safely returned to the clients
type PortfoliosPayload struct {
	Results []*PortfolioPayload `json:"results"`
}

// ExportPublic returns a ProfilesPayload containing only the fields that are safe to
// be seen by anyone
func (p Portfolios) ExportPublic() *PortfoliosPayload {
	pld := &PortfoliosPayload{}
	pld.Results = make([]*PortfolioPayload, len(p))
	for i, p := range p {
		pld.Results[i] = p.ExportPublic()
	}
	return pld
}

// ExportPrivate returns a ProfilesPayload containing all the fields
func (p Portfolios) ExportPrivate() *PortfoliosPayload {
	pld := &PortfoliosPayload{}
	pld.Results = make([]*PortfolioPayload, len(p))
	for i, p := range p {
		pld.Results[i] = p.ExportPrivate()
	}
	return pld
}
