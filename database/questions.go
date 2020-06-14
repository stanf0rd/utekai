package database

import (
	"fmt"
	"log"
	"strings"
)

// Question - basic question type
// Body - the question itself
// Order - the preferred order for asking this question (firts|random|last)
type Question struct {
	ID    int
	Body  string
	Order string
}

func initQuestions() error {
	count, err := GetQuestionCount()
	if err != nil {
		return fmt.Errorf("Cannot read questions table: %v", err)
	}

	if count == 0 {
		log.Println("Questions table is empty, inserting initial questions...")

		if err := AddQuestions(questions); err != nil {
			return fmt.Errorf("Cannot add initial questions: %v", err)
		}
	} else {
		log.Printf("Found %d already added questions", count)
	}

	return nil
}

// GetQuestionByID returns questions found by id
func GetQuestionByID(ID int) (*Question, error) {
	row := db.QueryRow(`
		SELECT id, body, "order" FROM questions
		WHERE id = $1
	`, ID)

	var q Question
	if err := row.Scan(&q.ID, &q.Body, &q.Order); err != nil {
		return nil, fmt.Errorf("Cannot get question by ID: %v", err)
	}

	return &q, nil
}

// AddQuestions adds all questions from array to DB
// writes userID in struct
func AddQuestions(questions []Question) error {
	query := `
		INSERT INTO questions(body, "order")
		VALUES
	`
	values := []interface{}{}

	for i, q := range questions {
		query += fmt.Sprintf("($%d, $%d),", (i+1)*2-1, (i+1)*2)
		values = append(values, q.Body, q.Order)
	}

	query = strings.TrimSuffix(query, ",")

	statement, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Unable to create statement: %v", err)
	}

	result, err := statement.Exec(values...)
	if err != nil {
		return fmt.Errorf("Unable to add questions: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Printf("Questions added, but cannot get stats")
	} else {
		log.Printf("%d questions added", count)
	}

	return nil
}

// GetQuestionCount gets count of questions already added to database
func GetQuestionCount() (count int, err error) {
	row := db.QueryRow("SELECT COUNT(*) FROM questions;")

	if err = row.Scan(&count); err != nil {
		return 0, fmt.Errorf("Cannot get count of questions: %v", err)
	}

	return count, nil
}

// GetAllQuestions returns all questions collected in database
func GetAllQuestions() ([]Question, error) {
	rows, err := db.Query("SELECT \"id\", \"order\", body FROM questions")
	if err != nil {
		return nil, fmt.Errorf("Cannot get questions from database: %v", err)
	}
	defer rows.Close()
	questions := make([]Question, 0)

	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.Order, &q.Body); err != nil {
			return nil, fmt.Errorf("Cannot read questions from database: %v", err)
		}
		questions = append(questions, q)
	}

	return questions, nil
}

// GetQuestionIDsByOrder gets question ids by order from database (like these linting rules)
func GetQuestionIDsByOrder(order string) (res []int, err error) {
	rows, err := db.Query("SELECT id FROM questions WHERE \"order\" = $1", order)
	if err != nil {
		return nil, fmt.Errorf("Cannot get questions from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ID int
		if err := rows.Scan(&ID); err != nil {
			return nil, fmt.Errorf("Cannot read question IDs from database: %v", err)
		}
		res = append(res, ID)
	}

	return res, nil
}

var questions = []Question{
	{Body: "Как ты остановился? Как ты реагируешь на призыв остановиться?", Order: "first"},
	{Body: "Обрати внимание на тело и опиши, что оно чувствует. Самые яркие ощущения. Теперь едва уловимые.", Order: "random"},
	{Body: "На что первое ты обратил внимание, когда остановился?", Order: "random"},
	{Body: "Опиши процесс, который ты прервал.", Order: "random"},
	{Body: "Опиши поверхность, на которую ты опираешься и как ты её чувствуешь.", Order: "random"},
	{Body: "Что двигается, пока ты смотришь?", Order: "random"},
	{Body: "Какой предмет вокруг больше всего похож на тебя? На твои чувства?", Order: "random"},
	{Body: "Попробуй переназвать любой предмет вокруг, исходя из его вида и ощущений от него. Опиши процесс.", Order: "random"},
	{Body: "Представь ощущение от прикосновения к любому предмету недалеко от себя. Потрогай его. Похоже? Чем отличается? ", Order: "random"},
	{Body: "Обрати внимание на дыхание, какое оно?", Order: "random"},
	{Body: "Чего ты сейчас хочешь?", Order: "random"},
	{Body: "На что смотреть приятнее всего? Почему это приятно?", Order: "random"},
	{Body: "Где ты находишься? И как тебе здесь? Что бы ты изменил?", Order: "random"},
	{Body: "Что ты сейчас чувствуешь?", Order: "random"},
	{Body: "Какая скорость у тебя была до остановки? Какая сейчас?", Order: "random"},
	{Body: "Как ты ощущаешь прикосновение одежды?", Order: "random"},
	{Body: "Какие звуки ты слышишь? Они вызывают какие-то чувства?", Order: "random"},
	{Body: "Издай любой звук. Какой?", Order: "random"},
	{Body: "Какой кусочек тела самый холодный? Тёплый?", Order: "random"},
	{Body: "Осмотрись. Что тебе здесь нравится? А что нет?", Order: "random"},
	{Body: "Попробуй почувствовать материальность любого предмета рядом. Как это?", Order: "random"},
	{Body: "Какой цвет тебе сейчас ближе из всех вокруг? Как ты себя чувствуешь, смотря на него?", Order: "random"},
	{Body: "Почему ты последовал/не последовал предложению остановиться?", Order: "random"},
	{Body: "Что-то изменилось за время, которое ты останавливаешься?", Order: "last"},
}
