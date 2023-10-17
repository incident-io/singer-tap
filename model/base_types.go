package model

type dateTime struct{}

var DateTime dateTime

func (dateTime) Schema() Property {
	return Property{
		Types:        []string{"string"},
		CustomFormat: "date-time",
	}
}

func (dateTime) Serialize(input string) any {
	return input
}
