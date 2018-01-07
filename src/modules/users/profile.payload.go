package users

// ProfilePayload represents the public information of a user
type ProfilePayload struct {
	*Payload // User payload
}

// ExportPublic returns a ProfilePayload containing only the fields that are safe to
// be seen by anyone
func (p *Profile) ExportPublic() *ProfilePayload {
	// It's OK to export a nil experience
	if p == nil {
		return nil
	}

	return &ProfilePayload{
		Payload: NewPayload(p.User),
	}
}

// ExportPrivate returns a ProfilePayload containing all the fields
func (p *Profile) ExportPrivate() *ProfilePayload {
	// It's OK to export a nil experience
	if p == nil {
		return nil
	}

	pld := p.ExportPublic()
	pld.Payload = NewPrivatePayload(p.User)
	return pld
}

// ProfilesPayload represents a list of Profiles that can be
// safely returned to the clients
type ProfilesPayload struct {
	Results []*ProfilePayload `json:"results"`
}

// ExportPublic returns a ProfilesPayload containing only the fields that are safe to
// be seen by anyone
func (profiles Profiles) ExportPublic() *ProfilesPayload {
	pld := &ProfilesPayload{}
	pld.Results = make([]*ProfilePayload, len(profiles))
	for i, p := range profiles {
		pld.Results[i] = p.ExportPublic()
	}
	return pld
}

// ExportPrivate returns a ProfilesPayload containing all the fields
func (profiles Profiles) ExportPrivate() *ProfilesPayload {
	pld := &ProfilesPayload{}
	pld.Results = make([]*ProfilePayload, len(profiles))
	for i, p := range profiles {
		pld.Results[i] = p.ExportPrivate()
	}
	return pld
}
