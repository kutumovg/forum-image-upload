package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// initDB initializes the database connection and creates tables if they don't exist.
func initDB() (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	createTables(db)

	return db, nil
}

// createTables defines the SQL schema for the forum database and creates tables if they don't exist.
func createTables(db *sql.DB) {
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT UNIQUE,
        username TEXT UNIQUE,
        password TEXT
    );`

	createPostsTable := `
    CREATE TABLE IF NOT EXISTS posts (
        id TEXT PRIMARY KEY,
        user_id TEXT,
        content TEXT,
        created_at DATETIME,
        likes INTEGER DEFAULT 0,
        dislikes INTEGER DEFAULT 0,
		image_path TEXT,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

	createPostLikesTable := `
    CREATE TABLE IF NOT EXISTS post_likes (
        id TEXT PRIMARY KEY,
        user_id TEXT,
        post_id TEXT,
        is_like BOOLEAN,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (post_id) REFERENCES posts(id),
        UNIQUE (user_id, post_id)
    );`

	createCommentsTable := `
    CREATE TABLE IF NOT EXISTS comments (
        id TEXT PRIMARY KEY,
        post_id TEXT,
        user_id TEXT,
        content TEXT,
        created_at DATETIME,
        likes INTEGER DEFAULT 0,
        dislikes INTEGER DEFAULT 0,
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

	createCommentLikesTable := `
    CREATE TABLE IF NOT EXISTS comment_likes (
        id TEXT PRIMARY KEY,
        user_id TEXT,
        comment_id TEXT,
        is_like BOOLEAN,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (comment_id) REFERENCES comments(id),
        UNIQUE (user_id, comment_id)
    );`

	createCategoriesTable := `
    CREATE TABLE IF NOT EXISTS categories (
        id TEXT PRIMARY KEY,
        name TEXT UNIQUE
    );`

	createPostCategoriesTable := `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id TEXT,
		category_id TEXT,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`
	// Execute the table creation commands
	_, err := db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createCategoriesTable)
	if err != nil {
		log.Fatal(err)
	}
	seedCategories(db)

	_, err = db.Exec(createPostCategoriesTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createPostsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createPostLikesTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createCommentsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createCommentLikesTable)
	if err != nil {
		log.Fatal(err)
	}

	// seedData(db)
}

// seedCategories inserts default categories into the categories table.
func seedCategories(db *sql.DB) {
	categories := []string{"Autobiography", "Comedy", "Science Fiction", "Fantasy", "Mystery", "Other"}

	for _, category := range categories {
		categoryID, _ := uuid.NewV4()
		_, err := db.Exec("INSERT OR IGNORE INTO categories (id, name) VALUES (?, ?)", categoryID.String(), category)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// seedData populates the comments table with initial sample data for testing.
func seedData(db *sql.DB) error {
	comments := []struct {
		post_id    string
		user_id    string
		content    string
		created_at time.Time
		likes      int
		dislikes   int
	}{
		{"f0629e27-9ffa-4ec4-bb5b-1ea517c6960c", "4421ccd6-555f-4f4c-b478-1f0d41803e05", "Согласен, 'Преступление и наказание' - глубочайший роман!", time.Date(2024, 10, 31, 14, 30, 0, 0, time.UTC), 5, 2},
		{"f0629e27-9ffa-4ec4-bb5b-1ea517c6960c", "dcfac435-eac4-40ff-a685-7e4340b8413e", "Потрясающее произведение, заставляет задуматься!", time.Date(2024, 10, 31, 15, 0, 0, 0, time.UTC), 6, 3},
		{"7f1fa671-9d63-4fd4-99a8-3e4a1b154e59", "9fa6c47c-32cf-4f01-8478-57aaaf5d9872", "Булгаков - настоящий гений! Столько скрытых смыслов.", time.Date(2024, 10, 31, 16, 0, 0, 0, time.UTC), 4, 1},
		{"abb20c13-c391-465f-bab6-ac88c7b92bb4", "dcfac435-eac4-40ff-a685-7e4340b8413e", "Оруэлл был провидцем. '1984' и сегодня актуален.", time.Date(2024, 10, 31, 17, 0, 0, 0, time.UTC), 7, 1},
		{"3c37e9f7-baaa-407c-9ae6-f94038314ffe", "4421ccd6-555f-4f4c-b478-1f0d41803e05", "Трагедия и правда о любви и жизни в России.", time.Date(2024, 10, 31, 18, 0, 0, 0, time.UTC), 8, 2},
		{"3c37e9f7-baaa-407c-9ae6-f94038314ffe", "9fa6c47c-32cf-4f01-8478-57aaaf5d9872", "Толстой сделал огромный вклад в мировую литературу.", time.Date(2024, 10, 31, 19, 0, 0, 0, time.UTC), 5, 2},
		{"0506b9b2-b808-45c6-a1dd-a0115fa4ee57", "17cddc75-647d-4f74-b7fa-5d7953985f41", "Сатира Булгакова - шедевр!", time.Date(2024, 10, 31, 20, 0, 0, 0, time.UTC), 6, 4},
		{"41086f47-be31-49cb-ba97-cbedeba1656f", "f718aafc-34dd-47ad-b4a5-d6ffe6301e13", "Читал Шантарам на одном дыхании!", time.Date(2024, 10, 31, 21, 0, 0, 0, time.UTC), 4, 1},
		{"4c95cb45-d54c-45d7-86e0-0fc5ccbd2053", "36a30519-aa35-4dcd-9987-95ca9911e5bb", "Алиса - это не просто сказка. Здесь целая философия.", time.Date(2024, 10, 31, 22, 0, 0, 0, time.UTC), 10, 0},
		{"6e0a6731-41e0-43e2-bfee-a71e13401774", "8e40c285-40d4-47fe-9754-25013dffc8a1", "Маленький принц - особенная книга для каждого возраста.", time.Date(2024, 10, 31, 23, 0, 0, 0, time.UTC), 3, 1},
		{"84cbe293-1aea-4fd8-b8f0-7fdce6b37c87", "f718aafc-34dd-47ad-b4a5-d6ffe6301e13", "Камю заставляет задуматься о смысле жизни.", time.Date(2024, 10, 31, 14, 30, 0, 0, time.UTC), 6, 2},
		{"b8954b27-defc-4379-940e-6bcea0adcda5", "0add8655-f69a-4362-8ae2-1c8fb87ff1bb", "Улисс очень сложен, но интересен.", time.Date(2024, 10, 31, 15, 0, 0, 0, time.UTC), 9, 3},
		{"cb1fd515-07b4-460e-9e3a-2feb6e417399", "7e8bc625-6ee8-4717-ab12-e9d919ea6d71", "Гарри Поттер - культовая история.", time.Date(2024, 10, 31, 16, 0, 0, 0, time.UTC), 5, 0},
		{"a1bf6607-9ee2-458d-b296-3d47439e3645", "36a30519-aa35-4dcd-9987-95ca9911e5bb", "Люблю детективы Агаты Кристи.", time.Date(2024, 10, 31, 17, 0, 0, 0, time.UTC), 4, 1},
		{"bc1a75a0-3912-4909-9bb0-01c5784290b9", "17cddc75-647d-4f74-b7fa-5d7953985f41", "Грозовой перевал - драма на все времена.", time.Date(2024, 10, 31, 18, 0, 0, 0, time.UTC), 7, 2},
		{"10f05449-55b8-4e16-987a-bbcd385bcc42", "cc3e5995-938d-4259-bd6f-395e3f0205ff", "Шерлок - король детективов.", time.Date(2024, 10, 31, 19, 0, 0, 0, time.UTC), 3, 1},
		{"5ac23c7e-ca5c-4f28-a5dc-25a440ed40eb", "8e40c285-40d4-47fe-9754-25013dffc8a1", "Кафка создал свою уникальную вселенную.", time.Date(2024, 10, 31, 20, 0, 0, 0, time.UTC), 6, 2},
		{"1446dee2-e4c8-4e93-b4f1-47b5a3a03222", "cc3e5995-938d-4259-bd6f-395e3f0205ff", "Ремарк точно передал дух войны.", time.Date(2024, 10, 31, 21, 0, 0, 0, time.UTC), 8, 1},
		{"1761ffae-8b00-4809-b50b-63f29bc7c18a", "0add8655-f69a-4362-8ae2-1c8fb87ff1bb", "Ромео и Джульетта - шедевр романтики.", time.Date(2024, 10, 31, 22, 0, 0, 0, time.UTC), 7, 3},
		{"11685b40-0a1a-42e9-8bf0-e512e7b97cbf", "7e8bc625-6ee8-4717-ab12-e9d919ea6d71", "Фэнтези Толкина вдохновляет.", time.Date(2024, 10, 31, 23, 0, 0, 0, time.UTC), 5, 1},
	}

	for _, comment := range comments {
		id, _ := uuid.NewV4()
		_, err := db.Exec(
			`INSERT OR IGNORE INTO comments (id, post_id, user_id, content, created_at, likes, dislikes) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			id.String(), comment.post_id, comment.user_id, comment.content, comment.created_at.Format("2006-01-02 15:04:05"), comment.likes, comment.dislikes,
		)
		if err != nil {
			log.Printf("Error inserting user %s: %v\n", comment.post_id, err)
			return err
		}
	}

	return nil
}
