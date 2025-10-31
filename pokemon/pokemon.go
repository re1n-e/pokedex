package pokemon

const pokemonApi = "https://pokeapi.co/api/v2/pokemon/"

type Pokemon struct {
	Name      string         `json:"name"`
	Abilities []AbilityEntry `json:"abilities"`
	BaseExp   int            `json:"base_experience"`
	Height    int            `json:"height"`
	Weight    int            `json:"weight"`
	Moves     []MoveEntry    `json:"moves"`
	Stats     []Stats        `json:"stats"`
	Type      []Type         `json:"types"`
}

type AbilityEntry struct {
	Ability Ability `json:"ability"`
}

type Ability struct {
	Name string `json:"name"`
}

type MoveEntry struct {
	Move Move `josn:"move"`
}

type Move struct {
	Name string `json:"name"`
}

type Stats struct {
	Base int      `json:"base_stat"`
	Name StatName `json:"stat"`
}

type StatName struct {
	Name string `json:"name"`
}

type Type struct {
	TypeName TypeName `json:"type"`
}

type TypeName struct {
	Name string `json:"name"`
}
