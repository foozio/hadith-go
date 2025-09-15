package data

// Hadith represents a single hadith entry from a book JSON file.
type Hadith struct {
    Book   string `json:"book"`
    Number int    `json:"number"`
    Arab   string `json:"arab"`
    ID     string `json:"id"`
}

// Book holds all hadiths for a particular collection.
type Book struct {
    Name    string   `json:"name"`
    Hadiths []Hadith `json:"hadiths"`
}

