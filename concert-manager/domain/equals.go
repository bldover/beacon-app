package domain

func (e Event) Equals(o Event) bool {
	return e.ID.Primary == o.ID.Primary
}

func (e Event) EqualsFields(o Event) bool {
	return (e.MainAct == o.MainAct || e.MainAct.EqualsFields(*o.MainAct)) && e.Venue.EqualsFields(o.Venue) && e.Date == o.Date
}

func (a Artist) Equals(o Artist) bool {
	return a.ID.Primary == o.ID.Primary
}

func (a Artist) EqualsFields(o Artist) bool {
	return a.Name == o.Name
}

func (v Venue) Equals(o Venue) bool {
	return v.ID == o.ID
}

func (v Venue) EqualsFields(o Venue) bool {
	return v.Name == o.Name && v.City == o.City && v.State == o.State
}
