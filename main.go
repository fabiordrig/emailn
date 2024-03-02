package main

import (
	"emailn/domains/campaign"

	"github.com/go-playground/validator/v10"
)

func main() {

	campaign := campaign.Campaign{}

	validate := validator.New()

	err := validate.Struct(campaign)

	if err != nil {
		errors := err.(validator.ValidationErrors)

		for _, e := range errors {
			println(e.Error())
			println(e.Tag())
		}

	}

}
