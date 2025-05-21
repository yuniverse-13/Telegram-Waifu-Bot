package characters

import (
	"math/rand"
	"strings"
)

type Character struct {
	ID          int
	Name        string
	AltNames    []string
	Description string
	ImageURL    string
	Rating      float64
}

var charactersStore []Character

func init() {
	charactersStore = []Character{
		{
			ID:          1,
			Name:        "Фрирен",
			AltNames:    []string{"Frieren"},
			Description: "Эльфийская волшебница, бывший член отряда героев, победившего Короля Демонов.",
			ImageURL:    "https://i.pinimg.com/736x/87/b3/08/87b308eb264e7ee8c8bdc77774b87044.jpg",
			Rating:      9.8,
		},
		{
			ID:          2,
			Name:        "Мегумин",
			AltNames:    []string{"Мегу", "Megumin", "Megu"},
			Description: "Архимаг из клана Алых Демонов в фэнтезийном мире и первый человек, присоединившийся к группе Казумы.",
			ImageURL:    "https://i.pinimg.com/736x/be/f3/7a/bef37a878ab694f6722533b3c3222f1a.jpg",
			Rating:      9.9,
		},
		{
			ID:          3,
			Name:        "Ферн",
			AltNames:    []string{"Fern"},
			Description: "Юная, но удивительно серьёзная и ответственная волшебница, ученица Фрирен, которая часто выступает её своеобразным опекуном и голосом разума.",
			ImageURL:    "https://i.pinimg.com/736x/59/7b/2b/597b2b1f7906d73d87432150f181cb91.jpg",
			Rating:      9.5,
		},
		{
			ID:          4,
			Name:        "Сильфиетта",
			AltNames:    []string{"Сильфи", "Фиттс", "Sylphiette", "Sylphy"},
			Description: "Преданная эльфийка-полукровка, подруга детства и жена Рудеуса, выросшая из робкой девочки в талантливую волшебницу, также известная как Фиттс.",
			ImageURL:    "https://i.pinimg.com/736x/57/a7/4b/57a74bd586a877e4a75ad66e0b92fc5f.jpg",
			Rating:      8.8,
		},
		{
			ID:          5,
			Name:        "Рокси",
			AltNames:    []string{"Roxy"},
			Description: "Это талантливая и добрая демоница-мигурд, первая и глубоко уважаемая учительница магии Рудеуса.",
			ImageURL:    "https://i.pinimg.com/736x/5d/5c/93/5d5c93e86d2713c9dc7d758676382bc2.jpg",
			Rating:      8.9,
		},
	}
}

func GetCharacterByID(id int) (Character, bool) {
	for _, char := range charactersStore {
		if id == char.ID {
			return char, true
		}
	}
	return Character{}, false
}

func GetCharacterByNameOrAlt(name string) (Character, bool) {
	searchNameLower := strings.ToLower(name)
	for _, char := range charactersStore {
		if strings.ToLower(char.Name) == searchNameLower {
			return char, true
		}
		for _, altName := range char.AltNames {
			if strings.ToLower(altName) == searchNameLower {
				return char, true
			}
		}
	}
	return Character{}, false
}

func GetRandomCharacter() (Character, bool) {
	if len(charactersStore) == 0 {
		return Character{}, false
	}
	randomIndex := rand.Intn(len(charactersStore))
	return charactersStore[randomIndex], true
}