package portfolios

// Payload represents the information of a portfolio that can be
// returned to the user
type Payload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ExportPublic returns a PortfolioPayload containing only the fields that are safe to
// be seen by anyone
func (p *Portfolio) ExportPublic() *Payload {
	return &Payload{
		ID:   p.ID,
		Name: p.Name,
	}
}

// ExportPrivate returns a PortfolioPayload containing all the fields
func (p *Portfolio) ExportPrivate() *Payload {
	return p.ExportPublic()
}

// ListPayload represents a list of Portfolio that can be
// safely returned to the clients
type ListPayload struct {
	Results []*Payload `json:"results"`
}

// ExportPublic returns a ProfilesPayload containing only the fields that are safe to
// be seen by anyone
func (p Portfolios) ExportPublic() *ListPayload {
	pld := &ListPayload{}
	pld.Results = make([]*Payload, len(p))
	for i, p := range p {
		pld.Results[i] = p.ExportPublic()
	}
	return pld
}

// ExportPrivate returns a ProfilesPayload containing all the fields
func (p Portfolios) ExportPrivate() *ListPayload {
	pld := &ListPayload{}
	pld.Results = make([]*Payload, len(p))
	for i, p := range p {
		pld.Results[i] = p.ExportPrivate()
	}
	return pld
}
