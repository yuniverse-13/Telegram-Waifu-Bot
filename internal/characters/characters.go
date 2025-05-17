package characters

type Character struct {
	ID          int
	Name        string
	Description string
	ImageURL    string
	Rating      int
}

func GetSampleCharacter() Character {
	return Character{
		ID:          1,
		Name:        "Frieren",
		Description: "1000+ years old elf",
		ImageURL:    "https://i.pinimg.com/736x/99/31/15/9931153152416410957e712bb13ce8ad.jpg",
		Rating:      10,
	}
}

func GetAnotherSampleCharacter() Character {
	return Character{
		ID:          2,
		Name:        "Megumin",
		Description: "Red Crimson Demon",
		ImageURL:    "https://i.pinimg.com/736x/be/f3/7a/bef37a878ab694f6722533b3c3222f1a.jpg",
		Rating:      10,
	}
}
