package model2

type MCountryList struct {
	byNameJa map[string]*MCountry
}

func (mc *MCountryList) GetByNameJA(name string) (*MCountry, bool) {
	c, ok := mc.byNameJa[name]
	return c, ok
}

func NewMCountryList(values []*MCountry) *MCountryList {
	c := &MCountryList{
		byNameJa: make(map[string]*MCountry, 1000),
	}
	for _, v := range values {
		c.byNameJa[v.NameJa] = v
	}
	return c
}
