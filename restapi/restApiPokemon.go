package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Sprites struct {
	BackDefault       *string `json:"back_default"`
	BackFemale        *string `json:"back_female"`
	BackShiny         *string `json:"back_shiny"`
	BackShinyFemale   *string `json:"back_shiny_female"`
	FrontDefault      *string `json:"front_default"`
	FrontFemale       *string `json:"front_female"`
	FrontShiny        *string `json:"front_shiny"`
	FrontShinyFemale  *string `json:"front_shiny_female"`
}

type PokemonResponse struct {
	Stats   []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name *string `json:"name"`
			Url  *string `json:"url"`
		}
	}
	Name    string  `json:"name"`
	Sprites Sprites `json:"sprites"`
}

func getPokemonData(id string) (*PokemonResponse, error) {
	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", id)
	formURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-form/%s/", id)

	pokemonResp, err := http.Get(pokemonURL)
	if err != nil {
		return nil, err
	}
	defer pokemonResp.Body.Close()

	pokemonBody, err := ioutil.ReadAll(pokemonResp.Body)
	if err != nil {
		return nil, err
	}

	var pokemon PokemonResponse
	err = json.Unmarshal(pokemonBody, &pokemon)
	if err != nil {
		return nil, err
	}

	formResp, err := http.Get(formURL)
	if err != nil {
		return nil, err
	}
	defer formResp.Body.Close()

	formBody, err := ioutil.ReadAll(formResp.Body)
	if err != nil {
		return nil, err
	}

	var form struct {
		Name    string  `json:"name"`
		Sprites Sprites `json:"sprites"`
	}
	err = json.Unmarshal(formBody, &form)
	if err != nil {
		return nil, err
	}

	pokemon.Name = form.Name
	pokemon.Sprites = form.Sprites

	return &pokemon, nil
}

func handlePokemonRequest(c *gin.Context) {
	var requestData struct {
		ID string `json:"id"`
	}
	err := c.BindJSON(&requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pokemon, err := getPokemonData(requestData.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

func main() {
	router := gin.Default()
	router.POST("/pokemon", handlePokemonRequest)
	fmt.Println("Starting server on port 8080...")
	router.Run(":8080")
}